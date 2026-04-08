package cbinding

type EOS_Auth_LoginOptions struct {
	CredentialType EOS_ELoginCredentialType
	ID             string
	Token          string
	ScopeFlags     EOS_EAuthScopeFlags
	ExternalType   EOS_EExternalCredentialType
}

type EOS_Auth_LogoutOptions struct {
	LocalUserId EOS_EpicAccountId
}

type EOS_Auth_DeletePersistentAuthOptions struct {
	RefreshToken string
}

type EOS_Auth_LoginCallbackInfo struct {
	ResultCode        EOS_EResult
	LocalUserId       EOS_EpicAccountId
	PinGrantInfo      *EOS_Auth_PinGrantInfo
	ContinuanceToken  EOS_ContinuanceToken
	SelectedAccountId EOS_EpicAccountId
}

type EOS_Auth_LogoutCallbackInfo struct {
	ResultCode  EOS_EResult
	LocalUserId EOS_EpicAccountId
}

type EOS_Auth_DeletePersistentAuthCallbackInfo struct {
	ResultCode EOS_EResult
}

type EOS_Auth_LoginStatusChangedCallbackInfo struct {
	LocalUserId   EOS_EpicAccountId
	PrevStatus    EOS_ELoginStatus
	CurrentStatus EOS_ELoginStatus
}

type EOS_Auth_PinGrantInfo struct {
	UserCode                string
	VerificationURI         string
	ExpiresIn               int32
	VerificationURIComplete string
}

type EOS_Auth_Token struct {
	App              string
	ClientId         string
	AccountId        EOS_EpicAccountId
	AccessToken      string
	ExpiresIn        float64
	ExpiresAt        string
	AuthType         int32
	RefreshToken     string
	RefreshExpiresIn float64
	RefreshExpiresAt string
}
