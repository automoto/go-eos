package types

type LoginCredentialType int

const (
	LoginCredentialPassword      LoginCredentialType = 0
	LoginCredentialExchangeCode  LoginCredentialType = 1
	LoginCredentialPersistentAuth LoginCredentialType = 2
	LoginCredentialDeviceCode    LoginCredentialType = 3
	LoginCredentialDeveloper     LoginCredentialType = 4
	LoginCredentialRefreshToken  LoginCredentialType = 5
	LoginCredentialAccountPortal LoginCredentialType = 6
	LoginCredentialExternalAuth  LoginCredentialType = 7
)

type LoginStatus int

const (
	LoginStatusNotLoggedIn      LoginStatus = 0
	LoginStatusUsingLocalProfile LoginStatus = 1
	LoginStatusLoggedIn         LoginStatus = 2
)

type NATType int

const (
	NATTypeUnknown  NATType = iota
	NATTypeOpen
	NATTypeModerate
	NATTypeStrict
)

type LogLevel int

const (
	LogOff         LogLevel = 0
	LogFatal       LogLevel = 100
	LogError       LogLevel = 200
	LogWarning     LogLevel = 300
	LogInfo        LogLevel = 400
	LogVerbose     LogLevel = 500
	LogVeryVerbose LogLevel = 600
)

type AuthScopeFlags uint64

const (
	AuthScopeNoFlags          AuthScopeFlags = 0x0
	AuthScopeBasicProfile     AuthScopeFlags = 0x1
	AuthScopeFriendsList      AuthScopeFlags = 0x2
	AuthScopePresence         AuthScopeFlags = 0x4
	AuthScopeFriendsManagement AuthScopeFlags = 0x8
	AuthScopeEmail            AuthScopeFlags = 0x10
	AuthScopeCountry          AuthScopeFlags = 0x20
)

type ExternalCredentialType int

const (
	ExternalCredentialEpic               ExternalCredentialType = 0
	ExternalCredentialSteamAppTicket     ExternalCredentialType = 1
	ExternalCredentialPSNIDToken         ExternalCredentialType = 2
	ExternalCredentialXBLXSTSToken       ExternalCredentialType = 3
	ExternalCredentialDiscordAccessToken ExternalCredentialType = 4
	ExternalCredentialGOGSessionTicket   ExternalCredentialType = 5
	ExternalCredentialNintendoIDToken    ExternalCredentialType = 6
	ExternalCredentialNintendoNSAIDToken ExternalCredentialType = 7
	ExternalCredentialDeviceIDAccessToken ExternalCredentialType = 10
	ExternalCredentialAppleIDToken       ExternalCredentialType = 11
	ExternalCredentialGoogleIDToken      ExternalCredentialType = 12
	ExternalCredentialEpicIDToken        ExternalCredentialType = 16
	ExternalCredentialAmazonAccessToken  ExternalCredentialType = 17
	ExternalCredentialSteamSessionTicket ExternalCredentialType = 18
)

type ExternalAccountType int

const (
	ExternalAccountEpic    ExternalAccountType = 0
	ExternalAccountSteam   ExternalAccountType = 1
	ExternalAccountPSN     ExternalAccountType = 2
	ExternalAccountXBL     ExternalAccountType = 3
	ExternalAccountDiscord ExternalAccountType = 4
	ExternalAccountGOG     ExternalAccountType = 5
	ExternalAccountNintendo ExternalAccountType = 6
	ExternalAccountApple   ExternalAccountType = 9
	ExternalAccountGoogle  ExternalAccountType = 10
	ExternalAccountAmazon  ExternalAccountType = 13
)
