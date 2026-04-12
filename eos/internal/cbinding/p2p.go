//go:build !eosstub

package cbinding

/*
#include "p2p_wrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// withCSocketName invokes fn with a *C.char allocated from the Go socket name
// (or nil for an empty/optional name). Caller-frees lifetime is bounded by fn.
func withCSocketName(name string, fn func(*C.char)) {
	if name == "" {
		fn(nil)
		return
	}
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))
	fn(c)
}

func EOS_P2P_SendPacket(handle EOS_HP2P, opts *EOS_P2P_SendPacketOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	cName := C.CString(socketName)
	defer C.free(unsafe.Pointer(cName))

	var dataPtr unsafe.Pointer
	if len(opts.Data) > 0 {
		dataPtr = unsafe.Pointer(&opts.Data[0])
	}

	allowDelayed := C.int(0)
	if opts.AllowDelayedDelivery {
		allowDelayed = 1
	}
	disableAutoAccept := C.int(0)
	if opts.DisableAutoAcceptConnection {
		disableAutoAccept = 1
	}

	return EOS_EResult(C.eos_p2p_send_packet(
		C.uintptr_t(handle),
		C.uintptr_t(opts.LocalUserId),
		C.uintptr_t(opts.RemoteUserId),
		cName,
		C.uint8_t(opts.Channel),
		dataPtr,
		C.uint32_t(len(opts.Data)),
		allowDelayed,
		C.int(opts.Reliability),
		disableAutoAccept,
	))
}

// EOS_P2P_GetNextReceivedPacketSize returns the size of the next pending
// packet for localUserId. If channel is non-nil, only packets on that channel
// are considered.
func EOS_P2P_GetNextReceivedPacketSize(handle EOS_HP2P, localUserId EOS_ProductUserId, channel *uint8) (uint32, EOS_EResult) {
	var hasChannel C.int
	var ch C.uint8_t
	if channel != nil {
		hasChannel = 1
		ch = C.uint8_t(*channel)
	}
	var size C.uint32_t
	result := EOS_EResult(C.eos_p2p_get_next_received_packet_size(
		C.uintptr_t(handle), C.uintptr_t(localUserId), hasChannel, ch, &size))
	return uint32(size), result
}

// EOS_P2P_ReceivePacket reads the next pending packet for localUserId into a
// fresh Go-allocated buffer of size maxBytes. The single SDK→Go copy is
// intentional — see MEM-5 in docs/prd.md and the package note in m3.md.
func EOS_P2P_ReceivePacket(handle EOS_HP2P, localUserId EOS_ProductUserId, maxBytes uint32, channel *uint8) (*EOS_P2P_ReceivePacketResult, EOS_EResult) {
	var hasChannel C.int
	var ch C.uint8_t
	if channel != nil {
		hasChannel = 1
		ch = C.uint8_t(*channel)
	}

	data := make([]byte, maxBytes)
	socketName := make([]byte, EOS_P2P_SOCKETID_SOCKETNAME_SIZE)
	var peerId C.uintptr_t
	var outChannel C.uint8_t
	var bytesWritten C.uint32_t

	var dataPtr unsafe.Pointer
	if maxBytes > 0 {
		dataPtr = unsafe.Pointer(&data[0])
	}

	result := EOS_EResult(C.eos_p2p_receive_packet(
		C.uintptr_t(handle),
		C.uintptr_t(localUserId),
		C.uint32_t(maxBytes),
		hasChannel, ch,
		dataPtr,
		&peerId,
		(*C.char)(unsafe.Pointer(&socketName[0])),
		&outChannel,
		&bytesWritten,
	))
	if result != EOS_EResult_Success {
		return nil, result
	}

	return &EOS_P2P_ReceivePacketResult{
		PeerId:   EOS_ProductUserId(peerId),
		SocketId: EOS_P2P_SocketId{Name: cStringFromBytes(socketName)},
		Channel:  uint8(outChannel),
		Data:     data[:bytesWritten:bytesWritten],
	}, EOS_EResult_Success
}

// cStringFromBytes returns the contents of a NUL-terminated byte slice up to
// (but not including) the first NUL byte.
func cStringFromBytes(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func EOS_P2P_AcceptConnection(handle EOS_HP2P, opts *EOS_P2P_AcceptConnectionOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	var result EOS_EResult
	withCSocketName(socketName, func(c *C.char) {
		result = EOS_EResult(C.eos_p2p_accept_connection(
			C.uintptr_t(handle),
			C.uintptr_t(opts.LocalUserId),
			C.uintptr_t(opts.RemoteUserId),
			c,
		))
	})
	return result
}

func EOS_P2P_CloseConnection(handle EOS_HP2P, opts *EOS_P2P_CloseConnectionOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	var result EOS_EResult
	withCSocketName(socketName, func(c *C.char) {
		result = EOS_EResult(C.eos_p2p_close_connection(
			C.uintptr_t(handle),
			C.uintptr_t(opts.LocalUserId),
			C.uintptr_t(opts.RemoteUserId),
			c,
		))
	})
	return result
}

func EOS_P2P_CloseConnections(handle EOS_HP2P, opts *EOS_P2P_CloseConnectionsOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	var result EOS_EResult
	withCSocketName(socketName, func(c *C.char) {
		result = EOS_EResult(C.eos_p2p_close_connections(
			C.uintptr_t(handle),
			C.uintptr_t(opts.LocalUserId),
			c,
		))
	})
	return result
}

func EOS_P2P_QueryNATType(handle EOS_HP2P, clientData uintptr) {
	C.eos_p2p_query_nat_type(C.uintptr_t(handle), C.uintptr_t(clientData))
}

func EOS_P2P_GetNATType(handle EOS_HP2P) (EOS_ENATType, EOS_EResult) {
	var natType C.int
	result := EOS_EResult(C.eos_p2p_get_nat_type(C.uintptr_t(handle), &natType))
	return EOS_ENATType(natType), result
}

func EOS_P2P_SetRelayControl(handle EOS_HP2P, opts *EOS_P2P_SetRelayControlOptions) EOS_EResult {
	return EOS_EResult(C.eos_p2p_set_relay_control(C.uintptr_t(handle), C.int(opts.RelayControl)))
}

func EOS_P2P_GetRelayControl(handle EOS_HP2P) (EOS_ERelayControl, EOS_EResult) {
	var rc C.int
	result := EOS_EResult(C.eos_p2p_get_relay_control(C.uintptr_t(handle), &rc))
	return EOS_ERelayControl(rc), result
}

func EOS_P2P_SetPortRange(handle EOS_HP2P, opts *EOS_P2P_SetPortRangeOptions) EOS_EResult {
	return EOS_EResult(C.eos_p2p_set_port_range(
		C.uintptr_t(handle),
		C.uint16_t(opts.Port),
		C.uint16_t(opts.MaxAdditionalPortsToTry),
	))
}

func EOS_P2P_GetPortRange(handle EOS_HP2P) (uint16, uint16, EOS_EResult) {
	var port, maxAdditional C.uint16_t
	result := EOS_EResult(C.eos_p2p_get_port_range(C.uintptr_t(handle), &port, &maxAdditional))
	return uint16(port), uint16(maxAdditional), result
}

func EOS_P2P_SetPacketQueueSize(handle EOS_HP2P, opts *EOS_P2P_SetPacketQueueSizeOptions) EOS_EResult {
	return EOS_EResult(C.eos_p2p_set_packet_queue_size(
		C.uintptr_t(handle),
		C.uint64_t(opts.IncomingPacketQueueMaxSizeBytes),
		C.uint64_t(opts.OutgoingPacketQueueMaxSizeBytes),
	))
}

func EOS_P2P_GetPacketQueueInfo(handle EOS_HP2P) (*EOS_P2P_PacketQueueInfo, EOS_EResult) {
	var inMax, inCur, inPkts, outMax, outCur, outPkts C.uint64_t
	result := EOS_EResult(C.eos_p2p_get_packet_queue_info(
		C.uintptr_t(handle), &inMax, &inCur, &inPkts, &outMax, &outCur, &outPkts))
	if result != EOS_EResult_Success {
		return nil, result
	}
	return &EOS_P2P_PacketQueueInfo{
		IncomingPacketQueueMaxSizeBytes:       uint64(inMax),
		IncomingPacketQueueCurrentSizeBytes:   uint64(inCur),
		IncomingPacketQueueCurrentPacketCount: uint64(inPkts),
		OutgoingPacketQueueMaxSizeBytes:       uint64(outMax),
		OutgoingPacketQueueCurrentSizeBytes:   uint64(outCur),
		OutgoingPacketQueueCurrentPacketCount: uint64(outPkts),
	}, EOS_EResult_Success
}

func EOS_P2P_ClearPacketQueue(handle EOS_HP2P, opts *EOS_P2P_ClearPacketQueueOptions) EOS_EResult {
	socketName := ""
	if opts.SocketId != nil {
		socketName = opts.SocketId.Name
	}
	var result EOS_EResult
	withCSocketName(socketName, func(c *C.char) {
		result = EOS_EResult(C.eos_p2p_clear_packet_queue(
			C.uintptr_t(handle),
			C.uintptr_t(opts.LocalUserId),
			C.uintptr_t(opts.RemoteUserId),
			c,
		))
	})
	return result
}

// addNotifyWithSocket is a helper for the four AddNotifyPeer* notifications
// which all share the same shape: localUserId + optional socket filter +
// clientData → notificationId.
type addNotifyFunc func(handle, localUserId C.uintptr_t, hasSocket C.int, name *C.char, clientData C.uintptr_t) C.uint64_t

func addNotifyWithSocket(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr, fn addNotifyFunc) EOS_NotificationId {
	hasSocket := C.int(0)
	name := ""
	if socket != nil {
		hasSocket = 1
		name = socket.Name
	}
	var id C.uint64_t
	withCSocketName(name, func(c *C.char) {
		id = fn(C.uintptr_t(handle), C.uintptr_t(localUserId), hasSocket, c, C.uintptr_t(clientData))
	})
	return EOS_NotificationId(id)
}

func EOS_P2P_AddNotifyPeerConnectionRequest(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return addNotifyWithSocket(handle, localUserId, socket, clientData,
		func(h, lu C.uintptr_t, hs C.int, n *C.char, cd C.uintptr_t) C.uint64_t {
			return C.eos_p2p_add_notify_peer_connection_request(h, lu, hs, n, cd)
		})
}

func EOS_P2P_RemoveNotifyPeerConnectionRequest(handle EOS_HP2P, id EOS_NotificationId) {
	C.eos_p2p_remove_notify_peer_connection_request(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_P2P_AddNotifyPeerConnectionEstablished(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return addNotifyWithSocket(handle, localUserId, socket, clientData,
		func(h, lu C.uintptr_t, hs C.int, n *C.char, cd C.uintptr_t) C.uint64_t {
			return C.eos_p2p_add_notify_peer_connection_established(h, lu, hs, n, cd)
		})
}

func EOS_P2P_RemoveNotifyPeerConnectionEstablished(handle EOS_HP2P, id EOS_NotificationId) {
	C.eos_p2p_remove_notify_peer_connection_established(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_P2P_AddNotifyPeerConnectionInterrupted(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return addNotifyWithSocket(handle, localUserId, socket, clientData,
		func(h, lu C.uintptr_t, hs C.int, n *C.char, cd C.uintptr_t) C.uint64_t {
			return C.eos_p2p_add_notify_peer_connection_interrupted(h, lu, hs, n, cd)
		})
}

func EOS_P2P_RemoveNotifyPeerConnectionInterrupted(handle EOS_HP2P, id EOS_NotificationId) {
	C.eos_p2p_remove_notify_peer_connection_interrupted(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_P2P_AddNotifyPeerConnectionClosed(handle EOS_HP2P, localUserId EOS_ProductUserId, socket *EOS_P2P_SocketId, clientData uintptr) EOS_NotificationId {
	return addNotifyWithSocket(handle, localUserId, socket, clientData,
		func(h, lu C.uintptr_t, hs C.int, n *C.char, cd C.uintptr_t) C.uint64_t {
			return C.eos_p2p_add_notify_peer_connection_closed(h, lu, hs, n, cd)
		})
}

func EOS_P2P_RemoveNotifyPeerConnectionClosed(handle EOS_HP2P, id EOS_NotificationId) {
	C.eos_p2p_remove_notify_peer_connection_closed(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_P2P_AddNotifyIncomingPacketQueueFull(handle EOS_HP2P, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_p2p_add_notify_incoming_packet_queue_full(
		C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_P2P_RemoveNotifyIncomingPacketQueueFull(handle EOS_HP2P, id EOS_NotificationId) {
	C.eos_p2p_remove_notify_incoming_packet_queue_full(C.uintptr_t(handle), C.uint64_t(id))
}
