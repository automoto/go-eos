package cbinding

// Sessions lifecycle options

type EOS_Sessions_CreateSessionModificationOptions struct {
	SessionName string
	BucketId    string
	MaxPlayers  uint32
	LocalUserId EOS_ProductUserId
}

type EOS_Sessions_DestroySessionOptions struct {
	SessionName string
}

type EOS_Sessions_JoinSessionOptions struct {
	SessionName   string
	SessionHandle EOS_HSessionDetails
	LocalUserId   EOS_ProductUserId
}

type EOS_Sessions_StartSessionOptions struct {
	SessionName string
}

type EOS_Sessions_EndSessionOptions struct {
	SessionName string
}

type EOS_Sessions_RegisterPlayersOptions struct {
	SessionName string
	PlayerIds   []EOS_ProductUserId
}

type EOS_Sessions_UnregisterPlayersOptions struct {
	SessionName string
	PlayerIds   []EOS_ProductUserId
}

// Sessions callback info

type EOS_Sessions_UpdateSessionCallbackInfo struct {
	ResultCode  EOS_EResult
	SessionName string
	SessionId   string
}

type EOS_Sessions_DestroySessionCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Sessions_JoinSessionCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Sessions_StartSessionCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Sessions_EndSessionCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Sessions_RegisterPlayersCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Sessions_UnregisterPlayersCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_SessionSearch_FindCallbackInfo struct {
	ResultCode EOS_EResult
}

// Notification callback info

type EOS_Sessions_SessionInviteReceivedCallbackInfo struct {
	LocalUserId  EOS_ProductUserId
	TargetUserId EOS_ProductUserId
	InviteId     string
}

// Session details info (from CopyInfo)

type EOS_SessionDetails_Info struct {
	SessionId               string
	HostAddress             string
	NumOpenPublicConnections uint32
	OwnerUserId             EOS_ProductUserId
	BucketId                string
	NumPublicConnections    uint32
	AllowJoinInProgress     bool
	PermissionLevel         EOS_EOnlineSessionPermissionLevel
	InvitesAllowed          bool
}

// Active session info

type EOS_ActiveSession_Info struct {
	SessionName string
	LocalUserId EOS_ProductUserId
	State       EOS_EOnlineSessionState
}

// Session attribute (flattened from C union)

type EOS_Sessions_Attribute struct {
	Key               string
	ValueType         EOS_EAttributeType
	AsInt64           int64
	AsDouble          float64
	AsBool            bool
	AsString          string
	AdvertisementType EOS_ESessionAttributeAdvertisementType
}
