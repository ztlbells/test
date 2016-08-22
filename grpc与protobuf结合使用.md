## grpc 与 protobuf 结合案例
#### 1、背景
protocolbuffer(以下简称PB)是google 的一种数据交换的格式，它独立于语言，独立于平台。google 提供了多种语言的实现：java、c#、c++、go 和 python，每一种实现都包含了相应语言的编译器以及库文件。由于它是一种二进制的格式，比使用 xml 进行数据交换快许多。可以把它用于分布式应用之间的数据通信或者异构环境下的数据交换。作为一种效率和兼容性都很优秀的二进制数据传输格式，可以用于诸如网络传输、配置文件、数据存储等诸多领域。

gRPC是一个高性能、通用的开源RPC框架，其由Google主要面向移动应用开发并基于HTTP/2协议标准而设计，基于ProtoBuf(Protocol Buffers)序列化协议开发，且支持众多开发语言。gRPC提供了一种简单的方法来精确地定义服务和为iOS、Android和后台支持服务自动生成可靠性很强的客户端功能库。

##### 2、protobuf grpc http2 stream


```
type PeerServer interface {
	// Accepts a stream of Message during chat session, while receiving
	// other Message (e.g. from other peers).
	Chat(Peer_ChatServer) error
	// Process a transaction from a remote source.
	ProcessTransaction(context.Context, *Transaction) (*Response, error)
}
```
```
type PeerClient interface {
	// Accepts a stream of Message during chat session, while receiving
	// other Message (e.g. from other peers).
	Chat(ctx context.Context, opts ...grpc.CallOption) (Peer_ChatClient, error)
	// Process a transaction from a remote source.
	ProcessTransaction(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*Response, error)
}
```

PeerServer端的实现：
```
// Chat implementation of the the Chat bidi streaming RPC function
func (p *PeerImpl) Chat(stream pb.Peer_ChatServer) error {
	return p.handleChat(stream.Context(), stream, false)
}
```

```
// ProcessTransaction implementation of the ProcessTransaction RPC function
func (p *PeerImpl) ProcessTransaction(ctx context.Context, tx *pb.Transaction) (response *pb.Response, err error) {
	peerLogger.Debugf("ProcessTransaction processing transaction uuid = %s", tx.Uuid)
	// Need to validate the Tx's signature if we are a validator.
	if p.isValidator {
		// Verify transaction signature if security is enabled
		secHelper := p.secHelper
		if nil != secHelper {
			peerLogger.Debugf("Verifying transaction signature %s", tx.Uuid)
			if tx, err = secHelper.TransactionPreValidation(tx); err != nil {
				peerLogger.Errorf("ProcessTransaction failed to verify transaction %v", err)
				return &pb.Response{Status: pb.Response_FAILURE, Msg: []byte(err.Error())}, nil
			}
		}

	}
	return p.ExecuteTransaction(tx), err
}
```

PeerClient端的调用：
```
func (p *PeerImpl) chatWithPeer(address string) error {
	peerLogger.Debugf("Initiating Chat with peer address: %s", address)
	conn, err := NewPeerClientConnectionWithAddress(address)
	if err != nil {
		peerLogger.Errorf("Error creating connection to peer address %s: %s", address, err)
		return err
	}
	serverClient := pb.NewPeerClient(conn)
	ctx := context.Background()
	stream, err := serverClient.Chat(ctx)
	if err != nil {
		peerLogger.Errorf("Error establishing chat with peer address %s: %s", address, err)
		return err
	}
	peerLogger.Debugf("Established Chat with peer address: %s", address)
	err = p.handleChat(ctx, stream, true)
	stream.CloseSend()
	if err != nil {
		peerLogger.Errorf("Ending Chat with peer address %s due to error: %s", address, err)
		return err
	}
	return nil
}
```


PeerImpl本身实现了所有的 MessageHandlerCoordinator 本身，包括：
```
RegisterHandler(messageHandler MessageHandler) error
DeregisterHandler(messageHandler MessageHandler) error
```

在 peer.go 中有一个struct
```
type ChatStream interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}
```
每次使用的时候都是被传入进去的，stream 的值是 `pb.Peer_ChatServer`。

##### 3、对于重点分析的代码：

```
    peer.engine, err = engFactory(peer)
	if err != nil {
		return nil, err
	}
	peer.handlerFactory = peer.engine.GetHandlerFactory()
```

messageHandler，MessageHandlerCoordinator，handler，PeerImpl，ConsensusHandler。
MessageHandlerCoordinator 是一套方法集合。
messageHandler 也是一套方法集合。

handler、ConsensusHandler、PeerImpl 都是结构体。
它们的定义如下：
PeerImpl:
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

MessageHandlerCoordinator:
```
// MessageHandlerCoordinator responsible for coordinating between the registered MessageHandler's
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

在 peer.go 中,有许多：
```
func (p *PeerImpl) GetPeers() (*pb.PeersMessage, error) {
```
可见 MessageHandlerCoordinator 的方法全部被 PeerImpl 实现。MessageHandlerCoordinator 提供了许多方法可以操纵修改 PeerImpl 的数据结构。

MessageHandler 定义：
```
// MessageHandler standard interface for handling Openchain messages.
type MessageHandler interface {
	RemoteLedger
	HandleMessage(msg *pb.Message) error
	SendMessage(msg *pb.Message) error
	To() (pb.PeerEndpoint, error)
	Stop() error
}
```

MessageHandler 也是一套方法集合 ，被 handler 全部实现，并且目前ConsensusHandler 部分实现。（因此`ConsensusHandler`不能作为在函数中`MessageHandler`参数的替代体）
handler 定义：
```
// Handler peer handler implementation.
type Handler struct {
	chatMutex                     sync.Mutex
	ToPeerEndpoint                *pb.PeerEndpoint
	Coordinator                   MessageHandlerCoordinator
	ChatStream                    ChatStream
	doneChan                      chan struct{}
	FSM                           *fsm.FSM
	initiatedStream               bool // Was the stream initiated within this Peer
	registered                    bool
	syncBlocks                    chan *pb.SyncBlocks
	snapshotRequestHandler        *syncStateSnapshotRequestHandler
	syncStateDeltasRequestHandler *syncStateDeltasHandler
	syncBlocksRequestHandler      *syncBlocksRequestHandler
}
```
ChatStream 就是 pb 的流，可以用来 send，resv。

ConsensusHandler：
```
// ConsensusHandler handles consensus messages.
// It also implements the Stack.
type ConsensusHandler struct {
	peer.MessageHandler
	consenterChan chan *util.Message
	coordinator   peer.MessageHandlerCoordinator
}
```


##### 4、消息通信函数及消息处理函数（包括使用FSM）
所有的消息通信函数如下：
#####1、chatstream
```
// ChatStream interface supported by stream between Peers
type ChatStream interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}
```
`ChatStream` 就是 pb.chatServer 的流，可以用来 send，resv。

```
type Peer_ChatServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}
```
```
type Peer_ChatClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}
```
上面两者的`send`、`Recv`在不同的时机被不同的函数被调用。

#####2、MessageHandler
```
    RemoteLedger
	HandleMessage(msg *pb.Message) error
	SendMessage(msg *pb.Message) error
	To() (pb.PeerEndpoint, error)
	Stop() error
```


```
type RemoteLedgers interface {
    GetRemoteBlocks(peerID uint64, start, finish uint64) (<-chan *pb.SyncBlocks, error)
    GetRemoteStateSnapshot(peerID uint64) (<-chan *pb.SyncStateSnapshot, error)
    GetRemoteStateDeltas(peerID uint64, start, finish uint64) (<-chan *pb.SyncStateDeltas, error)
    }
    
RemoteLedgers 接口的存在主要是为了启用状态转移，和向其它副本询问区块链的状态。和WritableLedger接口一样，这不是给正常的操作使用，而是为追赶，错误恢复等操作而设计的。这个接口中的所有函数调用这都有责任来处理超时。这个接口包含下面这些函数：

GetRemoteBlocks(peerID uint64, start, finish uint64) (<-chan *pb.SyncBlocks, error)
这个函数尝试从由peerID指定的 peer 中取出由start和finish标识的范围中的*pb.SyncBlocks流。一般情况下，由于区块链必须是从结束到开始这样的顺序来验证的，所以start是比finish更高的块编号。由于慢速的结构，其它请求的返回可能出现在这个通道中，所以调用者必须验证返回的是期望的块。第二次以同样的peerID来调用这个方法会导致第一次的通道关闭。

GetRemoteStateSnapshot(peerID uint64) (<-chan *pb.SyncStateSnapshot, error)
这个函数尝试从由peerID指定的 peer 中取出*pb.SyncStateSnapshot流。为了应用结果，首先需要通过WritableLedger的EmptyState调用来清空存在在状态，然后顺序应用包含在流中的变化量。

 GetRemoteStateDeltas(peerID uint64, start, finish uint64) (<-chan *pb.SyncStateDeltas, error)
这个函数尝试从由peerID指定的 peer 中取出由start和finish标识的范围中的*pb.SyncStateDeltas流。由于慢速的结构，其它请求的返回可能出现在这个通道中，所以调用者必须验证返回的是期望的块变化量。第二次以同样的peerID来调用这个方法会导致第一次的通道关闭。
```
##### 3、Handler
```
// Handler peer handler implementation.
type Handler struct {
	chatMutex                     sync.Mutex
	ToPeerEndpoint                *pb.PeerEndpoint
	Coordinator                   MessageHandlerCoordinator
	ChatStream                    ChatStream
	doneChan                      chan struct{}
	FSM                           *fsm.FSM
	initiatedStream               bool // Was the stream initiated within this Peer
	registered                    bool
	syncBlocks                    chan *pb.SyncBlocks
	snapshotRequestHandler        *syncStateSnapshotRequestHandler
	syncStateDeltasRequestHandler *syncStateDeltasHandler
	syncBlocksRequestHandler      *syncBlocksRequestHandler
}
```
Handler 实现MessageHandler通信的方法有：
```
func (d *Handler) To() (pb.PeerEndpoint, error) {
func (d *Handler) Stop() error {
func (d *Handler) HandleMessage(msg *pb.Message) error {
func (d *Handler) SendMessage(msg *pb.Message) error {
func (d *Handler) RequestBlocks(syncBlockRange *pb.SyncBlockRange) (<-chan *pb.SyncBlocks, error) {
func (d *Handler) RequestStateSnapshot() (<-chan *pb.SyncStateSnapshot, error) {
func (d *Handler) RequestStateDeltas(syncBlockRange *pb.SyncBlockRange) (<-chan *pb.SyncStateDeltas, error) {
```
Handler 自己的方法有：
```
func (d *Handler) enterState(e *fsm.Event) {
func (d *Handler) deregister() error {
func (d *Handler) beforeHello(e *fsm.Event) {
func (d *Handler) beforeGetPeers(e *fsm.Event) {
func (d *Handler) beforePeers(e *fsm.Event) {
func (d *Handler) beforeBlockAdded(e *fsm.Event) {
func (d *Handler) when(stateToCheck string) bool {
func (d *Handler) start() error {
func (d *Handler) beforeSyncGetBlocks(e *fsm.Event) {
func (d *Handler) beforeSyncBlocks(e *fsm.Event) {
func (d *Handler) sendBlocks(syncBlockRange *pb.SyncBlockRange) {
func (d *Handler) beforeSyncStateGetSnapshot(e *fsm.Event) {
func (d *Handler) beforeSyncStateSnapshot(e *fsm.Event) {
func (d *Handler) sendStateSnapshot(syncStateSnapshotRequest *pb.SyncStateSnapshotRequest) {
func (d *Handler) beforeSyncStateGetDeltas(e *fsm.Event) {
func (d *Handler) sendStateDeltas(syncStateDeltasRequest *pb.SyncStateDeltasRequest) {
func (d *Handler) beforeSyncStateDeltas(e *fsm.Event) {
```

#####4、MessageHandlerCoordinator
```
// MessageHandlerCoordinator responsible for coordinating between the registered MessageHandler's
type MessageHandlerCoordinator interface {
	Peer                      分析完成
	SecurityAccessor          分析完成
	BlockChainAccessor        分析完成
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
```
func (p *PeerImpl) GetPeerEndpoint() (*pb.PeerEndpoint, error) {
func (p *PeerImpl) NewOpenchainDiscoveryHello() (*pb.Message, error) {

type SecurityAccessor interface {
	GetSecHelper() crypto.Peer
}
GetSecHelper 是安全相关，获取加密相关方法
```
```
// BlockChainAccessor interface for retreiving blocks by block number
type BlockChainAccessor interface {
	GetBlockByNumber(blockNumber uint64) (*pb.Block, error)
	GetBlockchainSize() uint64
	GetCurrentStateHash() (stateHash []byte, err error)
}

账本相关的方法：
func (p *PeerImpl) GetBlockByNumber(blockNumber uint64) (*pb.Block, error) {
func (p *PeerImpl) GetBlockchainSize() uint64 {
func (p *PeerImpl) GetCurrentStateHash() (stateHash []byte, err error) {

账本相关的方法，与网络无关，不详细说明：
// BlockChainModifier interface for applying changes to the block chain
type BlockChainModifier interface {
	ApplyStateDelta(id interface{}, delta *statemgmt.StateDelta) error
	RollbackStateDelta(id interface{}) error
	CommitStateDelta(id interface{}) error
	EmptyState() error
	PutBlock(blockNumber uint64, block *pb.Block) error
}

账本相关的方法，与网络无关，不详细说明：
// BlockChainUtil interface for interrogating the block chain
type BlockChainUtil interface {
	HashBlock(block *pb.Block) ([]byte, error)
	VerifyBlockchain(start, finish uint64) (uint64, error)
}
```
```

与网络请求有关，但不是核心消息。
// StateAccessor interface for retreiving blocks by block number
type StateAccessor interface {
	GetStateSnapshot() (*state.StateSnapshot, error)
	GetStateDelta(blockNumber uint64) (*statemgmt.StateDelta, error)
}

将 MessageHandler 写入到 PeerImpl 的 HandlerMap 中
// RegisterHandler register a MessageHandler with this coordinator
func (p *PeerImpl) RegisterHandler(messageHandler MessageHandler) error {

// DeregisterHandler deregisters an already registered MessageHandler for this coordinator
func (p *PeerImpl) DeregisterHandler(messageHandler MessageHandler) error {


func (p *PeerImpl) Broadcast(msg *pb.Message, typ pb.PeerEndpoint_Type) []error {
// Unicast sends a message to a specific peer.
func (p *PeerImpl) Unicast(msg *pb.Message, receiverHandle *pb.PeerID) error {
func (p *PeerImpl) GetPeers() (*pb.PeersMessage, error) {

func (p *PeerImpl) GetRemoteLedger(receiverHandle *pb.PeerID) (RemoteLedger, error) {

RemoteLedger自带三个接口：

type RemoteLedgers interface {
    GetRemoteBlocks(peerID uint64, start, finish uint64) (<-chan *pb.SyncBlocks, error)
    GetRemoteStateSnapshot(peerID uint64) (<-chan *pb.SyncStateSnapshot, error)
    GetRemoteStateDeltas(peerID uint64, start, finish uint64) (<-chan *pb.SyncStateDeltas, error)
    }


func (p *PeerImpl) PeersDiscovered(peersMessage *pb.PeersMessage) error {
```
```
发送交易消息：
func (p *PeerImpl) ExecuteTransaction(transaction *pb.Transaction) (response *pb.Response) {

//ExecuteTransaction executes transactions decides to do execute in dev or prod mode
func (p *PeerImpl) ExecuteTransaction(transaction *pb.Transaction) (response *pb.Response) {
	if p.isValidator {
		response = p.sendTransactionsToLocalEngine(transaction)
	} else {
		peerAddresses := p.discHelper.GetRandomNodes(1)
		response = p.SendTransactionsToPeer(peerAddresses[0], transaction)
	}
	return response
}
```
最后discovery 是自己的一个发现工具包。

真正的通信入口是`peer/handler.go`下main的
```
helloMessage, err := d.Coordinator.NewOpenchainDiscoveryHello()
```

#####5、PeerImpl
PeerImpl 实现MessageHandlerCoordinator 的函数,，不一一举例。

PeerImpl 自己实现的函数及一些工具函数：
```

```