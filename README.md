##NewPeerWithEngine分析

在`main.go` 的`func serve`中，有一句：
```
		peerServer, err = peer.NewPeerWithEngine(secHelperFunc, helper.GetEngine)
```
本文件对这句代码所执行的操作进行分析。

在`core/peer/peer.go`文件中，NewPeerWithEngine 代码实现如下：

```
// NewPeerWithEngine returns a Peer which uses the supplied handler factory function for creating new handlers on new Chat service invocations.

func NewPeerWithEngine(secHelperFunc func() crypto.Peer, engFactory EngineFactory) (peer *PeerImpl, err error) {
       peer = new(PeerImpl)
       peerNodes := peer.initDiscovery()

       peer.handlerMap = &handlerMap{m: make(map[pb.PeerID]MessageHandler)}

       peer.isValidator = ValidatorEnabled()
       peer.secHelper = secHelperFunc()

       // Install security object for peer
       if SecurityEnabled() {
              if peer.secHelper == nil {
                     return nil, fmt.Errorf("Security helper not provided")
              }
       }

       // Initialize the ledger before the engine, as consensus may want to begin interrogating the ledger immediately
       ledgerPtr, err := ledger.GetLedger()
       if err != nil {
              return nil, fmt.Errorf("Error constructing NewPeerWithHandler: %s", err)
       }
       peer.ledgerWrapper = &ledgerWrapper{ledger: ledgerPtr}

       peer.engine, err = engFactory(peer)
       if err != nil {
              return nil, err
       }
       peer.handlerFactory = peer.engine.GetHandlerFactory()
       if peer.handlerFactory == nil {
              return nil, errors.New("Cannot supply nil handler factory")
       }

       peer.chatWithSomePeers(peerNodes)
       return peer, nil

}
```
####1、在第一行中，new了PeerImpl 一个实例
```
peer = new(PeerImpl)
```
PeerImpl定义如下：
```
// PeerImpl implementation of the Peer service
type PeerImpl struct {
	handlerFactory HandlerFactory
	handlerMap     *handlerMap
	ledgerWrapper  *ledgerWrapper
	secHelper      crypto.Peer
	engine         Engine
	isValidator    bool
	reconnectOnce  sync.Once
	discHelper     discovery.Discovery
	discPersist    bool
}
```
注：PeerImpl 实现了 MessageHandlerCoordinator 接口中的所有方法。
PeerImpl 实现的方法集：
```
type MessageHandlerCoordinator interface {
	Peer
	SecurityAccessor
	BlockChainAccessor
	BlockChainModifier
	BlockChainUtil
	StateAccessor
	RegisterHandler(messageHandler MessageHandler) error
	DeregisterHandler(messageHandler MessageHandler) error
	Broadcast(*pb.Message, pb.PeerEndpoint_Type) []error
	Unicast(*pb.Message, *pb.PeerID) error
	GetPeers() (*pb.PeersMessage, error)
	GetRemoteLedger(receiver *pb.PeerID) (RemoteLedger, error)
	PeersDiscovered(*pb.PeersMessage) error
	ExecuteTransaction(transaction *pb.Transaction) *pb.Response
	Discoverer
}
```
PeerImpl 结构中定义的变量：
- HandlerFactory 是用来创建一个 MessageHandlers 。MessageHandler 是 Peer 的一个消息处理器。这里面只是定义了方法，可能会有多种不同的实现。
- HandlerMap 是存储多个 MessageHandler 的一个 map ，实现了 peerId 与 MessageHandler 的映射。
- LedgerWrapper 存放了 ledger 的指针。
- secHelper 是 crypto.Peer 接口，是加密签名等方法的集合。
- engine 是管理网络通信和交易处理的引擎，和共识算法有关。并且提供了 GetHandlerFactory方法。
- isValidator 是一个节点是 vp 节点还是 nvp 节点的标志。
- reconnectOnce 是只执行一次。用来与节点之间进行连接。
- disHelpler 是 discovery 的一个工具类
- discPersist 是将发现的节点是否需要存储在数据库的一个标志。

####2、在第二步中，初始化 Discovery 过程。
```
peerNodes := peer.initDiscovery()
```
在`peer.initDiscovery`中，首先初始化相关实例
```
    p.discHelper = discovery.NewDiscoveryImpl()
	p.discPersist = viper.GetBool("peer.discovery.persist")
```
之后根据`p.discPersist`从数据库中读取之前存储的节点：
```
addresses, err := p.LoadDiscoveryList()
```

之后根据配置好的`CORE_PEER_DISCOVERY_ROOTNODE`环境变量，得到一个address数组。在目前的4节点pbft集群中，vp0的address为空，vp1~vp3的address数组长度均为1.

到目前为止，所有vp节点都没有获取到整个vp网络所有节点的信息。

####3、第三步中，初始化 handlerMap 变量
```
peer.handlerMap = &handlerMap{m: make(map[pb.PeerID]MessageHandler)}
```
只完成了初始化，其中并没有任何数据。

####4、第四步中，主要完成的工作是初始化 isValidator 变量
```
peer.isValidator = ValidatorEnabled()
```
ValidatorEnabled 实现的功能是读取配置文件，设置该节点VP flag以及设置该节点是否运行自动检测协议。

####5、第五步中，主要完成的是 crypto.Peer 的注册和初始化。并且通过helper的形式调用相关的方法。
```
peer.secHelper = secHelperFunc()
```
secHelperFunc() 是 newPeerWithEngine 被调用的时候实现的一个参数。
对于 secHelperFunc（），主要完成两个功能，节点的注册和初始化。
```
crypto.RegisterValidator(enrollID, nil, enrollID, enrollSecret)
crypto.RegisterPeer(enrollID, nil, enrollID, enrollSecret)
```
是通过向 membersrvc 进行注册，变成一个 VP 节点。
以 register 为例，主要函数调用过程如下：
- crypto.RegisterValidator(enrollID, nil, enrollID, enrollSecret)
- validator.register(name, pwd, enrollID, enrollPWD, nil)
- validator.peerImpl.register(NodeValidator, id, pwd, enrollID, enrollPWD, nil)
- peer.nodeImpl.register(eType, name, pwd, enrollID, enrollPWD, regFunc)
- node.nodeRegister(eType, name, pwd, enrollID, enrollPWD)
- node.registerCryptoEngine(enrollID, enrollPWD)
- node.retrieveECACertsChain(enrollID)，node.retrieveTCACertsChain(enrollID)， node.retrieveTLSCertificate(enrollID, enrollPWD)
- ecaCertRaw, err := node.getECACertificate()
- responce, err := node.callECAReadCACertificate(context.Background())
- conn, err := node.getClientConn(node.conf.getECAPAddr(), node.conf.getECAServerName())
- return viper.GetString(conf.ecaPAddressProperty)
- conf.ecaPAddressProperty = "peer.pki.eca.paddr"

其中`peer.pki.eca.paddr`即为配置在`membersrvc.yaml`中的环境变量：
```
CORE_PEER_PKI_ECA_PADDR=membersrvc:7054
```
通过上述一系列函数，获得到了 ECA、TCA、TLSCA 的证书。

####6、初始化ledger
```
    ledgerPtr, err := ledger.GetLedger()
	if err != nil {
		return nil, fmt.Errorf("Error constructing NewPeerWithHandler: %s", err)
	}
	peer.ledgerWrapper = &ledgerWrapper{ledger: ledgerPtr}
```
`GetLedger` 创建了一个单例 ledger，并且 ledger 的创建应该放在共识引擎 engine 初始化之前。

####7、初始化共识引擎 engine 及创建相应的MessageHandler
```
    peer.engine, err = engFactory(peer)
	if err != nil {
		return nil, err
	}
	peer.handlerFactory = peer.engine.GetHandlerFactory()
	if peer.handlerFactory == nil {
		return nil, errors.New("Cannot supply nil handler factory")
	}
```

`handler`、`messageHandler`、`MessageHandlerCoordinator`的作用.
消息通信相关：
FSM、Handler、Consensus engine。
