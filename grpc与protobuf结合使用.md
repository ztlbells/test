## grpc �� protobuf ��ϰ���
#### 1������
protocolbuffer(���¼��PB)��google ��һ�����ݽ����ĸ�ʽ�������������ԣ�������ƽ̨��google �ṩ�˶������Ե�ʵ�֣�java��c#��c++��go �� python��ÿһ��ʵ�ֶ���������Ӧ���Եı������Լ����ļ�����������һ�ֶ����Ƶĸ�ʽ����ʹ�� xml �������ݽ�������ࡣ���԰������ڷֲ�ʽӦ��֮�������ͨ�Ż����칹�����µ����ݽ�������Ϊһ��Ч�ʺͼ����Զ�������Ķ��������ݴ����ʽ�����������������紫�䡢�����ļ������ݴ洢���������

gRPC��һ�������ܡ�ͨ�õĿ�ԴRPC��ܣ�����Google��Ҫ�����ƶ�Ӧ�ÿ���������HTTP/2Э���׼����ƣ�����ProtoBuf(Protocol Buffers)���л�Э�鿪������֧���ڶ࿪�����ԡ�gRPC�ṩ��һ�ּ򵥵ķ�������ȷ�ض�������ΪiOS��Android�ͺ�̨֧�ַ����Զ����ɿɿ��Ժ�ǿ�Ŀͻ��˹��ܿ⡣

##### 2��protobuf grpc http2 stream


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

PeerServer�˵�ʵ�֣�
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

PeerClient�˵ĵ��ã�
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


PeerImpl����ʵ�������е� MessageHandlerCoordinator ����������
```
RegisterHandler(messageHandler MessageHandler) error
DeregisterHandler(messageHandler MessageHandler) error
```

�� peer.go ����һ��struct
```
type ChatStream interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}
```
ÿ��ʹ�õ�ʱ���Ǳ������ȥ�ģ�stream ��ֵ�� `pb.Peer_ChatServer`��

##### 3�������ص�����Ĵ��룺

```
    peer.engine, err = engFactory(peer)
	if err != nil {
		return nil, err
	}
	peer.handlerFactory = peer.engine.GetHandlerFactory()
```

messageHandler��MessageHandlerCoordinator��handler��PeerImpl��ConsensusHandler��
MessageHandlerCoordinator ��һ�׷������ϡ�
messageHandler Ҳ��һ�׷������ϡ�

handler��ConsensusHandler��PeerImpl ���ǽṹ�塣
���ǵĶ������£�
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

�� peer.go ��,����ࣺ
```
func (p *PeerImpl) GetPeers() (*pb.PeersMessage, error) {
```
�ɼ� MessageHandlerCoordinator �ķ���ȫ���� PeerImpl ʵ�֡�MessageHandlerCoordinator �ṩ����෽�����Բ����޸� PeerImpl �����ݽṹ��

MessageHandler ���壺
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

MessageHandler Ҳ��һ�׷������� ���� handler ȫ��ʵ�֣�����ĿǰConsensusHandler ����ʵ�֡������`ConsensusHandler`������Ϊ�ں�����`MessageHandler`����������壩
handler ���壺
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
ChatStream ���� pb �������������� send��resv��

ConsensusHandler��
```
// ConsensusHandler handles consensus messages.
// It also implements the Stack.
type ConsensusHandler struct {
	peer.MessageHandler
	consenterChan chan *util.Message
	coordinator   peer.MessageHandlerCoordinator
}
```


##### 4����Ϣͨ�ź�������Ϣ������������ʹ��FSM��
���е���Ϣͨ�ź������£�
#####1��chatstream
```
// ChatStream interface supported by stream between Peers
type ChatStream interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}
```
`ChatStream` ���� pb.chatServer �������������� send��resv��

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
�������ߵ�`send`��`Recv`�ڲ�ͬ��ʱ������ͬ�ĺ��������á�

#####2��MessageHandler
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
    
RemoteLedgers �ӿڵĴ�����Ҫ��Ϊ������״̬ת�ƣ�������������ѯ����������״̬����WritableLedger�ӿ�һ�����ⲻ�Ǹ������Ĳ���ʹ�ã�����Ϊ׷�ϣ�����ָ��Ȳ�������Ƶġ�����ӿ��е����к��������ⶼ������������ʱ������ӿڰ���������Щ������

GetRemoteBlocks(peerID uint64, start, finish uint64) (<-chan *pb.SyncBlocks, error)
����������Դ���peerIDָ���� peer ��ȡ����start��finish��ʶ�ķ�Χ�е�*pb.SyncBlocks����һ������£����������������Ǵӽ�������ʼ������˳������֤�ģ�����start�Ǳ�finish���ߵĿ��š��������ٵĽṹ����������ķ��ؿ��ܳ��������ͨ���У����Ե����߱�����֤���ص��������Ŀ顣�ڶ�����ͬ����peerID��������������ᵼ�µ�һ�ε�ͨ���رա�

GetRemoteStateSnapshot(peerID uint64) (<-chan *pb.SyncStateSnapshot, error)
����������Դ���peerIDָ���� peer ��ȡ��*pb.SyncStateSnapshot����Ϊ��Ӧ�ý����������Ҫͨ��WritableLedger��EmptyState��������մ�����״̬��Ȼ��˳��Ӧ�ð��������еı仯����

 GetRemoteStateDeltas(peerID uint64, start, finish uint64) (<-chan *pb.SyncStateDeltas, error)
����������Դ���peerIDָ���� peer ��ȡ����start��finish��ʶ�ķ�Χ�е�*pb.SyncStateDeltas�����������ٵĽṹ����������ķ��ؿ��ܳ��������ͨ���У����Ե����߱�����֤���ص��������Ŀ�仯�����ڶ�����ͬ����peerID��������������ᵼ�µ�һ�ε�ͨ���رա�
```
##### 3��Handler
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
Handler ʵ��MessageHandlerͨ�ŵķ����У�
```
func (d *Handler) To() (pb.PeerEndpoint, error) {
func (d *Handler) Stop() error {
func (d *Handler) HandleMessage(msg *pb.Message) error {
func (d *Handler) SendMessage(msg *pb.Message) error {
func (d *Handler) RequestBlocks(syncBlockRange *pb.SyncBlockRange) (<-chan *pb.SyncBlocks, error) {
func (d *Handler) RequestStateSnapshot() (<-chan *pb.SyncStateSnapshot, error) {
func (d *Handler) RequestStateDeltas(syncBlockRange *pb.SyncBlockRange) (<-chan *pb.SyncStateDeltas, error) {
```
Handler �Լ��ķ����У�
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

#####4��MessageHandlerCoordinator
```
// MessageHandlerCoordinator responsible for coordinating between the registered MessageHandler's
type MessageHandlerCoordinator interface {
	Peer                      �������
	SecurityAccessor          �������
	BlockChainAccessor        �������
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
GetSecHelper �ǰ�ȫ��أ���ȡ������ط���
```
```
// BlockChainAccessor interface for retreiving blocks by block number
type BlockChainAccessor interface {
	GetBlockByNumber(blockNumber uint64) (*pb.Block, error)
	GetBlockchainSize() uint64
	GetCurrentStateHash() (stateHash []byte, err error)
}

�˱���صķ�����
func (p *PeerImpl) GetBlockByNumber(blockNumber uint64) (*pb.Block, error) {
func (p *PeerImpl) GetBlockchainSize() uint64 {
func (p *PeerImpl) GetCurrentStateHash() (stateHash []byte, err error) {

�˱���صķ������������޹أ�����ϸ˵����
// BlockChainModifier interface for applying changes to the block chain
type BlockChainModifier interface {
	ApplyStateDelta(id interface{}, delta *statemgmt.StateDelta) error
	RollbackStateDelta(id interface{}) error
	CommitStateDelta(id interface{}) error
	EmptyState() error
	PutBlock(blockNumber uint64, block *pb.Block) error
}

�˱���صķ������������޹أ�����ϸ˵����
// BlockChainUtil interface for interrogating the block chain
type BlockChainUtil interface {
	HashBlock(block *pb.Block) ([]byte, error)
	VerifyBlockchain(start, finish uint64) (uint64, error)
}
```
```

�����������йأ������Ǻ�����Ϣ��
// StateAccessor interface for retreiving blocks by block number
type StateAccessor interface {
	GetStateSnapshot() (*state.StateSnapshot, error)
	GetStateDelta(blockNumber uint64) (*statemgmt.StateDelta, error)
}

�� MessageHandler д�뵽 PeerImpl �� HandlerMap ��
// RegisterHandler register a MessageHandler with this coordinator
func (p *PeerImpl) RegisterHandler(messageHandler MessageHandler) error {

// DeregisterHandler deregisters an already registered MessageHandler for this coordinator
func (p *PeerImpl) DeregisterHandler(messageHandler MessageHandler) error {


func (p *PeerImpl) Broadcast(msg *pb.Message, typ pb.PeerEndpoint_Type) []error {
// Unicast sends a message to a specific peer.
func (p *PeerImpl) Unicast(msg *pb.Message, receiverHandle *pb.PeerID) error {
func (p *PeerImpl) GetPeers() (*pb.PeersMessage, error) {

func (p *PeerImpl) GetRemoteLedger(receiverHandle *pb.PeerID) (RemoteLedger, error) {

RemoteLedger�Դ������ӿڣ�

type RemoteLedgers interface {
    GetRemoteBlocks(peerID uint64, start, finish uint64) (<-chan *pb.SyncBlocks, error)
    GetRemoteStateSnapshot(peerID uint64) (<-chan *pb.SyncStateSnapshot, error)
    GetRemoteStateDeltas(peerID uint64, start, finish uint64) (<-chan *pb.SyncStateDeltas, error)
    }


func (p *PeerImpl) PeersDiscovered(peersMessage *pb.PeersMessage) error {
```
```
���ͽ�����Ϣ��
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
���discovery ���Լ���һ�����ֹ��߰���

������ͨ�������`peer/handler.go`��main��
```
helloMessage, err := d.Coordinator.NewOpenchainDiscoveryHello()
```

#####5��PeerImpl
PeerImpl ʵ��MessageHandlerCoordinator �ĺ���,����һһ������

PeerImpl �Լ�ʵ�ֵĺ�����һЩ���ߺ�����
```

```