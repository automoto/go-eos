#ifndef LOBBY_WRAPPER_H
#define LOBBY_WRAPPER_H

#include <stdint.h>

/* Core lifecycle (async) */
void eos_lobby_create(uintptr_t handle, uintptr_t localUserId, uint32_t maxMembers,
					  int permissionLevel, int allowInvites, const char* bucketId,
					  uintptr_t clientData);
void eos_lobby_destroy(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
					   uintptr_t clientData);
void eos_lobby_join(uintptr_t handle, uintptr_t detailsHandle, uintptr_t localUserId,
					uintptr_t clientData);
void eos_lobby_leave(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
					 uintptr_t clientData);
void eos_lobby_update(uintptr_t handle, uintptr_t modHandle, uintptr_t clientData);
void eos_lobby_send_invite(uintptr_t handle, const char* lobbyId, uintptr_t localUserId,
						   uintptr_t targetUserId, uintptr_t clientData);
void eos_lobby_query_invites(uintptr_t handle, uintptr_t localUserId, uintptr_t clientData);

/* Modification handle */
int eos_lobby_update_lobby_modification(uintptr_t handle, uintptr_t localUserId,
										const char* lobbyId, uintptr_t* outMod);
int eos_lobby_mod_set_bucket_id(uintptr_t mod, const char* bucketId);
int eos_lobby_mod_set_permission_level(uintptr_t mod, int level);
int eos_lobby_mod_set_max_members(uintptr_t mod, uint32_t maxMembers);
int eos_lobby_mod_set_invites_allowed(uintptr_t mod, int allowed);
int eos_lobby_mod_add_attr_int64(uintptr_t mod, const char* key, int64_t val, int visibility);
int eos_lobby_mod_add_attr_double(uintptr_t mod, const char* key, double val, int visibility);
int eos_lobby_mod_add_attr_bool(uintptr_t mod, const char* key, int val, int visibility);
int eos_lobby_mod_add_attr_string(uintptr_t mod, const char* key, const char* val, int visibility);
int eos_lobby_mod_remove_attr(uintptr_t mod, const char* key);
int eos_lobby_mod_add_member_attr_int64(uintptr_t mod, const char* key, int64_t val,
										int visibility);
int eos_lobby_mod_add_member_attr_double(uintptr_t mod, const char* key, double val,
										 int visibility);
int eos_lobby_mod_add_member_attr_bool(uintptr_t mod, const char* key, int val, int visibility);
int eos_lobby_mod_add_member_attr_string(uintptr_t mod, const char* key, const char* val,
										 int visibility);
int eos_lobby_mod_remove_member_attr(uintptr_t mod, const char* key);
void eos_lobby_mod_release(uintptr_t mod);

/* Search */
int eos_lobby_create_search(uintptr_t handle, uint32_t maxResults, uintptr_t* outSearch);
void eos_lobby_search_find(uintptr_t search, uintptr_t localUserId, uintptr_t clientData);
int eos_lobby_search_set_param_int64(uintptr_t search, const char* key, int64_t val, int op);
int eos_lobby_search_set_param_double(uintptr_t search, const char* key, double val, int op);
int eos_lobby_search_set_param_bool(uintptr_t search, const char* key, int val, int op);
int eos_lobby_search_set_param_string(uintptr_t search, const char* key, const char* val, int op);
int eos_lobby_search_set_lobby_id(uintptr_t search, const char* lobbyId);
int eos_lobby_search_set_target_user_id(uintptr_t search, uintptr_t targetUserId);
int eos_lobby_search_set_max_results(uintptr_t search, uint32_t maxResults);
uint32_t eos_lobby_search_get_search_result_count(uintptr_t search);
int eos_lobby_search_copy_search_result_by_index(uintptr_t search, uint32_t index,
												 uintptr_t* outDetails);
void eos_lobby_search_release(uintptr_t search);

/* Details */
int eos_lobby_copy_details_handle(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
								  uintptr_t* outDetails);
uintptr_t eos_lobby_details_get_owner(uintptr_t details);
uint32_t eos_lobby_details_get_member_count(uintptr_t details);
uintptr_t eos_lobby_details_get_member_by_index(uintptr_t details, uint32_t index);
uint32_t eos_lobby_details_get_attribute_count(uintptr_t details);
uint32_t eos_lobby_details_get_member_attribute_count(uintptr_t details, uintptr_t targetUserId);

/* Details info accessors (from CopyInfo) */
int eos_lobby_details_copy_info(uintptr_t details, uintptr_t* outInfo);
const char* eos_lobby_info_get_lobby_id(uintptr_t info);
uintptr_t eos_lobby_info_get_owner(uintptr_t info);
int eos_lobby_info_get_permission_level(uintptr_t info);
uint32_t eos_lobby_info_get_available_slots(uintptr_t info);
uint32_t eos_lobby_info_get_max_members(uintptr_t info);
int eos_lobby_info_get_allow_invites(uintptr_t info);
const char* eos_lobby_info_get_bucket_id(uintptr_t info);
void eos_lobby_details_info_release(uintptr_t info);

/* Attribute accessors (from CopyAttributeByIndex/ByKey) */
int eos_lobby_details_copy_attr_by_index(uintptr_t details, uint32_t index, uintptr_t* outAttr);
int eos_lobby_details_copy_attr_by_key(uintptr_t details, const char* key, uintptr_t* outAttr);
int eos_lobby_details_copy_member_attr_by_index(uintptr_t details, uintptr_t targetUserId,
												uint32_t index, uintptr_t* outAttr);
int eos_lobby_details_copy_member_attr_by_key(uintptr_t details, uintptr_t targetUserId,
											  const char* key, uintptr_t* outAttr);
const char* eos_lobby_attr_get_key(uintptr_t attr);
int eos_lobby_attr_get_type(uintptr_t attr);
int64_t eos_lobby_attr_get_int64(uintptr_t attr);
double eos_lobby_attr_get_double(uintptr_t attr);
int eos_lobby_attr_get_bool(uintptr_t attr);
const char* eos_lobby_attr_get_string(uintptr_t attr);
int eos_lobby_attr_get_visibility(uintptr_t attr);
void eos_lobby_attr_release(uintptr_t attr);
void eos_lobby_details_release(uintptr_t details);

/* Invite helpers */
uint32_t eos_lobby_get_invite_count(uintptr_t handle, uintptr_t localUserId);
int eos_lobby_get_invite_id_by_index(uintptr_t handle, uintptr_t localUserId, uint32_t index,
									 char* outBuf, int32_t* outBufLen);
int eos_lobby_copy_details_handle_by_invite_id(uintptr_t handle, const char* inviteId,
											   uintptr_t* outDetails);

/* Notifications */
uint64_t eos_lobby_add_notify_update_received(uintptr_t handle, uintptr_t clientData);
void eos_lobby_remove_notify_update_received(uintptr_t handle, uint64_t id);
uint64_t eos_lobby_add_notify_member_update_received(uintptr_t handle, uintptr_t clientData);
void eos_lobby_remove_notify_member_update_received(uintptr_t handle, uint64_t id);
uint64_t eos_lobby_add_notify_member_status_received(uintptr_t handle, uintptr_t clientData);
void eos_lobby_remove_notify_member_status_received(uintptr_t handle, uint64_t id);
uint64_t eos_lobby_add_notify_invite_received(uintptr_t handle, uintptr_t clientData);
void eos_lobby_remove_notify_invite_received(uintptr_t handle, uint64_t id);

#endif
