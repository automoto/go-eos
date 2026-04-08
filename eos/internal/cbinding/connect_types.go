package cbinding

type EOS_Connect_LoginOptions struct {
	CredentialType EOS_EExternalCredentialType
	Token          string
	DisplayName    string
}

type EOS_Connect_CreateUserOptions struct {
	ContinuanceToken EOS_ContinuanceToken
}

type EOS_Connect_LinkAccountOptions struct {
	LocalUserId      EOS_ProductUserId
	ContinuanceToken EOS_ContinuanceToken
}

type EOS_Connect_LoginCallbackInfo struct {
	ResultCode       EOS_EResult
	LocalUserId      EOS_ProductUserId
	ContinuanceToken EOS_ContinuanceToken
}

type EOS_Connect_CreateUserCallbackInfo struct {
	ResultCode  EOS_EResult
	LocalUserId EOS_ProductUserId
}

type EOS_Connect_LinkAccountCallbackInfo struct {
	ResultCode  EOS_EResult
	LocalUserId EOS_ProductUserId
}

type EOS_Connect_AuthExpirationCallbackInfo struct {
	LocalUserId EOS_ProductUserId
}

type EOS_Connect_LoginStatusChangedCallbackInfo struct {
	LocalUserId    EOS_ProductUserId
	PreviousStatus EOS_ELoginStatus
	CurrentStatus  EOS_ELoginStatus
}
