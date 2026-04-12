//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

// Stub state — global because the cbinding layer has no per-handle storage.
// Tests that mutate these MUST reset them in cleanup.

var (
	stubP2PNotifyCounter uint64
	stubP2PMu            sync.Mutex
	stubP2PQueue         []EOS_P2P_ReceivePacketResult
	stubP2PNATType              = EOS_NAT_Open
	stubP2PRelayControl         = EOS_RC_AllowRelays
	stubP2PPort          uint16 = 7777
	stubP2PMaxAdditional uint16 = 99
)

// StubP2PReset clears all stub-side P2P state. Tests that touch the global
// state SHOULD call this in their cleanup.
func StubP2PReset() {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	stubP2PQueue = nil
	stubP2PNATType = EOS_NAT_Open
	stubP2PRelayControl = EOS_RC_AllowRelays
	stubP2PPort = 7777
	stubP2PMaxAdditional = 99
}

// StubP2PSetNATType lets tests force a specific NAT type for both
// QueryNATType and GetNATType.
func StubP2PSetNATType(nat EOS_ENATType) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	stubP2PNATType = nat
}

func EOS_P2P_SendPacket(handle EOS_HP2P, opts *EOS_P2P_SendPacketOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	dataCopy := make([]byte, len(opts.Data))
	copy(dataCopy, opts.Data)

	stubP2PMu.Lock()
	stubP2PQueue = append(stubP2PQueue, EOS_P2P_ReceivePacketResult{
		PeerId:   opts.LocalUserId,
		SocketId: EOS_P2P_SocketId{Name: socketName},
		Channel:  opts.Channel,
		Data:     dataCopy,
	})
	stubP2PMu.Unlock()
	return EOS_EResult_Success
}

func EOS_P2P_GetNextReceivedPacketSize(handle EOS_HP2P, localUserId EOS_ProductUserId, channel *uint8) (uint32, EOS_EResult) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	for i := range stubP2PQueue {
		if channel == nil || stubP2PQueue[i].Channel == *channel {
			return uint32(len(stubP2PQueue[i].Data)), EOS_EResult_Success
		}
	}
	return 0, EOS_EResult_NotFound
}

func EOS_P2P_ReceivePacket(handle EOS_HP2P, localUserId EOS_ProductUserId, maxBytes uint32, channel *uint8) (*EOS_P2P_ReceivePacketResult, EOS_EResult) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	for i := range stubP2PQueue {
		if channel != nil && stubP2PQueue[i].Channel != *channel {
			continue
		}
		pkt := stubP2PQueue[i]
		stubP2PQueue = append(stubP2PQueue[:i], stubP2PQueue[i+1:]...)
		return &pkt, EOS_EResult_Success
	}
	return nil, EOS_EResult_NotFound
}

func EOS_P2P_AcceptConnection(handle EOS_HP2P, opts *EOS_P2P_AcceptConnectionOptions) EOS_EResult {
	return EOS_EResult_Success
}

func EOS_P2P_CloseConnection(handle EOS_HP2P, opts *EOS_P2P_CloseConnectionOptions) EOS_EResult {
	return EOS_EResult_Success
}

func EOS_P2P_CloseConnections(handle EOS_HP2P, opts *EOS_P2P_CloseConnectionsOptions) EOS_EResult {
	return EOS_EResult_Success
}

func EOS_P2P_QueryNATType(handle EOS_HP2P, clientData uintptr) {
	stubP2PMu.Lock()
	natType := stubP2PNATType
	stubP2PMu.Unlock()
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_P2P_OnQueryNATTypeCompleteInfo{
				ResultCode: EOS_EResult_Success,
				NATType:    natType,
			},
		})
	}()
}

func EOS_P2P_GetNATType(handle EOS_HP2P) (EOS_ENATType, EOS_EResult) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	return stubP2PNATType, EOS_EResult_Success
}

func EOS_P2P_SetRelayControl(handle EOS_HP2P, opts *EOS_P2P_SetRelayControlOptions) EOS_EResult {
	stubP2PMu.Lock()
	stubP2PRelayControl = opts.RelayControl
	stubP2PMu.Unlock()
	return EOS_EResult_Success
}

func EOS_P2P_GetRelayControl(handle EOS_HP2P) (EOS_ERelayControl, EOS_EResult) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	return stubP2PRelayControl, EOS_EResult_Success
}

func EOS_P2P_SetPortRange(handle EOS_HP2P, opts *EOS_P2P_SetPortRangeOptions) EOS_EResult {
	stubP2PMu.Lock()
	stubP2PPort = opts.Port
	stubP2PMaxAdditional = opts.MaxAdditionalPortsToTry
	stubP2PMu.Unlock()
	return EOS_EResult_Success
}

func EOS_P2P_GetPortRange(handle EOS_HP2P) (uint16, uint16, EOS_EResult) {
	stubP2PMu.Lock()
	defer stubP2PMu.Unlock()
	return stubP2PPort, stubP2PMaxAdditional, EOS_EResult_Success
}

func EOS_P2P_SetPacketQueueSize(handle EOS_HP2P, opts *EOS_P2P_SetPacketQueueSizeOptions) EOS_EResult {
	return EOS_EResult_Success
}

func EOS_P2P_GetPacketQueueInfo(handle EOS_HP2P) (*EOS_P2P_PacketQueueInfo, EOS_EResult) {
	return &EOS_P2P_PacketQueueInfo{}, EOS_EResult_Success
}

func EOS_P2P_ClearPacketQueue(handle EOS_HP2P, opts *EOS_P2P_ClearPacketQueueOptions) EOS_EResult {
	stubP2PMu.Lock()
	stubP2PQueue = nil
	stubP2PMu.Unlock()
	return EOS_EResult_Success
}

func EOS_P2P_AddNotifyPeerConnectionRequest(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubP2PNotifyCounter, 1))
}

func EOS_P2P_RemoveNotifyPeerConnectionRequest(handle EOS_HP2P, id EOS_NotificationId) {}

func EOS_P2P_AddNotifyPeerConnectionEstablished(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubP2PNotifyCounter, 1))
}

func EOS_P2P_RemoveNotifyPeerConnectionEstablished(handle EOS_HP2P, id EOS_NotificationId) {}

func EOS_P2P_AddNotifyPeerConnectionInterrupted(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubP2PNotifyCounter, 1))
}

func EOS_P2P_RemoveNotifyPeerConnectionInterrupted(handle EOS_HP2P, id EOS_NotificationId) {}

func EOS_P2P_AddNotifyPeerConnectionClosed(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubP2PNotifyCounter, 1))
}

func EOS_P2P_RemoveNotifyPeerConnectionClosed(handle EOS_HP2P, id EOS_NotificationId) {}

func EOS_P2P_AddNotifyIncomingPacketQueueFull(handle EOS_HP2P, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubP2PNotifyCounter, 1))
}

func EOS_P2P_RemoveNotifyIncomingPacketQueueFull(handle EOS_HP2P, id EOS_NotificationId) {}
