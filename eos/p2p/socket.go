package p2p

// SocketId is the application-defined name for a P2P connection. The same
// Name on both peers identifies a logical channel; different Names can share
// the same underlying physical socket.
//
// The SDK's wire format reserves 33 bytes (32 chars + NUL terminator); the
// allowed character set is alphanumeric plus '-_+= .'.
type SocketId struct {
	Name string
}

// Validate checks the Name length. Returns ErrInvalidSocketId if Name is
// empty or longer than 32 characters. Character-set validation is left to
// the SDK.
func (s SocketId) Validate() error {
	if len(s.Name) == 0 || len(s.Name) > 32 {
		return ErrInvalidSocketId
	}
	return nil
}
