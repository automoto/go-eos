package types

// LoginCredentialType identifies the method used for authentication (EOS_ELoginCredentialType).
type LoginCredentialType int

const (
	// LoginCredentialPassword authenticates with email and password.
	LoginCredentialPassword LoginCredentialType = 0
	// LoginCredentialExchangeCode authenticates with a one-time exchange code.
	LoginCredentialExchangeCode LoginCredentialType = 1
	// LoginCredentialPersistentAuth authenticates using a locally persisted token.
	LoginCredentialPersistentAuth LoginCredentialType = 2
	// LoginCredentialDeviceCode authenticates with a device code.
	LoginCredentialDeviceCode LoginCredentialType = 3
	// LoginCredentialDeveloper authenticates via the developer authentication tool.
	LoginCredentialDeveloper LoginCredentialType = 4
	// LoginCredentialRefreshToken authenticates with a refresh token.
	LoginCredentialRefreshToken LoginCredentialType = 5
	// LoginCredentialAccountPortal authenticates through the account portal.
	LoginCredentialAccountPortal LoginCredentialType = 6
	// LoginCredentialExternalAuth authenticates using an external auth provider.
	LoginCredentialExternalAuth LoginCredentialType = 7
)

// LoginStatus represents the current authentication state of a user (EOS_ELoginStatus).
type LoginStatus int

const (
	// LoginStatusNotLoggedIn indicates the user is not logged in.
	LoginStatusNotLoggedIn LoginStatus = 0
	// LoginStatusUsingLocalProfile indicates a local profile is in use but not fully authenticated.
	LoginStatusUsingLocalProfile LoginStatus = 1
	// LoginStatusLoggedIn indicates the user is fully logged in.
	LoginStatusLoggedIn LoginStatus = 2
)

// NATType represents the detected NAT traversal type (EOS_ENATType).
type NATType int

const (
	// NATTypeUnknown indicates the NAT type has not yet been determined.
	NATTypeUnknown NATType = iota
	// NATTypeOpen indicates an open NAT with no restrictions.
	NATTypeOpen
	// NATTypeModerate indicates a moderate NAT with some restrictions.
	NATTypeModerate
	// NATTypeStrict indicates a strict NAT that may limit connectivity.
	NATTypeStrict
)

// LogLevel controls the verbosity of EOS SDK log output (EOS_ELogLevel).
type LogLevel int

const (
	// LogOff disables all log output.
	LogOff LogLevel = 0
	// LogFatal logs only fatal errors.
	LogFatal LogLevel = 100
	// LogError logs errors and fatal errors.
	LogError LogLevel = 200
	// LogWarning logs warnings and above.
	LogWarning LogLevel = 300
	// LogInfo logs informational messages and above.
	LogInfo LogLevel = 400
	// LogVerbose logs verbose debug messages and above.
	LogVerbose LogLevel = 500
	// LogVeryVerbose logs the most detailed debug messages.
	LogVeryVerbose LogLevel = 600
)

// AuthScopeFlags is a bitmask of permission scopes requested during authentication (EOS_EAuthScopeFlags).
type AuthScopeFlags uint64

const (
	// AuthScopeNoFlags requests no additional scopes.
	AuthScopeNoFlags AuthScopeFlags = 0x0
	// AuthScopeBasicProfile requests access to the user's basic profile.
	AuthScopeBasicProfile AuthScopeFlags = 0x1
	// AuthScopeFriendsList requests access to the user's friends list.
	AuthScopeFriendsList AuthScopeFlags = 0x2
	// AuthScopePresence requests access to the user's presence information.
	AuthScopePresence AuthScopeFlags = 0x4
	// AuthScopeFriendsManagement requests permission to manage the user's friends list.
	AuthScopeFriendsManagement AuthScopeFlags = 0x8
	// AuthScopeEmail requests access to the user's email address.
	AuthScopeEmail AuthScopeFlags = 0x10
	// AuthScopeCountry requests access to the user's country.
	AuthScopeCountry AuthScopeFlags = 0x20
)

// ExternalCredentialType identifies the type of external credential used for Connect login (EOS_EExternalCredentialType).
type ExternalCredentialType int

const (
	// ExternalCredentialEpic uses an Epic Games token.
	ExternalCredentialEpic ExternalCredentialType = 0
	// ExternalCredentialSteamAppTicket uses a Steam encrypted app ticket.
	ExternalCredentialSteamAppTicket ExternalCredentialType = 1
	// ExternalCredentialPSNIDToken uses a PlayStation Network ID token.
	ExternalCredentialPSNIDToken ExternalCredentialType = 2
	// ExternalCredentialXBLXSTSToken uses an Xbox Live XSTS token.
	ExternalCredentialXBLXSTSToken ExternalCredentialType = 3
	// ExternalCredentialDiscordAccessToken uses a Discord access token.
	ExternalCredentialDiscordAccessToken ExternalCredentialType = 4
	// ExternalCredentialGOGSessionTicket uses a GOG Galaxy session ticket.
	ExternalCredentialGOGSessionTicket ExternalCredentialType = 5
	// ExternalCredentialNintendoIDToken uses a Nintendo Account ID token.
	ExternalCredentialNintendoIDToken ExternalCredentialType = 6
	// ExternalCredentialNintendoNSAIDToken uses a Nintendo NSA ID token.
	ExternalCredentialNintendoNSAIDToken ExternalCredentialType = 7
	// ExternalCredentialDeviceIDAccessToken uses a Device ID access token.
	ExternalCredentialDeviceIDAccessToken ExternalCredentialType = 10
	// ExternalCredentialAppleIDToken uses an Apple ID token.
	ExternalCredentialAppleIDToken ExternalCredentialType = 11
	// ExternalCredentialGoogleIDToken uses a Google ID token.
	ExternalCredentialGoogleIDToken ExternalCredentialType = 12
	// ExternalCredentialEpicIDToken uses an Epic ID token.
	ExternalCredentialEpicIDToken ExternalCredentialType = 16
	// ExternalCredentialAmazonAccessToken uses an Amazon access token.
	ExternalCredentialAmazonAccessToken ExternalCredentialType = 17
	// ExternalCredentialSteamSessionTicket uses a Steam session ticket.
	ExternalCredentialSteamSessionTicket ExternalCredentialType = 18
)

// ExternalAccountType identifies an external platform account type (EOS_EExternalAccountType).
type ExternalAccountType int

const (
	// ExternalAccountEpic represents an Epic Games account.
	ExternalAccountEpic ExternalAccountType = 0
	// ExternalAccountSteam represents a Steam account.
	ExternalAccountSteam ExternalAccountType = 1
	// ExternalAccountPSN represents a PlayStation Network account.
	ExternalAccountPSN ExternalAccountType = 2
	// ExternalAccountXBL represents an Xbox Live account.
	ExternalAccountXBL ExternalAccountType = 3
	// ExternalAccountDiscord represents a Discord account.
	ExternalAccountDiscord ExternalAccountType = 4
	// ExternalAccountGOG represents a GOG Galaxy account.
	ExternalAccountGOG ExternalAccountType = 5
	// ExternalAccountNintendo represents a Nintendo account.
	ExternalAccountNintendo ExternalAccountType = 6
	// ExternalAccountApple represents an Apple account.
	ExternalAccountApple ExternalAccountType = 9
	// ExternalAccountGoogle represents a Google account.
	ExternalAccountGoogle ExternalAccountType = 10
	// ExternalAccountAmazon represents an Amazon account.
	ExternalAccountAmazon ExternalAccountType = 13
)
