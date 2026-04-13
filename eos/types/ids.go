package types

// EpicAccountId represents an Epic Games account identifier (EOS_EpicAccountId).
type EpicAccountId string

// IsValid reports whether the account ID is non-empty.
func (id EpicAccountId) IsValid() bool {
	return id != ""
}

// String returns the account ID as a plain string.
func (id EpicAccountId) String() string {
	return string(id)
}

// ProductUserId represents a product-scoped user identifier (EOS_ProductUserId).
type ProductUserId string

// IsValid reports whether the product user ID is non-empty.
func (id ProductUserId) IsValid() bool {
	return id != ""
}

// String returns the product user ID as a plain string.
func (id ProductUserId) String() string {
	return string(id)
}
