// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_sessions.h"
#include "cgo_helpers.h"
#include <stdint.h>

/* Forward declarations for Go export functions */
extern void goSessionsUpdateCallback(int resultCode, uintptr_t clientData, const char* sessionName,
									 const char* sessionId);
extern void goSessionsDestroyCallback(int resultCode, uintptr_t clientData);
extern void goSessionsJoinCallback(int resultCode, uintptr_t clientData);
extern void goSessionsStartCallback(int resultCode, uintptr_t clientData);
extern void goSessionsEndCallback(int resultCode, uintptr_t clientData);
extern void goSessionsRegisterPlayersCallback(int resultCode, uintptr_t clientData);
extern void goSessionsUnregisterPlayersCallback(int resultCode, uintptr_t clientData);
extern void goSessionSearchFindCallback(int resultCode, uintptr_t clientData);
extern void goSessionsInviteReceivedCallback(uintptr_t clientData, uintptr_t localUserId,
											 uintptr_t targetUserId, const char* inviteId);

/* Trampolines */

static void sessionsUpdateTrampoline(const EOS_Sessions_UpdateSessionCallbackInfo* data) {
	goSessionsUpdateCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
							 safe_str(data->SessionName), safe_str(data->SessionId));
}

static void sessionsDestroyTrampoline(const EOS_Sessions_DestroySessionCallbackInfo* data) {
	goSessionsDestroyCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void sessionsJoinTrampoline(const EOS_Sessions_JoinSessionCallbackInfo* data) {
	goSessionsJoinCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void sessionsStartTrampoline(const EOS_Sessions_StartSessionCallbackInfo* data) {
	goSessionsStartCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void sessionsEndTrampoline(const EOS_Sessions_EndSessionCallbackInfo* data) {
	goSessionsEndCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void
sessionsRegisterPlayersTrampoline(const EOS_Sessions_RegisterPlayersCallbackInfo* data) {
	goSessionsRegisterPlayersCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void
sessionsUnregisterPlayersTrampoline(const EOS_Sessions_UnregisterPlayersCallbackInfo* data) {
	goSessionsUnregisterPlayersCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void sessionSearchFindTrampoline(const EOS_SessionSearch_FindCallbackInfo* data) {
	goSessionSearchFindCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void
sessionsInviteReceivedTrampoline(const EOS_Sessions_SessionInviteReceivedCallbackInfo* data) {
	goSessionsInviteReceivedCallback((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
									 (uintptr_t)data->TargetUserId, safe_str(data->InviteId));
}

/* Core lifecycle wrappers */

void eos_sessions_update_session(uintptr_t handle, uintptr_t modHandle, uintptr_t clientData) {
	EOS_Sessions_UpdateSessionOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_UPDATESESSION_API_LATEST;
	opts.SessionModificationHandle = (EOS_HSessionModification)modHandle;
	EOS_Sessions_UpdateSession((EOS_HSessions)handle, &opts, (void*)clientData,
							   &sessionsUpdateTrampoline);
}

void eos_sessions_destroy_session(uintptr_t handle, const char* sessionName, uintptr_t clientData) {
	EOS_Sessions_DestroySessionOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_DESTROYSESSION_API_LATEST;
	opts.SessionName = sessionName;
	EOS_Sessions_DestroySession((EOS_HSessions)handle, &opts, (void*)clientData,
								&sessionsDestroyTrampoline);
}

void eos_sessions_join_session(uintptr_t handle, const char* sessionName, uintptr_t sessionDetails,
							   uintptr_t localUserId, uintptr_t clientData) {
	EOS_Sessions_JoinSessionOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_JOINSESSION_API_LATEST;
	opts.SessionName = sessionName;
	opts.SessionHandle = (EOS_HSessionDetails)sessionDetails;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_Sessions_JoinSession((EOS_HSessions)handle, &opts, (void*)clientData,
							 &sessionsJoinTrampoline);
}

void eos_sessions_start_session(uintptr_t handle, const char* sessionName, uintptr_t clientData) {
	EOS_Sessions_StartSessionOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_STARTSESSION_API_LATEST;
	opts.SessionName = sessionName;
	EOS_Sessions_StartSession((EOS_HSessions)handle, &opts, (void*)clientData,
							  &sessionsStartTrampoline);
}

void eos_sessions_end_session(uintptr_t handle, const char* sessionName, uintptr_t clientData) {
	EOS_Sessions_EndSessionOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_ENDSESSION_API_LATEST;
	opts.SessionName = sessionName;
	EOS_Sessions_EndSession((EOS_HSessions)handle, &opts, (void*)clientData,
							&sessionsEndTrampoline);
}

void eos_sessions_register_players(uintptr_t handle, const char* sessionName, uintptr_t* playerIds,
								   uint32_t count, uintptr_t clientData) {
	EOS_Sessions_RegisterPlayersOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_REGISTERPLAYERS_API_LATEST;
	opts.SessionName = sessionName;
	opts.PlayersToRegister = (EOS_ProductUserId*)playerIds;
	opts.PlayersToRegisterCount = count;
	EOS_Sessions_RegisterPlayers((EOS_HSessions)handle, &opts, (void*)clientData,
								 &sessionsRegisterPlayersTrampoline);
}

void eos_sessions_unregister_players(uintptr_t handle, const char* sessionName,
									 uintptr_t* playerIds, uint32_t count, uintptr_t clientData) {
	EOS_Sessions_UnregisterPlayersOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_UNREGISTERPLAYERS_API_LATEST;
	opts.SessionName = sessionName;
	opts.PlayersToUnregister = (EOS_ProductUserId*)playerIds;
	opts.PlayersToUnregisterCount = count;
	EOS_Sessions_UnregisterPlayers((EOS_HSessions)handle, &opts, (void*)clientData,
								   &sessionsUnregisterPlayersTrampoline);
}

/* Modification handle */

int eos_sessions_create_session_modification(uintptr_t handle, const char* sessionName,
											 const char* bucketId, uint32_t maxPlayers,
											 uintptr_t localUserId, uintptr_t* outMod) {
	EOS_Sessions_CreateSessionModificationOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_CREATESESSIONMODIFICATION_API_LATEST;
	opts.SessionName = sessionName;
	opts.BucketId = bucketId;
	opts.MaxPlayers = maxPlayers;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_HSessionModification mod = NULL;
	EOS_EResult result = EOS_Sessions_CreateSessionModification((EOS_HSessions)handle, &opts, &mod);
	*outMod = (uintptr_t)mod;
	return (int)result;
}

int eos_sessions_update_session_modification(uintptr_t handle, const char* sessionName,
											 uintptr_t* outMod) {
	EOS_Sessions_UpdateSessionModificationOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_UPDATESESSIONMODIFICATION_API_LATEST;
	opts.SessionName = sessionName;
	EOS_HSessionModification mod = NULL;
	EOS_EResult result = EOS_Sessions_UpdateSessionModification((EOS_HSessions)handle, &opts, &mod);
	*outMod = (uintptr_t)mod;
	return (int)result;
}

int eos_session_mod_set_bucket_id(uintptr_t mod, const char* bucketId) {
	EOS_SessionModification_SetBucketIdOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETBUCKETID_API_LATEST;
	opts.BucketId = bucketId;
	return (int)EOS_SessionModification_SetBucketId((EOS_HSessionModification)mod, &opts);
}

int eos_session_mod_set_permission_level(uintptr_t mod, int level) {
	EOS_SessionModification_SetPermissionLevelOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETPERMISSIONLEVEL_API_LATEST;
	opts.PermissionLevel = (EOS_EOnlineSessionPermissionLevel)level;
	return (int)EOS_SessionModification_SetPermissionLevel((EOS_HSessionModification)mod, &opts);
}

int eos_session_mod_set_max_players(uintptr_t mod, uint32_t maxPlayers) {
	EOS_SessionModification_SetMaxPlayersOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETMAXPLAYERS_API_LATEST;
	opts.MaxPlayers = maxPlayers;
	return (int)EOS_SessionModification_SetMaxPlayers((EOS_HSessionModification)mod, &opts);
}

int eos_session_mod_set_join_in_progress_allowed(uintptr_t mod, int allowed) {
	EOS_SessionModification_SetJoinInProgressAllowedOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETJOININPROGRESSALLOWED_API_LATEST;
	opts.bAllowJoinInProgress = allowed ? EOS_TRUE : EOS_FALSE;
	return (int)EOS_SessionModification_SetJoinInProgressAllowed((EOS_HSessionModification)mod,
																 &opts);
}

int eos_session_mod_set_invites_allowed(uintptr_t mod, int allowed) {
	EOS_SessionModification_SetInvitesAllowedOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETINVITESALLOWED_API_LATEST;
	opts.bInvitesAllowed = allowed ? EOS_TRUE : EOS_FALSE;
	return (int)EOS_SessionModification_SetInvitesAllowed((EOS_HSessionModification)mod, &opts);
}

int eos_session_mod_set_host_address(uintptr_t mod, const char* addr) {
	EOS_SessionModification_SetHostAddressOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_SETHOSTADDRESS_API_LATEST;
	opts.HostAddress = addr;
	return (int)EOS_SessionModification_SetHostAddress((EOS_HSessionModification)mod, &opts);
}

static int session_mod_add_attr(uintptr_t mod, EOS_Sessions_AttributeData* data, int advType) {
	EOS_SessionModification_AddAttributeOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_ADDATTRIBUTE_API_LATEST;
	opts.SessionAttribute = data;
	opts.AdvertisementType = (EOS_ESessionAttributeAdvertisementType)advType;
	return (int)EOS_SessionModification_AddAttribute((EOS_HSessionModification)mod, &opts);
}

int eos_session_mod_add_attr_int64(uintptr_t mod, const char* key, int64_t val, int advType) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsInt64 = val;
	data.ValueType = EOS_AT_INT64;
	return session_mod_add_attr(mod, &data, advType);
}

int eos_session_mod_add_attr_double(uintptr_t mod, const char* key, double val, int advType) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsDouble = val;
	data.ValueType = EOS_AT_DOUBLE;
	return session_mod_add_attr(mod, &data, advType);
}

int eos_session_mod_add_attr_bool(uintptr_t mod, const char* key, int val, int advType) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsBool = val ? EOS_TRUE : EOS_FALSE;
	data.ValueType = EOS_AT_BOOLEAN;
	return session_mod_add_attr(mod, &data, advType);
}

int eos_session_mod_add_attr_string(uintptr_t mod, const char* key, const char* val, int advType) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsUtf8 = val;
	data.ValueType = EOS_AT_STRING;
	return session_mod_add_attr(mod, &data, advType);
}

int eos_session_mod_remove_attr(uintptr_t mod, const char* key) {
	EOS_SessionModification_RemoveAttributeOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONMODIFICATION_REMOVEATTRIBUTE_API_LATEST;
	opts.Key = key;
	return (int)EOS_SessionModification_RemoveAttribute((EOS_HSessionModification)mod, &opts);
}

void eos_session_mod_release(uintptr_t mod) {
	EOS_SessionModification_Release((EOS_HSessionModification)mod);
}

/* Search */

int eos_sessions_create_search(uintptr_t handle, uint32_t maxResults, uintptr_t* outSearch) {
	EOS_Sessions_CreateSessionSearchOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_CREATESESSIONSEARCH_API_LATEST;
	opts.MaxSearchResults = maxResults;
	EOS_HSessionSearch search = NULL;
	EOS_EResult result = EOS_Sessions_CreateSessionSearch((EOS_HSessions)handle, &opts, &search);
	*outSearch = (uintptr_t)search;
	return (int)result;
}

void eos_session_search_find(uintptr_t search, uintptr_t localUserId, uintptr_t clientData) {
	EOS_SessionSearch_FindOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_FIND_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_SessionSearch_Find((EOS_HSessionSearch)search, &opts, (void*)clientData,
						   &sessionSearchFindTrampoline);
}

static int session_search_set_param(uintptr_t search, EOS_Sessions_AttributeData* data, int op) {
	EOS_SessionSearch_SetParameterOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_SETPARAMETER_API_LATEST;
	opts.Parameter = data;
	opts.ComparisonOp = (EOS_EComparisonOp)op;
	return (int)EOS_SessionSearch_SetParameter((EOS_HSessionSearch)search, &opts);
}

int eos_session_search_set_param_int64(uintptr_t search, const char* key, int64_t val, int op) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsInt64 = val;
	data.ValueType = EOS_AT_INT64;
	return session_search_set_param(search, &data, op);
}

int eos_session_search_set_param_double(uintptr_t search, const char* key, double val, int op) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsDouble = val;
	data.ValueType = EOS_AT_DOUBLE;
	return session_search_set_param(search, &data, op);
}

int eos_session_search_set_param_bool(uintptr_t search, const char* key, int val, int op) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsBool = val ? EOS_TRUE : EOS_FALSE;
	data.ValueType = EOS_AT_BOOLEAN;
	return session_search_set_param(search, &data, op);
}

int eos_session_search_set_param_string(uintptr_t search, const char* key, const char* val,
										int op) {
	EOS_Sessions_AttributeData data = {0};
	data.ApiVersion = EOS_SESSIONS_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsUtf8 = val;
	data.ValueType = EOS_AT_STRING;
	return session_search_set_param(search, &data, op);
}

int eos_session_search_set_session_id(uintptr_t search, const char* sessionId) {
	EOS_SessionSearch_SetSessionIdOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_SETSESSIONID_API_LATEST;
	opts.SessionId = sessionId;
	return (int)EOS_SessionSearch_SetSessionId((EOS_HSessionSearch)search, &opts);
}

int eos_session_search_set_target_user_id(uintptr_t search, uintptr_t targetUserId) {
	EOS_SessionSearch_SetTargetUserIdOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_SETTARGETUSERID_API_LATEST;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	return (int)EOS_SessionSearch_SetTargetUserId((EOS_HSessionSearch)search, &opts);
}

int eos_session_search_set_max_results(uintptr_t search, uint32_t maxResults) {
	EOS_SessionSearch_SetMaxResultsOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_SETMAXSEARCHRESULTS_API_LATEST;
	opts.MaxSearchResults = maxResults;
	return (int)EOS_SessionSearch_SetMaxResults((EOS_HSessionSearch)search, &opts);
}

uint32_t eos_session_search_get_result_count(uintptr_t search) {
	EOS_SessionSearch_GetSearchResultCountOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_GETSEARCHRESULTCOUNT_API_LATEST;
	return EOS_SessionSearch_GetSearchResultCount((EOS_HSessionSearch)search, &opts);
}

int eos_session_search_copy_result_by_index(uintptr_t search, uint32_t index,
											uintptr_t* outDetails) {
	EOS_SessionSearch_CopySearchResultByIndexOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONSEARCH_COPYSEARCHRESULTBYINDEX_API_LATEST;
	opts.SessionIndex = index;
	EOS_HSessionDetails details = NULL;
	EOS_EResult result =
		EOS_SessionSearch_CopySearchResultByIndex((EOS_HSessionSearch)search, &opts, &details);
	*outDetails = (uintptr_t)details;
	return (int)result;
}

void eos_session_search_release(uintptr_t search) {
	EOS_SessionSearch_Release((EOS_HSessionSearch)search);
}

/* Session details */

int eos_session_details_copy_info(uintptr_t details, uintptr_t* outInfo) {
	EOS_SessionDetails_CopyInfoOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONDETAILS_COPYINFO_API_LATEST;
	EOS_SessionDetails_Info* info = NULL;
	EOS_EResult result = EOS_SessionDetails_CopyInfo((EOS_HSessionDetails)details, &opts, &info);
	*outInfo = (uintptr_t)info;
	return (int)result;
}

const char* eos_session_info_get_session_id(uintptr_t info) {
	return safe_str(((EOS_SessionDetails_Info*)info)->SessionId);
}
const char* eos_session_info_get_host_address(uintptr_t info) {
	return safe_str(((EOS_SessionDetails_Info*)info)->HostAddress);
}
uint32_t eos_session_info_get_num_open_connections(uintptr_t info) {
	return ((EOS_SessionDetails_Info*)info)->NumOpenPublicConnections;
}
uintptr_t eos_session_info_get_owner(uintptr_t info) {
	return (uintptr_t)((EOS_SessionDetails_Info*)info)->OwnerUserId;
}
const char* eos_session_info_get_bucket_id(uintptr_t info) {
	const EOS_SessionDetails_Settings* s = ((EOS_SessionDetails_Info*)info)->Settings;
	return s ? safe_str(s->BucketId) : "";
}
uint32_t eos_session_info_get_num_connections(uintptr_t info) {
	const EOS_SessionDetails_Settings* s = ((EOS_SessionDetails_Info*)info)->Settings;
	return s ? s->NumPublicConnections : 0;
}
int eos_session_info_get_allow_join_in_progress(uintptr_t info) {
	const EOS_SessionDetails_Settings* s = ((EOS_SessionDetails_Info*)info)->Settings;
	return s ? (s->bAllowJoinInProgress == EOS_TRUE) : 0;
}
int eos_session_info_get_permission_level(uintptr_t info) {
	const EOS_SessionDetails_Settings* s = ((EOS_SessionDetails_Info*)info)->Settings;
	return s ? (int)s->PermissionLevel : 0;
}
int eos_session_info_get_invites_allowed(uintptr_t info) {
	const EOS_SessionDetails_Settings* s = ((EOS_SessionDetails_Info*)info)->Settings;
	return s ? (s->bInvitesAllowed == EOS_TRUE) : 0;
}
void eos_session_details_info_release(uintptr_t info) {
	EOS_SessionDetails_Info_Release((EOS_SessionDetails_Info*)info);
}

uint32_t eos_session_details_get_attribute_count(uintptr_t details) {
	EOS_SessionDetails_GetSessionAttributeCountOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONDETAILS_GETSESSIONATTRIBUTECOUNT_API_LATEST;
	return EOS_SessionDetails_GetSessionAttributeCount((EOS_HSessionDetails)details, &opts);
}

int eos_session_details_copy_attr_by_index(uintptr_t details, uint32_t index, uintptr_t* outAttr) {
	EOS_SessionDetails_CopySessionAttributeByIndexOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONDETAILS_COPYSESSIONATTRIBUTEBYINDEX_API_LATEST;
	opts.AttrIndex = index;
	EOS_SessionDetails_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_SessionDetails_CopySessionAttributeByIndex((EOS_HSessionDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

int eos_session_details_copy_attr_by_key(uintptr_t details, const char* key, uintptr_t* outAttr) {
	EOS_SessionDetails_CopySessionAttributeByKeyOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONDETAILS_COPYSESSIONATTRIBUTEBYKEY_API_LATEST;
	opts.AttrKey = key;
	EOS_SessionDetails_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_SessionDetails_CopySessionAttributeByKey((EOS_HSessionDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

const char* eos_session_attr_get_key(uintptr_t attr) {
	return safe_str(((EOS_SessionDetails_Attribute*)attr)->Data->Key);
}
int eos_session_attr_get_type(uintptr_t attr) {
	return (int)((EOS_SessionDetails_Attribute*)attr)->Data->ValueType;
}
int64_t eos_session_attr_get_int64(uintptr_t attr) {
	return ((EOS_SessionDetails_Attribute*)attr)->Data->Value.AsInt64;
}
double eos_session_attr_get_double(uintptr_t attr) {
	return ((EOS_SessionDetails_Attribute*)attr)->Data->Value.AsDouble;
}
int eos_session_attr_get_bool(uintptr_t attr) {
	return ((EOS_SessionDetails_Attribute*)attr)->Data->Value.AsBool == EOS_TRUE;
}
const char* eos_session_attr_get_string(uintptr_t attr) {
	return safe_str(((EOS_SessionDetails_Attribute*)attr)->Data->Value.AsUtf8);
}
int eos_session_attr_get_advertisement_type(uintptr_t attr) {
	return (int)((EOS_SessionDetails_Attribute*)attr)->AdvertisementType;
}
void eos_session_attr_release(uintptr_t attr) {
	EOS_SessionDetails_Attribute_Release((EOS_SessionDetails_Attribute*)attr);
}
void eos_session_details_release(uintptr_t details) {
	EOS_SessionDetails_Release((EOS_HSessionDetails)details);
}

/* Active session */

int eos_sessions_copy_active_session_handle(uintptr_t handle, const char* sessionName,
											uintptr_t* outActive) {
	EOS_Sessions_CopyActiveSessionHandleOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_COPYACTIVESESSIONHANDLE_API_LATEST;
	opts.SessionName = sessionName;
	EOS_HActiveSession active = NULL;
	EOS_EResult result =
		EOS_Sessions_CopyActiveSessionHandle((EOS_HSessions)handle, &opts, &active);
	*outActive = (uintptr_t)active;
	return (int)result;
}

const char* eos_active_session_get_name(uintptr_t active) {
	EOS_ActiveSession_CopyInfoOptions opts = {0};
	opts.ApiVersion = EOS_ACTIVESESSION_COPYINFO_API_LATEST;
	EOS_ActiveSession_Info* info = NULL;
	if (EOS_ActiveSession_CopyInfo((EOS_HActiveSession)active, &opts, &info) != EOS_Success)
		return "";
	const char* name = safe_str(info->SessionName);
	/* Note: name pointer is valid only while info is alive.
	   Caller should use C.GoString immediately. */
	return name;
}

uintptr_t eos_active_session_get_local_user(uintptr_t active) {
	EOS_ActiveSession_CopyInfoOptions opts = {0};
	opts.ApiVersion = EOS_ACTIVESESSION_COPYINFO_API_LATEST;
	EOS_ActiveSession_Info* info = NULL;
	if (EOS_ActiveSession_CopyInfo((EOS_HActiveSession)active, &opts, &info) != EOS_Success)
		return 0;
	uintptr_t user = (uintptr_t)info->LocalUserId;
	EOS_ActiveSession_Info_Release(info);
	return user;
}

int eos_active_session_get_state(uintptr_t active) {
	EOS_ActiveSession_CopyInfoOptions opts = {0};
	opts.ApiVersion = EOS_ACTIVESESSION_COPYINFO_API_LATEST;
	EOS_ActiveSession_Info* info = NULL;
	if (EOS_ActiveSession_CopyInfo((EOS_HActiveSession)active, &opts, &info) != EOS_Success)
		return 0;
	int state = (int)info->State;
	EOS_ActiveSession_Info_Release(info);
	return state;
}

uint32_t eos_active_session_get_registered_player_count(uintptr_t active) {
	EOS_ActiveSession_GetRegisteredPlayerCountOptions opts = {0};
	opts.ApiVersion = EOS_ACTIVESESSION_GETREGISTEREDPLAYERCOUNT_API_LATEST;
	return EOS_ActiveSession_GetRegisteredPlayerCount((EOS_HActiveSession)active, &opts);
}

uintptr_t eos_active_session_get_registered_player_by_index(uintptr_t active, uint32_t index) {
	EOS_ActiveSession_GetRegisteredPlayerByIndexOptions opts = {0};
	opts.ApiVersion = EOS_ACTIVESESSION_GETREGISTEREDPLAYERBYINDEX_API_LATEST;
	opts.PlayerIndex = index;
	return (uintptr_t)EOS_ActiveSession_GetRegisteredPlayerByIndex((EOS_HActiveSession)active,
																   &opts);
}

void eos_active_session_info_release(uintptr_t info) {
	EOS_ActiveSession_Info_Release((EOS_ActiveSession_Info*)info);
}

void eos_active_session_release(uintptr_t active) {
	EOS_ActiveSession_Release((EOS_HActiveSession)active);
}

/* Notifications */

uint64_t eos_sessions_add_notify_invite_received(uintptr_t handle, uintptr_t clientData) {
	EOS_Sessions_AddNotifySessionInviteReceivedOptions opts = {0};
	opts.ApiVersion = EOS_SESSIONS_ADDNOTIFYSESSIONINVITERECEIVED_API_LATEST;
	return (uint64_t)EOS_Sessions_AddNotifySessionInviteReceived(
		(EOS_HSessions)handle, &opts, (void*)clientData, &sessionsInviteReceivedTrampoline);
}

void eos_sessions_remove_notify_invite_received(uintptr_t handle, uint64_t id) {
	EOS_Sessions_RemoveNotifySessionInviteReceived((EOS_HSessions)handle, (EOS_NotificationId)id);
}

#endif /* EOS_CGO */
