// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_p2p.h"
#include <stdint.h>
#include <string.h>

/* Forward declarations for Go export functions (implemented in p2p_callback.go) */
extern void goP2PQueryNATTypeCallback(int resultCode, uintptr_t clientData, int natType);
extern void goP2POnIncomingConnectionRequest(uintptr_t clientData, uintptr_t localUserId,
											 uintptr_t remoteUserId, const char* socketName);
extern void goP2POnPeerConnectionEstablished(uintptr_t clientData, uintptr_t localUserId,
											 uintptr_t remoteUserId, const char* socketName,
											 int connectionType, int networkType);
extern void goP2POnPeerConnectionInterrupted(uintptr_t clientData, uintptr_t localUserId,
											 uintptr_t remoteUserId, const char* socketName);
extern void goP2POnPeerConnectionClosed(uintptr_t clientData, uintptr_t localUserId,
										uintptr_t remoteUserId, const char* socketName, int reason);
extern void goP2POnIncomingPacketQueueFull(uintptr_t clientData, uint64_t maxSizeBytes,
										   uint64_t currentSizeBytes, uintptr_t localUserId,
										   uint8_t channel, uint32_t packetSizeBytes);

/* Trampolines */

static void p2pQueryNATTypeTrampoline(const EOS_P2P_OnQueryNATTypeCompleteInfo* data) {
	goP2PQueryNATTypeCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
							  (int)data->NATType);
}

static void
p2pIncomingConnectionRequestTrampoline(const EOS_P2P_OnIncomingConnectionRequestInfo* data) {
	goP2POnIncomingConnectionRequest((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
									 (uintptr_t)data->RemoteUserId,
									 data->SocketId ? data->SocketId->SocketName : "");
}

static void
p2pPeerConnectionEstablishedTrampoline(const EOS_P2P_OnPeerConnectionEstablishedInfo* data) {
	goP2POnPeerConnectionEstablished((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
									 (uintptr_t)data->RemoteUserId,
									 data->SocketId ? data->SocketId->SocketName : "",
									 (int)data->ConnectionType, (int)data->NetworkType);
}

static void
p2pPeerConnectionInterruptedTrampoline(const EOS_P2P_OnPeerConnectionInterruptedInfo* data) {
	goP2POnPeerConnectionInterrupted((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
									 (uintptr_t)data->RemoteUserId,
									 data->SocketId ? data->SocketId->SocketName : "");
}

static void p2pPeerConnectionClosedTrampoline(const EOS_P2P_OnRemoteConnectionClosedInfo* data) {
	goP2POnPeerConnectionClosed(
		(uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId, (uintptr_t)data->RemoteUserId,
		data->SocketId ? data->SocketId->SocketName : "", (int)data->Reason);
}

static void
p2pIncomingPacketQueueFullTrampoline(const EOS_P2P_OnIncomingPacketQueueFullInfo* data) {
	goP2POnIncomingPacketQueueFull((uintptr_t)data->ClientData, data->PacketQueueMaxSizeBytes,
								   data->PacketQueueCurrentSizeBytes,
								   (uintptr_t)data->OverflowPacketLocalUserId,
								   data->OverflowPacketChannel, data->OverflowPacketSizeBytes);
}

/* Helper: build a stack EOS_P2P_SocketId from a name. Caller must keep the
 * returned struct alive for the duration of the SDK call. */
static EOS_P2P_SocketId makeSocketId(const char* name) {
	EOS_P2P_SocketId s = {0};
	s.ApiVersion = EOS_P2P_SOCKETID_API_LATEST;
	if (name != NULL) {
		strncpy(s.SocketName, name, EOS_P2P_SOCKETID_SOCKETNAME_SIZE - 1);
	}
	return s;
}

/* Synchronous send / receive */

int eos_p2p_send_packet(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
						const char* socketName, uint8_t channel, const void* data, uint32_t dataLen,
						int allowDelayedDelivery, int reliability,
						int disableAutoAcceptConnection) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);

	EOS_P2P_SendPacketOptions opts = {0};
	opts.ApiVersion = EOS_P2P_SENDPACKET_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.RemoteUserId = (EOS_ProductUserId)remoteUserId;
	opts.SocketId = &socket;
	opts.Channel = channel;
	opts.DataLengthBytes = dataLen;
	opts.Data = data;
	opts.bAllowDelayedDelivery = allowDelayedDelivery ? EOS_TRUE : EOS_FALSE;
	opts.Reliability = (EOS_EPacketReliability)reliability;
	opts.bDisableAutoAcceptConnection = disableAutoAcceptConnection ? EOS_TRUE : EOS_FALSE;

	return (int)EOS_P2P_SendPacket((EOS_HP2P)handle, &opts);
}

int eos_p2p_get_next_received_packet_size(uintptr_t handle, uintptr_t localUserId, int hasChannel,
										  uint8_t channel, uint32_t* outSizeBytes) {
	EOS_P2P_GetNextReceivedPacketSizeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_GETNEXTRECEIVEDPACKETSIZE_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.RequestedChannel = hasChannel ? &channel : NULL;
	return (int)EOS_P2P_GetNextReceivedPacketSize((EOS_HP2P)handle, &opts, outSizeBytes);
}

int eos_p2p_receive_packet(uintptr_t handle, uintptr_t localUserId, uint32_t maxDataBytes,
						   int hasChannel, uint8_t channel, void* outData, uintptr_t* outPeerId,
						   char* outSocketName, uint8_t* outChannel, uint32_t* outBytesWritten) {
	EOS_P2P_ReceivePacketOptions opts = {0};
	opts.ApiVersion = EOS_P2P_RECEIVEPACKET_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.MaxDataSizeBytes = maxDataBytes;
	opts.RequestedChannel = hasChannel ? &channel : NULL;

	EOS_ProductUserId peerId = NULL;
	EOS_P2P_SocketId socket = {0};
	socket.ApiVersion = EOS_P2P_SOCKETID_API_LATEST;

	int result = (int)EOS_P2P_ReceivePacket((EOS_HP2P)handle, &opts, &peerId, &socket, outChannel,
											outData, outBytesWritten);
	*outPeerId = (uintptr_t)peerId;
	if (outSocketName != NULL) {
		memcpy(outSocketName, socket.SocketName, EOS_P2P_SOCKETID_SOCKETNAME_SIZE);
	}
	return result;
}

/* Connection management */

int eos_p2p_accept_connection(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							  const char* socketName) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_AcceptConnectionOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ACCEPTCONNECTION_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.RemoteUserId = (EOS_ProductUserId)remoteUserId;
	opts.SocketId = &socket;
	return (int)EOS_P2P_AcceptConnection((EOS_HP2P)handle, &opts);
}

int eos_p2p_close_connection(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							 const char* socketName) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_CloseConnectionOptions opts = {0};
	opts.ApiVersion = EOS_P2P_CLOSECONNECTION_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.RemoteUserId = (EOS_ProductUserId)remoteUserId;
	opts.SocketId = &socket;
	return (int)EOS_P2P_CloseConnection((EOS_HP2P)handle, &opts);
}

int eos_p2p_close_connections(uintptr_t handle, uintptr_t localUserId, const char* socketName) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_CloseConnectionsOptions opts = {0};
	opts.ApiVersion = EOS_P2P_CLOSECONNECTIONS_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.SocketId = &socket;
	return (int)EOS_P2P_CloseConnections((EOS_HP2P)handle, &opts);
}

/* NAT / relay / port range */

void eos_p2p_query_nat_type(uintptr_t handle, uintptr_t clientData) {
	EOS_P2P_QueryNATTypeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_QUERYNATTYPE_API_LATEST;
	EOS_P2P_QueryNATType((EOS_HP2P)handle, &opts, (void*)clientData, &p2pQueryNATTypeTrampoline);
}

int eos_p2p_get_nat_type(uintptr_t handle, int* outNATType) {
	EOS_P2P_GetNATTypeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_GETNATTYPE_API_LATEST;
	EOS_ENATType natType = EOS_NAT_Unknown;
	int result = (int)EOS_P2P_GetNATType((EOS_HP2P)handle, &opts, &natType);
	*outNATType = (int)natType;
	return result;
}

int eos_p2p_set_relay_control(uintptr_t handle, int relayControl) {
	EOS_P2P_SetRelayControlOptions opts = {0};
	opts.ApiVersion = EOS_P2P_SETRELAYCONTROL_API_LATEST;
	opts.RelayControl = (EOS_ERelayControl)relayControl;
	return (int)EOS_P2P_SetRelayControl((EOS_HP2P)handle, &opts);
}

int eos_p2p_get_relay_control(uintptr_t handle, int* outRelayControl) {
	EOS_P2P_GetRelayControlOptions opts = {0};
	opts.ApiVersion = EOS_P2P_GETRELAYCONTROL_API_LATEST;
	EOS_ERelayControl rc = EOS_RC_AllowRelays;
	int result = (int)EOS_P2P_GetRelayControl((EOS_HP2P)handle, &opts, &rc);
	*outRelayControl = (int)rc;
	return result;
}

int eos_p2p_set_port_range(uintptr_t handle, uint16_t port, uint16_t maxAdditional) {
	EOS_P2P_SetPortRangeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_SETPORTRANGE_API_LATEST;
	opts.Port = port;
	opts.MaxAdditionalPortsToTry = maxAdditional;
	return (int)EOS_P2P_SetPortRange((EOS_HP2P)handle, &opts);
}

int eos_p2p_get_port_range(uintptr_t handle, uint16_t* outPort, uint16_t* outMaxAdditional) {
	EOS_P2P_GetPortRangeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_GETPORTRANGE_API_LATEST;
	return (int)EOS_P2P_GetPortRange((EOS_HP2P)handle, &opts, outPort, outMaxAdditional);
}

/* Packet queue */

int eos_p2p_set_packet_queue_size(uintptr_t handle, uint64_t incomingMax, uint64_t outgoingMax) {
	EOS_P2P_SetPacketQueueSizeOptions opts = {0};
	opts.ApiVersion = EOS_P2P_SETPACKETQUEUESIZE_API_LATEST;
	opts.IncomingPacketQueueMaxSizeBytes = incomingMax;
	opts.OutgoingPacketQueueMaxSizeBytes = outgoingMax;
	return (int)EOS_P2P_SetPacketQueueSize((EOS_HP2P)handle, &opts);
}

int eos_p2p_get_packet_queue_info(uintptr_t handle, uint64_t* outIncomingMax,
								  uint64_t* outIncomingCurrent, uint64_t* outIncomingPackets,
								  uint64_t* outOutgoingMax, uint64_t* outOutgoingCurrent,
								  uint64_t* outOutgoingPackets) {
	EOS_P2P_GetPacketQueueInfoOptions opts = {0};
	opts.ApiVersion = EOS_P2P_GETPACKETQUEUEINFO_API_LATEST;
	EOS_P2P_PacketQueueInfo info = {0};
	int result = (int)EOS_P2P_GetPacketQueueInfo((EOS_HP2P)handle, &opts, &info);
	*outIncomingMax = info.IncomingPacketQueueMaxSizeBytes;
	*outIncomingCurrent = info.IncomingPacketQueueCurrentSizeBytes;
	*outIncomingPackets = info.IncomingPacketQueueCurrentPacketCount;
	*outOutgoingMax = info.OutgoingPacketQueueMaxSizeBytes;
	*outOutgoingCurrent = info.OutgoingPacketQueueCurrentSizeBytes;
	*outOutgoingPackets = info.OutgoingPacketQueueCurrentPacketCount;
	return result;
}

int eos_p2p_clear_packet_queue(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							   const char* socketName) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_ClearPacketQueueOptions opts = {0};
	opts.ApiVersion = EOS_P2P_CLEARPACKETQUEUE_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.RemoteUserId = (EOS_ProductUserId)remoteUserId;
	opts.SocketId = &socket;
	return (int)EOS_P2P_ClearPacketQueue((EOS_HP2P)handle, &opts);
}

/* Notifications */

uint64_t eos_p2p_add_notify_peer_connection_request(uintptr_t handle, uintptr_t localUserId,
													int hasSocket, const char* socketName,
													uintptr_t clientData) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_AddNotifyPeerConnectionRequestOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ADDNOTIFYPEERCONNECTIONREQUEST_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.SocketId = hasSocket ? &socket : NULL;
	return (uint64_t)EOS_P2P_AddNotifyPeerConnectionRequest(
		(EOS_HP2P)handle, &opts, (void*)clientData, &p2pIncomingConnectionRequestTrampoline);
}

void eos_p2p_remove_notify_peer_connection_request(uintptr_t handle, uint64_t id) {
	EOS_P2P_RemoveNotifyPeerConnectionRequest((EOS_HP2P)handle, (EOS_NotificationId)id);
}

uint64_t eos_p2p_add_notify_peer_connection_established(uintptr_t handle, uintptr_t localUserId,
														int hasSocket, const char* socketName,
														uintptr_t clientData) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_AddNotifyPeerConnectionEstablishedOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ADDNOTIFYPEERCONNECTIONESTABLISHED_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.SocketId = hasSocket ? &socket : NULL;
	return (uint64_t)EOS_P2P_AddNotifyPeerConnectionEstablished(
		(EOS_HP2P)handle, &opts, (void*)clientData, &p2pPeerConnectionEstablishedTrampoline);
}

void eos_p2p_remove_notify_peer_connection_established(uintptr_t handle, uint64_t id) {
	EOS_P2P_RemoveNotifyPeerConnectionEstablished((EOS_HP2P)handle, (EOS_NotificationId)id);
}

uint64_t eos_p2p_add_notify_peer_connection_interrupted(uintptr_t handle, uintptr_t localUserId,
														int hasSocket, const char* socketName,
														uintptr_t clientData) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_AddNotifyPeerConnectionInterruptedOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ADDNOTIFYPEERCONNECTIONINTERRUPTED_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.SocketId = hasSocket ? &socket : NULL;
	return (uint64_t)EOS_P2P_AddNotifyPeerConnectionInterrupted(
		(EOS_HP2P)handle, &opts, (void*)clientData, &p2pPeerConnectionInterruptedTrampoline);
}

void eos_p2p_remove_notify_peer_connection_interrupted(uintptr_t handle, uint64_t id) {
	EOS_P2P_RemoveNotifyPeerConnectionInterrupted((EOS_HP2P)handle, (EOS_NotificationId)id);
}

uint64_t eos_p2p_add_notify_peer_connection_closed(uintptr_t handle, uintptr_t localUserId,
												   int hasSocket, const char* socketName,
												   uintptr_t clientData) {
	EOS_P2P_SocketId socket = makeSocketId(socketName);
	EOS_P2P_AddNotifyPeerConnectionClosedOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ADDNOTIFYPEERCONNECTIONCLOSED_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.SocketId = hasSocket ? &socket : NULL;
	return (uint64_t)EOS_P2P_AddNotifyPeerConnectionClosed(
		(EOS_HP2P)handle, &opts, (void*)clientData, &p2pPeerConnectionClosedTrampoline);
}

void eos_p2p_remove_notify_peer_connection_closed(uintptr_t handle, uint64_t id) {
	EOS_P2P_RemoveNotifyPeerConnectionClosed((EOS_HP2P)handle, (EOS_NotificationId)id);
}

uint64_t eos_p2p_add_notify_incoming_packet_queue_full(uintptr_t handle, uintptr_t clientData) {
	EOS_P2P_AddNotifyIncomingPacketQueueFullOptions opts = {0};
	opts.ApiVersion = EOS_P2P_ADDNOTIFYINCOMINGPACKETQUEUEFULL_API_LATEST;
	return (uint64_t)EOS_P2P_AddNotifyIncomingPacketQueueFull(
		(EOS_HP2P)handle, &opts, (void*)clientData, &p2pIncomingPacketQueueFullTrampoline);
}

void eos_p2p_remove_notify_incoming_packet_queue_full(uintptr_t handle, uint64_t id) {
	EOS_P2P_RemoveNotifyIncomingPacketQueueFull((EOS_HP2P)handle, (EOS_NotificationId)id);
}

#endif /* EOS_CGO */
