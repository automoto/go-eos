package cbinding

// Lobby lifecycle options

type EOS_Lobby_CreateLobbyOptions struct {
	LocalUserId     EOS_ProductUserId
	MaxLobbyMembers uint32
	PermissionLevel EOS_ELobbyPermissionLevel
	AllowInvites    bool
	BucketId        string
}

type EOS_Lobby_DestroyLobbyOptions struct {
	LocalUserId EOS_ProductUserId
	LobbyId     string
}

type EOS_Lobby_JoinLobbyOptions struct {
	LobbyDetailsHandle EOS_HLobbyDetails
	LocalUserId        EOS_ProductUserId
}

type EOS_Lobby_LeaveLobbyOptions struct {
	LocalUserId EOS_ProductUserId
	LobbyId     string
}

type EOS_Lobby_UpdateLobbyOptions struct {
	LobbyModificationHandle EOS_HLobbyModification
}

type EOS_Lobby_SendInviteOptions struct {
	LobbyId      string
	LocalUserId  EOS_ProductUserId
	TargetUserId EOS_ProductUserId
}

type EOS_Lobby_QueryInvitesOptions struct {
	LocalUserId EOS_ProductUserId
}

// Lobby lifecycle callback info

type EOS_Lobby_CreateLobbyCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_DestroyLobbyCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_JoinLobbyCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_LeaveLobbyCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_UpdateLobbyCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_SendInviteCallbackInfo struct {
	ResultCode EOS_EResult
	LobbyId    string
}

type EOS_Lobby_QueryInvitesCallbackInfo struct {
	ResultCode  EOS_EResult
	LocalUserId EOS_ProductUserId
}

type EOS_LobbySearch_FindCallbackInfo struct {
	ResultCode EOS_EResult
}

// Notification callback info

type EOS_Lobby_LobbyUpdateReceivedCallbackInfo struct {
	LobbyId string
}

type EOS_Lobby_LobbyMemberUpdateReceivedCallbackInfo struct {
	LobbyId      string
	TargetUserId EOS_ProductUserId
}

type EOS_Lobby_LobbyMemberStatusReceivedCallbackInfo struct {
	LobbyId       string
	TargetUserId  EOS_ProductUserId
	CurrentStatus EOS_ELobbyMemberStatus
}

type EOS_Lobby_LobbyInviteReceivedCallbackInfo struct {
	InviteId     string
	LocalUserId  EOS_ProductUserId
	TargetUserId EOS_ProductUserId
}

// Lobby details info (from CopyInfo)

type EOS_LobbyDetails_Info struct {
	LobbyId          string
	LobbyOwnerUserId EOS_ProductUserId
	PermissionLevel  EOS_ELobbyPermissionLevel
	AvailableSlots   uint32
	MaxMembers       uint32
	AllowInvites     bool
	BucketId         string
}

// Attribute (flattened from C union)

type EOS_Lobby_Attribute struct {
	Key        string
	ValueType  EOS_EAttributeType
	AsInt64    int64
	AsDouble   float64
	AsBool     bool
	AsString   string
	Visibility EOS_ELobbyAttributeVisibility
}
