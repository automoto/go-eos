#ifndef P2P_WRAPPER_H
#define P2P_WRAPPER_H

#include <stdint.h>

/* Synchronous send / receive */

int eos_p2p_send_packet(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
						const char* socketName, uint8_t channel, const void* data, uint32_t dataLen,
						int allowDelayedDelivery, int reliability, int disableAutoAcceptConnection);

int eos_p2p_get_next_received_packet_size(uintptr_t handle, uintptr_t localUserId, int hasChannel,
										  uint8_t channel, uint32_t* outSizeBytes);

int eos_p2p_receive_packet(uintptr_t handle, uintptr_t localUserId, uint32_t maxDataBytes,
						   int hasChannel, uint8_t channel, void* outData, uintptr_t* outPeerId,
						   char* outSocketName, uint8_t* outChannel, uint32_t* outBytesWritten);

/* Connection management */

int eos_p2p_accept_connection(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							  const char* socketName);
int eos_p2p_close_connection(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							 const char* socketName);
int eos_p2p_close_connections(uintptr_t handle, uintptr_t localUserId, const char* socketName);

/* NAT / relay / port range */

void eos_p2p_query_nat_type(uintptr_t handle, uintptr_t clientData);
int eos_p2p_get_nat_type(uintptr_t handle, int* outNATType);
int eos_p2p_set_relay_control(uintptr_t handle, int relayControl);
int eos_p2p_get_relay_control(uintptr_t handle, int* outRelayControl);
int eos_p2p_set_port_range(uintptr_t handle, uint16_t port, uint16_t maxAdditional);
int eos_p2p_get_port_range(uintptr_t handle, uint16_t* outPort, uint16_t* outMaxAdditional);

/* Packet queue */

int eos_p2p_set_packet_queue_size(uintptr_t handle, uint64_t incomingMax, uint64_t outgoingMax);
int eos_p2p_get_packet_queue_info(uintptr_t handle, uint64_t* outIncomingMax,
								  uint64_t* outIncomingCurrent, uint64_t* outIncomingPackets,
								  uint64_t* outOutgoingMax, uint64_t* outOutgoingCurrent,
								  uint64_t* outOutgoingPackets);
int eos_p2p_clear_packet_queue(uintptr_t handle, uintptr_t localUserId, uintptr_t remoteUserId,
							   const char* socketName);

/* Notifications */

uint64_t eos_p2p_add_notify_peer_connection_request(uintptr_t handle, uintptr_t localUserId,
													int hasSocket, const char* socketName,
													uintptr_t clientData);
void eos_p2p_remove_notify_peer_connection_request(uintptr_t handle, uint64_t id);

uint64_t eos_p2p_add_notify_peer_connection_established(uintptr_t handle, uintptr_t localUserId,
														int hasSocket, const char* socketName,
														uintptr_t clientData);
void eos_p2p_remove_notify_peer_connection_established(uintptr_t handle, uint64_t id);

uint64_t eos_p2p_add_notify_peer_connection_interrupted(uintptr_t handle, uintptr_t localUserId,
														int hasSocket, const char* socketName,
														uintptr_t clientData);
void eos_p2p_remove_notify_peer_connection_interrupted(uintptr_t handle, uint64_t id);

uint64_t eos_p2p_add_notify_peer_connection_closed(uintptr_t handle, uintptr_t localUserId,
												   int hasSocket, const char* socketName,
												   uintptr_t clientData);
void eos_p2p_remove_notify_peer_connection_closed(uintptr_t handle, uint64_t id);

uint64_t eos_p2p_add_notify_incoming_packet_queue_full(uintptr_t handle, uintptr_t clientData);
void eos_p2p_remove_notify_incoming_packet_queue_full(uintptr_t handle, uint64_t id);

#endif
