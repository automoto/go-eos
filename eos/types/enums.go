package types

type LoginType int

const (
	LoginTypeExchangeCode LoginType = iota
	LoginTypePersistentAuth
	LoginTypeDeviceCode
	LoginTypeAccountPortal
	LoginTypeDeveloper
)

type NATType int

const (
	NATTypeUnknown NATType = iota
	NATTypeOpen
	NATTypeModerate
	NATTypeStrict
)

type LogLevel int

const (
	LogOff LogLevel = iota
	LogFatal
	LogError
	LogWarning
	LogInfo
	LogVerbose
	LogVeryVerbose
)
