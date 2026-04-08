#ifndef SESSIONS_WRAPPER_H
#define SESSIONS_WRAPPER_H

#include <stdint.h>

/* Core lifecycle (async) */
void eos_sessions_update_session(uintptr_t handle, uintptr_t modHandle, uintptr_t clientData);
void eos_sessions_destroy_session(uintptr_t handle, const char* sessionName, uintptr_t clientData);
void eos_sessions_join_session(uintptr_t handle, const char* sessionName, uintptr_t sessionDetails,
							   uintptr_t localUserId, uintptr_t clientData);
void eos_sessions_start_session(uintptr_t handle, const char* sessionName, uintptr_t clientData);
void eos_sessions_end_session(uintptr_t handle, const char* sessionName, uintptr_t clientData);
void eos_sessions_register_players(uintptr_t handle, const char* sessionName, uintptr_t* playerIds,
								   uint32_t count, uintptr_t clientData);
void eos_sessions_unregister_players(uintptr_t handle, const char* sessionName,
									 uintptr_t* playerIds, uint32_t count, uintptr_t clientData);

/* Modification handle */
int eos_sessions_create_session_modification(uintptr_t handle, const char* sessionName,
											 const char* bucketId, uint32_t maxPlayers,
											 uintptr_t localUserId, uintptr_t* outMod);
int eos_sessions_update_session_modification(uintptr_t handle, const char* sessionName,
											 uintptr_t* outMod);
int eos_session_mod_set_bucket_id(uintptr_t mod, const char* bucketId);
int eos_session_mod_set_permission_level(uintptr_t mod, int level);
int eos_session_mod_set_max_players(uintptr_t mod, uint32_t maxPlayers);
int eos_session_mod_set_join_in_progress_allowed(uintptr_t mod, int allowed);
int eos_session_mod_set_invites_allowed(uintptr_t mod, int allowed);
int eos_session_mod_set_host_address(uintptr_t mod, const char* addr);
int eos_session_mod_add_attr_int64(uintptr_t mod, const char* key, int64_t val, int advType);
int eos_session_mod_add_attr_double(uintptr_t mod, const char* key, double val, int advType);
int eos_session_mod_add_attr_bool(uintptr_t mod, const char* key, int val, int advType);
int eos_session_mod_add_attr_string(uintptr_t mod, const char* key, const char* val, int advType);
int eos_session_mod_remove_attr(uintptr_t mod, const char* key);
void eos_session_mod_release(uintptr_t mod);

/* Search */
int eos_sessions_create_search(uintptr_t handle, uint32_t maxResults, uintptr_t* outSearch);
void eos_session_search_find(uintptr_t search, uintptr_t localUserId, uintptr_t clientData);
int eos_session_search_set_param_int64(uintptr_t search, const char* key, int64_t val, int op);
int eos_session_search_set_param_double(uintptr_t search, const char* key, double val, int op);
int eos_session_search_set_param_bool(uintptr_t search, const char* key, int val, int op);
int eos_session_search_set_param_string(uintptr_t search, const char* key, const char* val, int op);
int eos_session_search_set_session_id(uintptr_t search, const char* sessionId);
int eos_session_search_set_target_user_id(uintptr_t search, uintptr_t targetUserId);
int eos_session_search_set_max_results(uintptr_t search, uint32_t maxResults);
uint32_t eos_session_search_get_result_count(uintptr_t search);
int eos_session_search_copy_result_by_index(uintptr_t search, uint32_t index,
											uintptr_t* outDetails);
void eos_session_search_release(uintptr_t search);

/* Session details */
int eos_session_details_copy_info(uintptr_t details, uintptr_t* outInfo);
const char* eos_session_info_get_session_id(uintptr_t info);
const char* eos_session_info_get_host_address(uintptr_t info);
uint32_t eos_session_info_get_num_open_connections(uintptr_t info);
uintptr_t eos_session_info_get_owner(uintptr_t info);
const char* eos_session_info_get_bucket_id(uintptr_t info);
uint32_t eos_session_info_get_num_connections(uintptr_t info);
int eos_session_info_get_allow_join_in_progress(uintptr_t info);
int eos_session_info_get_permission_level(uintptr_t info);
int eos_session_info_get_invites_allowed(uintptr_t info);
void eos_session_details_info_release(uintptr_t info);
uint32_t eos_session_details_get_attribute_count(uintptr_t details);
int eos_session_details_copy_attr_by_index(uintptr_t details, uint32_t index, uintptr_t* outAttr);
int eos_session_details_copy_attr_by_key(uintptr_t details, const char* key, uintptr_t* outAttr);
const char* eos_session_attr_get_key(uintptr_t attr);
int eos_session_attr_get_type(uintptr_t attr);
int64_t eos_session_attr_get_int64(uintptr_t attr);
double eos_session_attr_get_double(uintptr_t attr);
int eos_session_attr_get_bool(uintptr_t attr);
const char* eos_session_attr_get_string(uintptr_t attr);
int eos_session_attr_get_advertisement_type(uintptr_t attr);
void eos_session_attr_release(uintptr_t attr);
void eos_session_details_release(uintptr_t details);

/* Active session */
int eos_sessions_copy_active_session_handle(uintptr_t handle, const char* sessionName,
											uintptr_t* outActive);
const char* eos_active_session_get_name(uintptr_t active);
uintptr_t eos_active_session_get_local_user(uintptr_t active);
int eos_active_session_get_state(uintptr_t active);
uint32_t eos_active_session_get_registered_player_count(uintptr_t active);
uintptr_t eos_active_session_get_registered_player_by_index(uintptr_t active, uint32_t index);
void eos_active_session_info_release(uintptr_t info);
void eos_active_session_release(uintptr_t active);

/* Notifications */
uint64_t eos_sessions_add_notify_invite_received(uintptr_t handle, uintptr_t clientData);
void eos_sessions_remove_notify_invite_received(uintptr_t handle, uint64_t id);

#endif
