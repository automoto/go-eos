// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_lobby.h"
#include "cgo_helpers.h"
#include <stdint.h>
#include <string.h>

/* Forward declarations for Go export functions */
extern void goLobbyCreateCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbyDestroyCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbyJoinCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbyLeaveCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbyUpdateCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbySendInviteCallback(int resultCode, uintptr_t clientData, const char* lobbyId);
extern void goLobbyQueryInvitesCallback(int resultCode, uintptr_t clientData,
										uintptr_t localUserId);
extern void goLobbySearchFindCallback(int resultCode, uintptr_t clientData);
extern void goLobbyUpdateReceivedCallback(uintptr_t clientData, const char* lobbyId);
extern void goLobbyMemberUpdateReceivedCallback(uintptr_t clientData, const char* lobbyId,
												uintptr_t targetUserId);
extern void goLobbyMemberStatusReceivedCallback(uintptr_t clientData, const char* lobbyId,
												uintptr_t targetUserId, int currentStatus);
extern void goLobbyInviteReceivedCallback(uintptr_t clientData, const char* inviteId,
										  uintptr_t localUserId, uintptr_t targetUserId);

/* Trampolines — lifecycle */

static void lobbyCreateTrampoline(const EOS_Lobby_CreateLobbyCallbackInfo* data) {
	goLobbyCreateCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						  safe_str(data->LobbyId));
}

static void lobbyDestroyTrampoline(const EOS_Lobby_DestroyLobbyCallbackInfo* data) {
	goLobbyDestroyCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						   safe_str(data->LobbyId));
}

static void lobbyJoinTrampoline(const EOS_Lobby_JoinLobbyCallbackInfo* data) {
	goLobbyJoinCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						safe_str(data->LobbyId));
}

static void lobbyLeaveTrampoline(const EOS_Lobby_LeaveLobbyCallbackInfo* data) {
	goLobbyLeaveCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						 safe_str(data->LobbyId));
}

static void lobbyUpdateTrampoline(const EOS_Lobby_UpdateLobbyCallbackInfo* data) {
	goLobbyUpdateCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						  safe_str(data->LobbyId));
}

static void lobbySendInviteTrampoline(const EOS_Lobby_SendInviteCallbackInfo* data) {
	goLobbySendInviteCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
							  safe_str(data->LobbyId));
}

static void lobbyQueryInvitesTrampoline(const EOS_Lobby_QueryInvitesCallbackInfo* data) {
	goLobbyQueryInvitesCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
								(uintptr_t)data->LocalUserId);
}

static void lobbySearchFindTrampoline(const EOS_LobbySearch_FindCallbackInfo* data) {
	goLobbySearchFindCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

/* Trampolines — notifications */

static void lobbyUpdateReceivedTrampoline(const EOS_Lobby_LobbyUpdateReceivedCallbackInfo* data) {
	goLobbyUpdateReceivedCallback((uintptr_t)data->ClientData, safe_str(data->LobbyId));
}

static void
lobbyMemberUpdateReceivedTrampoline(const EOS_Lobby_LobbyMemberUpdateReceivedCallbackInfo* data) {
	goLobbyMemberUpdateReceivedCallback((uintptr_t)data->ClientData, safe_str(data->LobbyId),
										(uintptr_t)data->TargetUserId);
}

static void
lobbyMemberStatusReceivedTrampoline(const EOS_Lobby_LobbyMemberStatusReceivedCallbackInfo* data) {
	goLobbyMemberStatusReceivedCallback((uintptr_t)data->ClientData, safe_str(data->LobbyId),
										(uintptr_t)data->TargetUserId, (int)data->CurrentStatus);
}

static void lobbyInviteReceivedTrampoline(const EOS_Lobby_LobbyInviteReceivedCallbackInfo* data) {
	goLobbyInviteReceivedCallback((uintptr_t)data->ClientData, safe_str(data->InviteId),
								  (uintptr_t)data->LocalUserId, (uintptr_t)data->TargetUserId);
}

/* Core lifecycle wrappers */

void eos_lobby_create(uintptr_t handle, uintptr_t localUserId, uint32_t maxMembers,
					  int permissionLevel, int allowInvites, const char* bucketId,
					  uintptr_t clientData) {
	EOS_Lobby_CreateLobbyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_CREATELOBBY_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.MaxLobbyMembers = maxMembers;
	opts.PermissionLevel = (EOS_ELobbyPermissionLevel)permissionLevel;
	opts.bAllowInvites = allowInvites ? EOS_TRUE : EOS_FALSE;
	opts.BucketId = bucketId;
	EOS_Lobby_CreateLobby((EOS_HLobby)handle, &opts, (void*)clientData, &lobbyCreateTrampoline);
}

void eos_lobby_destroy(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
					   uintptr_t clientData) {
	EOS_Lobby_DestroyLobbyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_DESTROYLOBBY_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.LobbyId = lobbyId;
	EOS_Lobby_DestroyLobby((EOS_HLobby)handle, &opts, (void*)clientData, &lobbyDestroyTrampoline);
}

void eos_lobby_join(uintptr_t handle, uintptr_t detailsHandle, uintptr_t localUserId,
					uintptr_t clientData) {
	EOS_Lobby_JoinLobbyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_JOINLOBBY_API_LATEST;
	opts.LobbyDetailsHandle = (EOS_HLobbyDetails)detailsHandle;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_Lobby_JoinLobby((EOS_HLobby)handle, &opts, (void*)clientData, &lobbyJoinTrampoline);
}

void eos_lobby_leave(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
					 uintptr_t clientData) {
	EOS_Lobby_LeaveLobbyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_LEAVELOBBY_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.LobbyId = lobbyId;
	EOS_Lobby_LeaveLobby((EOS_HLobby)handle, &opts, (void*)clientData, &lobbyLeaveTrampoline);
}

void eos_lobby_update(uintptr_t handle, uintptr_t modHandle, uintptr_t clientData) {
	EOS_Lobby_UpdateLobbyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_UPDATELOBBY_API_LATEST;
	opts.LobbyModificationHandle = (EOS_HLobbyModification)modHandle;
	EOS_Lobby_UpdateLobby((EOS_HLobby)handle, &opts, (void*)clientData, &lobbyUpdateTrampoline);
}

void eos_lobby_send_invite(uintptr_t handle, const char* lobbyId, uintptr_t localUserId,
						   uintptr_t targetUserId, uintptr_t clientData) {
	EOS_Lobby_SendInviteOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_SENDINVITE_API_LATEST;
	opts.LobbyId = lobbyId;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	EOS_Lobby_SendInvite((EOS_HLobby)handle, &opts, (void*)clientData, &lobbySendInviteTrampoline);
}

void eos_lobby_query_invites(uintptr_t handle, uintptr_t localUserId, uintptr_t clientData) {
	EOS_Lobby_QueryInvitesOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_QUERYINVITES_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_Lobby_QueryInvites((EOS_HLobby)handle, &opts, (void*)clientData,
						   &lobbyQueryInvitesTrampoline);
}

/* Modification handle */

int eos_lobby_update_lobby_modification(uintptr_t handle, uintptr_t localUserId,
										const char* lobbyId, uintptr_t* outMod) {
	EOS_Lobby_UpdateLobbyModificationOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_UPDATELOBBYMODIFICATION_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.LobbyId = lobbyId;
	EOS_HLobbyModification mod = NULL;
	EOS_EResult result = EOS_Lobby_UpdateLobbyModification((EOS_HLobby)handle, &opts, &mod);
	*outMod = (uintptr_t)mod;
	return (int)result;
}

int eos_lobby_mod_set_bucket_id(uintptr_t mod, const char* bucketId) {
	EOS_LobbyModification_SetBucketIdOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_SETBUCKETID_API_LATEST;
	opts.BucketId = bucketId;
	return (int)EOS_LobbyModification_SetBucketId((EOS_HLobbyModification)mod, &opts);
}

int eos_lobby_mod_set_permission_level(uintptr_t mod, int level) {
	EOS_LobbyModification_SetPermissionLevelOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_SETPERMISSIONLEVEL_API_LATEST;
	opts.PermissionLevel = (EOS_ELobbyPermissionLevel)level;
	return (int)EOS_LobbyModification_SetPermissionLevel((EOS_HLobbyModification)mod, &opts);
}

int eos_lobby_mod_set_max_members(uintptr_t mod, uint32_t maxMembers) {
	EOS_LobbyModification_SetMaxMembersOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_SETMAXMEMBERS_API_LATEST;
	opts.MaxMembers = maxMembers;
	return (int)EOS_LobbyModification_SetMaxMembers((EOS_HLobbyModification)mod, &opts);
}

int eos_lobby_mod_set_invites_allowed(uintptr_t mod, int allowed) {
	EOS_LobbyModification_SetInvitesAllowedOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_SETINVITESALLOWED_API_LATEST;
	opts.bInvitesAllowed = allowed ? EOS_TRUE : EOS_FALSE;
	return (int)EOS_LobbyModification_SetInvitesAllowed((EOS_HLobbyModification)mod, &opts);
}

static int lobby_mod_add_attr(uintptr_t mod, const char* key, EOS_Lobby_AttributeData* data,
							  int visibility) {
	EOS_LobbyModification_AddAttributeOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_ADDATTRIBUTE_API_LATEST;
	opts.Attribute = data;
	opts.Visibility = (EOS_ELobbyAttributeVisibility)visibility;
	return (int)EOS_LobbyModification_AddAttribute((EOS_HLobbyModification)mod, &opts);
}

int eos_lobby_mod_add_attr_int64(uintptr_t mod, const char* key, int64_t val, int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsInt64 = val;
	data.ValueType = EOS_AT_INT64;
	return lobby_mod_add_attr(mod, key, &data, visibility);
}

int eos_lobby_mod_add_attr_double(uintptr_t mod, const char* key, double val, int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsDouble = val;
	data.ValueType = EOS_AT_DOUBLE;
	return lobby_mod_add_attr(mod, key, &data, visibility);
}

int eos_lobby_mod_add_attr_bool(uintptr_t mod, const char* key, int val, int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsBool = val ? EOS_TRUE : EOS_FALSE;
	data.ValueType = EOS_AT_BOOLEAN;
	return lobby_mod_add_attr(mod, key, &data, visibility);
}

int eos_lobby_mod_add_attr_string(uintptr_t mod, const char* key, const char* val, int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsUtf8 = val;
	data.ValueType = EOS_AT_STRING;
	return lobby_mod_add_attr(mod, key, &data, visibility);
}

int eos_lobby_mod_remove_attr(uintptr_t mod, const char* key) {
	EOS_LobbyModification_RemoveAttributeOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_REMOVEATTRIBUTE_API_LATEST;
	opts.Key = key;
	return (int)EOS_LobbyModification_RemoveAttribute((EOS_HLobbyModification)mod, &opts);
}

static int lobby_mod_add_member_attr(uintptr_t mod, EOS_Lobby_AttributeData* data, int visibility) {
	EOS_LobbyModification_AddMemberAttributeOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_ADDMEMBERATTRIBUTE_API_LATEST;
	opts.Attribute = data;
	opts.Visibility = (EOS_ELobbyAttributeVisibility)visibility;
	return (int)EOS_LobbyModification_AddMemberAttribute((EOS_HLobbyModification)mod, &opts);
}

int eos_lobby_mod_add_member_attr_int64(uintptr_t mod, const char* key, int64_t val,
										int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsInt64 = val;
	data.ValueType = EOS_AT_INT64;
	return lobby_mod_add_member_attr(mod, &data, visibility);
}

int eos_lobby_mod_add_member_attr_double(uintptr_t mod, const char* key, double val,
										 int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsDouble = val;
	data.ValueType = EOS_AT_DOUBLE;
	return lobby_mod_add_member_attr(mod, &data, visibility);
}

int eos_lobby_mod_add_member_attr_bool(uintptr_t mod, const char* key, int val, int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsBool = val ? EOS_TRUE : EOS_FALSE;
	data.ValueType = EOS_AT_BOOLEAN;
	return lobby_mod_add_member_attr(mod, &data, visibility);
}

int eos_lobby_mod_add_member_attr_string(uintptr_t mod, const char* key, const char* val,
										 int visibility) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsUtf8 = val;
	data.ValueType = EOS_AT_STRING;
	return lobby_mod_add_member_attr(mod, &data, visibility);
}

int eos_lobby_mod_remove_member_attr(uintptr_t mod, const char* key) {
	EOS_LobbyModification_RemoveMemberAttributeOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYMODIFICATION_REMOVEMEMBERATTRIBUTE_API_LATEST;
	opts.Key = key;
	return (int)EOS_LobbyModification_RemoveMemberAttribute((EOS_HLobbyModification)mod, &opts);
}

void eos_lobby_mod_release(uintptr_t mod) {
	EOS_LobbyModification_Release((EOS_HLobbyModification)mod);
}

/* Search */

int eos_lobby_create_search(uintptr_t handle, uint32_t maxResults, uintptr_t* outSearch) {
	EOS_Lobby_CreateLobbySearchOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_CREATELOBBYSEARCH_API_LATEST;
	opts.MaxResults = maxResults;
	EOS_HLobbySearch search = NULL;
	EOS_EResult result = EOS_Lobby_CreateLobbySearch((EOS_HLobby)handle, &opts, &search);
	*outSearch = (uintptr_t)search;
	return (int)result;
}

void eos_lobby_search_find(uintptr_t search, uintptr_t localUserId, uintptr_t clientData) {
	EOS_LobbySearch_FindOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_FIND_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_LobbySearch_Find((EOS_HLobbySearch)search, &opts, (void*)clientData,
						 &lobbySearchFindTrampoline);
}

static int lobby_search_set_param(uintptr_t search, EOS_Lobby_AttributeData* data, int op) {
	EOS_LobbySearch_SetParameterOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_SETPARAMETER_API_LATEST;
	opts.Parameter = data;
	opts.ComparisonOp = (EOS_EComparisonOp)op;
	return (int)EOS_LobbySearch_SetParameter((EOS_HLobbySearch)search, &opts);
}

int eos_lobby_search_set_param_int64(uintptr_t search, const char* key, int64_t val, int op) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsInt64 = val;
	data.ValueType = EOS_AT_INT64;
	return lobby_search_set_param(search, &data, op);
}

int eos_lobby_search_set_param_double(uintptr_t search, const char* key, double val, int op) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsDouble = val;
	data.ValueType = EOS_AT_DOUBLE;
	return lobby_search_set_param(search, &data, op);
}

int eos_lobby_search_set_param_bool(uintptr_t search, const char* key, int val, int op) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsBool = val ? EOS_TRUE : EOS_FALSE;
	data.ValueType = EOS_AT_BOOLEAN;
	return lobby_search_set_param(search, &data, op);
}

int eos_lobby_search_set_param_string(uintptr_t search, const char* key, const char* val, int op) {
	EOS_Lobby_AttributeData data = {0};
	data.ApiVersion = EOS_LOBBY_ATTRIBUTEDATA_API_LATEST;
	data.Key = key;
	data.Value.AsUtf8 = val;
	data.ValueType = EOS_AT_STRING;
	return lobby_search_set_param(search, &data, op);
}

int eos_lobby_search_set_lobby_id(uintptr_t search, const char* lobbyId) {
	EOS_LobbySearch_SetLobbyIdOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_SETLOBBYID_API_LATEST;
	opts.LobbyId = lobbyId;
	return (int)EOS_LobbySearch_SetLobbyId((EOS_HLobbySearch)search, &opts);
}

int eos_lobby_search_set_target_user_id(uintptr_t search, uintptr_t targetUserId) {
	EOS_LobbySearch_SetTargetUserIdOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_SETTARGETUSERID_API_LATEST;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	return (int)EOS_LobbySearch_SetTargetUserId((EOS_HLobbySearch)search, &opts);
}

int eos_lobby_search_set_max_results(uintptr_t search, uint32_t maxResults) {
	EOS_LobbySearch_SetMaxResultsOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_SETMAXRESULTS_API_LATEST;
	opts.MaxResults = maxResults;
	return (int)EOS_LobbySearch_SetMaxResults((EOS_HLobbySearch)search, &opts);
}

uint32_t eos_lobby_search_get_search_result_count(uintptr_t search) {
	EOS_LobbySearch_GetSearchResultCountOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_GETSEARCHRESULTCOUNT_API_LATEST;
	return EOS_LobbySearch_GetSearchResultCount((EOS_HLobbySearch)search, &opts);
}

int eos_lobby_search_copy_search_result_by_index(uintptr_t search, uint32_t index,
												 uintptr_t* outDetails) {
	EOS_LobbySearch_CopySearchResultByIndexOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYSEARCH_COPYSEARCHRESULTBYINDEX_API_LATEST;
	opts.LobbyIndex = index;
	EOS_HLobbyDetails details = NULL;
	EOS_EResult result =
		EOS_LobbySearch_CopySearchResultByIndex((EOS_HLobbySearch)search, &opts, &details);
	*outDetails = (uintptr_t)details;
	return (int)result;
}

void eos_lobby_search_release(uintptr_t search) {
	EOS_LobbySearch_Release((EOS_HLobbySearch)search);
}

/* Details */

int eos_lobby_copy_details_handle(uintptr_t handle, uintptr_t localUserId, const char* lobbyId,
								  uintptr_t* outDetails) {
	EOS_Lobby_CopyLobbyDetailsHandleOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_COPYLOBBYDETAILSHANDLE_API_LATEST;
	opts.LobbyId = lobbyId;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	EOS_HLobbyDetails details = NULL;
	EOS_EResult result = EOS_Lobby_CopyLobbyDetailsHandle((EOS_HLobby)handle, &opts, &details);
	*outDetails = (uintptr_t)details;
	return (int)result;
}

uintptr_t eos_lobby_details_get_owner(uintptr_t details) {
	EOS_LobbyDetails_GetLobbyOwnerOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_GETLOBBYOWNER_API_LATEST;
	return (uintptr_t)EOS_LobbyDetails_GetLobbyOwner((EOS_HLobbyDetails)details, &opts);
}

uint32_t eos_lobby_details_get_member_count(uintptr_t details) {
	EOS_LobbyDetails_GetMemberCountOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_GETMEMBERCOUNT_API_LATEST;
	return EOS_LobbyDetails_GetMemberCount((EOS_HLobbyDetails)details, &opts);
}

uintptr_t eos_lobby_details_get_member_by_index(uintptr_t details, uint32_t index) {
	EOS_LobbyDetails_GetMemberByIndexOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_GETMEMBERBYINDEX_API_LATEST;
	opts.MemberIndex = index;
	return (uintptr_t)EOS_LobbyDetails_GetMemberByIndex((EOS_HLobbyDetails)details, &opts);
}

uint32_t eos_lobby_details_get_attribute_count(uintptr_t details) {
	EOS_LobbyDetails_GetAttributeCountOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_GETATTRIBUTECOUNT_API_LATEST;
	return EOS_LobbyDetails_GetAttributeCount((EOS_HLobbyDetails)details, &opts);
}

uint32_t eos_lobby_details_get_member_attribute_count(uintptr_t details, uintptr_t targetUserId) {
	EOS_LobbyDetails_GetMemberAttributeCountOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_GETMEMBERATTRIBUTECOUNT_API_LATEST;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	return EOS_LobbyDetails_GetMemberAttributeCount((EOS_HLobbyDetails)details, &opts);
}

/* Details info accessors */

int eos_lobby_details_copy_info(uintptr_t details, uintptr_t* outInfo) {
	EOS_LobbyDetails_CopyInfoOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_COPYINFO_API_LATEST;
	EOS_LobbyDetails_Info* info = NULL;
	EOS_EResult result = EOS_LobbyDetails_CopyInfo((EOS_HLobbyDetails)details, &opts, &info);
	*outInfo = (uintptr_t)info;
	return (int)result;
}

const char* eos_lobby_info_get_lobby_id(uintptr_t info) {
	return safe_str(((EOS_LobbyDetails_Info*)info)->LobbyId);
}
uintptr_t eos_lobby_info_get_owner(uintptr_t info) {
	return (uintptr_t)((EOS_LobbyDetails_Info*)info)->LobbyOwnerUserId;
}
int eos_lobby_info_get_permission_level(uintptr_t info) {
	return (int)((EOS_LobbyDetails_Info*)info)->PermissionLevel;
}
uint32_t eos_lobby_info_get_available_slots(uintptr_t info) {
	return ((EOS_LobbyDetails_Info*)info)->AvailableSlots;
}
uint32_t eos_lobby_info_get_max_members(uintptr_t info) {
	return ((EOS_LobbyDetails_Info*)info)->MaxMembers;
}
int eos_lobby_info_get_allow_invites(uintptr_t info) {
	return ((EOS_LobbyDetails_Info*)info)->bAllowInvites == EOS_TRUE;
}
const char* eos_lobby_info_get_bucket_id(uintptr_t info) {
	return safe_str(((EOS_LobbyDetails_Info*)info)->BucketId);
}
void eos_lobby_details_info_release(uintptr_t info) {
	EOS_LobbyDetails_Info_Release((EOS_LobbyDetails_Info*)info);
}

/* Attribute accessors */

int eos_lobby_details_copy_attr_by_index(uintptr_t details, uint32_t index, uintptr_t* outAttr) {
	EOS_LobbyDetails_CopyAttributeByIndexOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_COPYATTRIBUTEBYINDEX_API_LATEST;
	opts.AttrIndex = index;
	EOS_Lobby_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_LobbyDetails_CopyAttributeByIndex((EOS_HLobbyDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

int eos_lobby_details_copy_attr_by_key(uintptr_t details, const char* key, uintptr_t* outAttr) {
	EOS_LobbyDetails_CopyAttributeByKeyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_COPYATTRIBUTEBYKEY_API_LATEST;
	opts.AttrKey = key;
	EOS_Lobby_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_LobbyDetails_CopyAttributeByKey((EOS_HLobbyDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

int eos_lobby_details_copy_member_attr_by_index(uintptr_t details, uintptr_t targetUserId,
												uint32_t index, uintptr_t* outAttr) {
	EOS_LobbyDetails_CopyMemberAttributeByIndexOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_COPYMEMBERATTRIBUTEBYINDEX_API_LATEST;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	opts.AttrIndex = index;
	EOS_Lobby_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_LobbyDetails_CopyMemberAttributeByIndex((EOS_HLobbyDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

int eos_lobby_details_copy_member_attr_by_key(uintptr_t details, uintptr_t targetUserId,
											  const char* key, uintptr_t* outAttr) {
	EOS_LobbyDetails_CopyMemberAttributeByKeyOptions opts = {0};
	opts.ApiVersion = EOS_LOBBYDETAILS_COPYMEMBERATTRIBUTEBYKEY_API_LATEST;
	opts.TargetUserId = (EOS_ProductUserId)targetUserId;
	opts.AttrKey = key;
	EOS_Lobby_Attribute* attr = NULL;
	EOS_EResult result =
		EOS_LobbyDetails_CopyMemberAttributeByKey((EOS_HLobbyDetails)details, &opts, &attr);
	*outAttr = (uintptr_t)attr;
	return (int)result;
}

const char* eos_lobby_attr_get_key(uintptr_t attr) {
	return safe_str(((EOS_Lobby_Attribute*)attr)->Data->Key);
}
int eos_lobby_attr_get_type(uintptr_t attr) {
	return (int)((EOS_Lobby_Attribute*)attr)->Data->ValueType;
}
int64_t eos_lobby_attr_get_int64(uintptr_t attr) {
	return ((EOS_Lobby_Attribute*)attr)->Data->Value.AsInt64;
}
double eos_lobby_attr_get_double(uintptr_t attr) {
	return ((EOS_Lobby_Attribute*)attr)->Data->Value.AsDouble;
}
int eos_lobby_attr_get_bool(uintptr_t attr) {
	return ((EOS_Lobby_Attribute*)attr)->Data->Value.AsBool == EOS_TRUE;
}
const char* eos_lobby_attr_get_string(uintptr_t attr) {
	return safe_str(((EOS_Lobby_Attribute*)attr)->Data->Value.AsUtf8);
}
int eos_lobby_attr_get_visibility(uintptr_t attr) {
	return (int)((EOS_Lobby_Attribute*)attr)->Visibility;
}
void eos_lobby_attr_release(uintptr_t attr) {
	EOS_Lobby_Attribute_Release((EOS_Lobby_Attribute*)attr);
}
void eos_lobby_details_release(uintptr_t details) {
	EOS_LobbyDetails_Release((EOS_HLobbyDetails)details);
}

/* Invite helpers */

uint32_t eos_lobby_get_invite_count(uintptr_t handle, uintptr_t localUserId) {
	EOS_Lobby_GetInviteCountOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_GETINVITECOUNT_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	return EOS_Lobby_GetInviteCount((EOS_HLobby)handle, &opts);
}

int eos_lobby_get_invite_id_by_index(uintptr_t handle, uintptr_t localUserId, uint32_t index,
									 char* outBuf, int32_t* outBufLen) {
	EOS_Lobby_GetInviteIdByIndexOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_GETINVITEIDBYINDEX_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.Index = index;
	return (int)EOS_Lobby_GetInviteIdByIndex((EOS_HLobby)handle, &opts, outBuf, outBufLen);
}

int eos_lobby_copy_details_handle_by_invite_id(uintptr_t handle, const char* inviteId,
											   uintptr_t* outDetails) {
	EOS_Lobby_CopyLobbyDetailsHandleByInviteIdOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_COPYLOBBYDETAILSHANDLEBYINVITEID_API_LATEST;
	opts.InviteId = inviteId;
	EOS_HLobbyDetails details = NULL;
	EOS_EResult result =
		EOS_Lobby_CopyLobbyDetailsHandleByInviteId((EOS_HLobby)handle, &opts, &details);
	*outDetails = (uintptr_t)details;
	return (int)result;
}

/* Notifications */

uint64_t eos_lobby_add_notify_update_received(uintptr_t handle, uintptr_t clientData) {
	EOS_Lobby_AddNotifyLobbyUpdateReceivedOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_ADDNOTIFYLOBBYUPDATERECEIVED_API_LATEST;
	return (uint64_t)EOS_Lobby_AddNotifyLobbyUpdateReceived(
		(EOS_HLobby)handle, &opts, (void*)clientData, &lobbyUpdateReceivedTrampoline);
}

void eos_lobby_remove_notify_update_received(uintptr_t handle, uint64_t id) {
	EOS_Lobby_RemoveNotifyLobbyUpdateReceived((EOS_HLobby)handle, (EOS_NotificationId)id);
}

uint64_t eos_lobby_add_notify_member_update_received(uintptr_t handle, uintptr_t clientData) {
	EOS_Lobby_AddNotifyLobbyMemberUpdateReceivedOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_ADDNOTIFYLOBBYMEMBERUPDATERECEIVED_API_LATEST;
	return (uint64_t)EOS_Lobby_AddNotifyLobbyMemberUpdateReceived(
		(EOS_HLobby)handle, &opts, (void*)clientData, &lobbyMemberUpdateReceivedTrampoline);
}

void eos_lobby_remove_notify_member_update_received(uintptr_t handle, uint64_t id) {
	EOS_Lobby_RemoveNotifyLobbyMemberUpdateReceived((EOS_HLobby)handle, (EOS_NotificationId)id);
}

uint64_t eos_lobby_add_notify_member_status_received(uintptr_t handle, uintptr_t clientData) {
	EOS_Lobby_AddNotifyLobbyMemberStatusReceivedOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_ADDNOTIFYLOBBYMEMBERSTATUSRECEIVED_API_LATEST;
	return (uint64_t)EOS_Lobby_AddNotifyLobbyMemberStatusReceived(
		(EOS_HLobby)handle, &opts, (void*)clientData, &lobbyMemberStatusReceivedTrampoline);
}

void eos_lobby_remove_notify_member_status_received(uintptr_t handle, uint64_t id) {
	EOS_Lobby_RemoveNotifyLobbyMemberStatusReceived((EOS_HLobby)handle, (EOS_NotificationId)id);
}

uint64_t eos_lobby_add_notify_invite_received(uintptr_t handle, uintptr_t clientData) {
	EOS_Lobby_AddNotifyLobbyInviteReceivedOptions opts = {0};
	opts.ApiVersion = EOS_LOBBY_ADDNOTIFYLOBBYINVITERECEIVED_API_LATEST;
	return (uint64_t)EOS_Lobby_AddNotifyLobbyInviteReceived(
		(EOS_HLobby)handle, &opts, (void*)clientData, &lobbyInviteReceivedTrampoline);
}

void eos_lobby_remove_notify_invite_received(uintptr_t handle, uint64_t id) {
	EOS_Lobby_RemoveNotifyLobbyInviteReceived((EOS_HLobby)handle, (EOS_NotificationId)id);
}

#endif /* EOS_CGO */
