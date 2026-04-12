package cbinding

const (
	EOS_P2P_SOCKETID_SOCKETNAME_SIZE = 33
	EOS_P2P_MAX_PACKET_SIZE          = 1170
)

// EOS_P2P_SocketId mirrors the SDK socket-id struct. Name is 1-32 alphanumeric
// chars; the C wrapper copies it into the SDK's fixed 33-byte array and sets
// ApiVersion (omitted here per the cbinding convention — version macros live
// in the C wrapper, see MEM-2/MEM-4).
type EOS_P2P_SocketId struct {
	Name string
}

type EOS_P2P_SendPacketOptions struct {
	LocalUserId                 EOS_ProductUserId
	RemoteUserId                EOS_ProductUserId
	SocketId                    *EOS_P2P_SocketId
	Channel                     uint8
	Data                        []byte
	AllowDelayedDelivery        bool
	Reliability                 EOS_EPacketReliability
	DisableAutoAcceptConnection bool
}

type EOS_P2P_ReceivePacketResult struct {
	PeerId   EOS_ProductUserId
	SocketId EOS_P2P_SocketId
	Channel  uint8
	Data     []byte
}

type EOS_P2P_AcceptConnectionOptions struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     *EOS_P2P_SocketId
}

type EOS_P2P_CloseConnectionOptions struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     *EOS_P2P_SocketId
}

type EOS_P2P_CloseConnectionsOptions struct {
	LocalUserId EOS_ProductUserId
	SocketId    *EOS_P2P_SocketId
}

type EOS_P2P_SetPortRangeOptions struct {
	Port                    uint16
	MaxAdditionalPortsToTry uint16
}

type EOS_P2P_SetRelayControlOptions struct {
	RelayControl EOS_ERelayControl
}

type EOS_P2P_SetPacketQueueSizeOptions struct {
	IncomingPacketQueueMaxSizeBytes uint64
	OutgoingPacketQueueMaxSizeBytes uint64
}

type EOS_P2P_PacketQueueInfo struct {
	IncomingPacketQueueMaxSizeBytes       uint64
	IncomingPacketQueueCurrentSizeBytes   uint64
	IncomingPacketQueueCurrentPacketCount uint64
	OutgoingPacketQueueMaxSizeBytes       uint64
	OutgoingPacketQueueCurrentSizeBytes   uint64
	OutgoingPacketQueueCurrentPacketCount uint64
}

// Notification info structs (received by Go callbacks via trampolines).

type EOS_P2P_OnIncomingConnectionRequestInfo struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     EOS_P2P_SocketId
}

type EOS_P2P_OnPeerConnectionEstablishedInfo struct {
	LocalUserId    EOS_ProductUserId
	RemoteUserId   EOS_ProductUserId
	SocketId       EOS_P2P_SocketId
	ConnectionType EOS_EConnectionEstablishedType
	NetworkType    EOS_ENetworkConnectionType
}

type EOS_P2P_OnRemoteConnectionClosedInfo struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     EOS_P2P_SocketId
	Reason       EOS_EConnectionClosedReason
}

type EOS_P2P_OnPeerConnectionInterruptedInfo struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     EOS_P2P_SocketId
}

type EOS_P2P_OnIncomingPacketQueueFullInfo struct {
	PacketQueueMaxSizeBytes     uint64
	PacketQueueCurrentSizeBytes uint64
	OverflowPacketLocalUserId   EOS_ProductUserId
	OverflowPacketChannel       uint8
	OverflowPacketSizeBytes     uint32
}

type EOS_P2P_OnQueryNATTypeCompleteInfo struct {
	ResultCode EOS_EResult
	NATType    EOS_ENATType
}

type EOS_P2P_ClearPacketQueueOptions struct {
	LocalUserId  EOS_ProductUserId
	RemoteUserId EOS_ProductUserId
	SocketId     *EOS_P2P_SocketId
}
