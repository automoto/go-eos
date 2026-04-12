//go:build !eosstub

package cbinding

/*
#include <stdint.h>
*/
import "C"

import (
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
)

//export goP2PQueryNATTypeCallback
func goP2PQueryNATTypeCallback(resultCode C.int, clientData C.uintptr_t, natType C.int) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_P2P_OnQueryNATTypeCompleteInfo{
			ResultCode: EOS_EResult(resultCode),
			NATType:    EOS_ENATType(natType),
		},
	})
}

//export goP2POnIncomingConnectionRequest
func goP2POnIncomingConnectionRequest(clientData C.uintptr_t, localUserId C.uintptr_t,
	remoteUserId C.uintptr_t, socketName *C.char) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_P2P_OnIncomingConnectionRequestInfo{
		LocalUserId:  EOS_ProductUserId(localUserId),
		RemoteUserId: EOS_ProductUserId(remoteUserId),
		SocketId:     EOS_P2P_SocketId{Name: C.GoString(socketName)},
	})
}

//export goP2POnPeerConnectionEstablished
func goP2POnPeerConnectionEstablished(clientData C.uintptr_t, localUserId C.uintptr_t,
	remoteUserId C.uintptr_t, socketName *C.char, connectionType C.int, networkType C.int) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_P2P_OnPeerConnectionEstablishedInfo{
		LocalUserId:    EOS_ProductUserId(localUserId),
		RemoteUserId:   EOS_ProductUserId(remoteUserId),
		SocketId:       EOS_P2P_SocketId{Name: C.GoString(socketName)},
		ConnectionType: EOS_EConnectionEstablishedType(connectionType),
		NetworkType:    EOS_ENetworkConnectionType(networkType),
	})
}

//export goP2POnPeerConnectionInterrupted
func goP2POnPeerConnectionInterrupted(clientData C.uintptr_t, localUserId C.uintptr_t,
	remoteUserId C.uintptr_t, socketName *C.char) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_P2P_OnPeerConnectionInterruptedInfo{
		LocalUserId:  EOS_ProductUserId(localUserId),
		RemoteUserId: EOS_ProductUserId(remoteUserId),
		SocketId:     EOS_P2P_SocketId{Name: C.GoString(socketName)},
	})
}

//export goP2POnPeerConnectionClosed
func goP2POnPeerConnectionClosed(clientData C.uintptr_t, localUserId C.uintptr_t,
	remoteUserId C.uintptr_t, socketName *C.char, reason C.int) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_P2P_OnRemoteConnectionClosedInfo{
		LocalUserId:  EOS_ProductUserId(localUserId),
		RemoteUserId: EOS_ProductUserId(remoteUserId),
		SocketId:     EOS_P2P_SocketId{Name: C.GoString(socketName)},
		Reason:       EOS_EConnectionClosedReason(reason),
	})
}

//export goP2POnIncomingPacketQueueFull
func goP2POnIncomingPacketQueueFull(clientData C.uintptr_t, maxSizeBytes C.uint64_t,
	currentSizeBytes C.uint64_t, localUserId C.uintptr_t, channel C.uint8_t,
	packetSizeBytes C.uint32_t) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_P2P_OnIncomingPacketQueueFullInfo{
		PacketQueueMaxSizeBytes:     uint64(maxSizeBytes),
		PacketQueueCurrentSizeBytes: uint64(currentSizeBytes),
		OverflowPacketLocalUserId:   EOS_ProductUserId(localUserId),
		OverflowPacketChannel:       uint8(channel),
		OverflowPacketSizeBytes:     uint32(packetSizeBytes),
	})
}
