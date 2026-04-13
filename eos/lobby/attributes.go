package lobby

import "github.com/mydev/go-eos/eos/internal/cbinding"

// AttributeVisibility controls who can see a lobby attribute. Maps to EOS_ELobbyAttributeVisibility.
type AttributeVisibility int

const (
	// VisibilityPublic makes the attribute visible to all users.
	VisibilityPublic AttributeVisibility = 0
	// VisibilityPrivate makes the attribute visible only to lobby members.
	VisibilityPrivate AttributeVisibility = 1
)

// ComparisonOp specifies how search parameters are compared. Maps to EOS_EComparisonOp.
type ComparisonOp int

const (
	// ComparisonEqual matches when values are equal.
	ComparisonEqual ComparisonOp = 0
	// ComparisonNotEqual matches when values are not equal.
	ComparisonNotEqual ComparisonOp = 1
	// ComparisonGreaterThan matches when the attribute is greater than the parameter.
	ComparisonGreaterThan ComparisonOp = 2
	// ComparisonGreaterThanOrEqual matches when the attribute is greater than or equal to the parameter.
	ComparisonGreaterThanOrEqual ComparisonOp = 3
	// ComparisonLessThan matches when the attribute is less than the parameter.
	ComparisonLessThan ComparisonOp = 4
	// ComparisonLessThanOrEqual matches when the attribute is less than or equal to the parameter.
	ComparisonLessThanOrEqual ComparisonOp = 5
	// ComparisonDistance matches based on absolute distance from the parameter value.
	ComparisonDistance ComparisonOp = 6
	// ComparisonAnyOf matches when the attribute equals any value in the parameter set.
	ComparisonAnyOf ComparisonOp = 7
	// ComparisonNotAnyOf matches when the attribute does not equal any value in the parameter set.
	ComparisonNotAnyOf ComparisonOp = 8
	// ComparisonOneOf matches when the parameter equals any value in the attribute set.
	ComparisonOneOf ComparisonOp = 9
	// ComparisonNotOneOf matches when the parameter does not equal any value in the attribute set.
	ComparisonNotOneOf ComparisonOp = 10
	// ComparisonContains matches when the attribute contains the parameter value.
	ComparisonContains ComparisonOp = 11
)

// PermissionLevel controls who can find and join a lobby. Maps to EOS_ELobbyPermissionLevel.
type PermissionLevel int

const (
	// PermissionPublicAdvertised allows the lobby to appear in search results.
	PermissionPublicAdvertised PermissionLevel = 0
	// PermissionJoinViaPresence allows joining only through presence information.
	PermissionJoinViaPresence PermissionLevel = 1
	// PermissionInviteOnly restricts joining to invited users only.
	PermissionInviteOnly PermissionLevel = 2
)

// MemberStatus represents a lobby member's current status. Maps to EOS_ELobbyMemberStatus.
type MemberStatus int

const (
	// MemberJoined indicates the member has joined the lobby.
	MemberJoined MemberStatus = 0
	// MemberLeft indicates the member has voluntarily left the lobby.
	MemberLeft MemberStatus = 1
	// MemberDisconnected indicates the member was disconnected.
	MemberDisconnected MemberStatus = 2
	// MemberKicked indicates the member was kicked from the lobby.
	MemberKicked MemberStatus = 3
	// MemberPromoted indicates the member was promoted to lobby owner.
	MemberPromoted MemberStatus = 4
	// MemberClosed indicates the lobby was closed.
	MemberClosed MemberStatus = 5
)

// Attribute represents a key-value pair attached to a lobby or lobby member.
type Attribute struct {
	Key        string
	Value      any
	Visibility AttributeVisibility
}

func attributeFromCBinding(attr *cbinding.EOS_Lobby_Attribute) Attribute {
	a := Attribute{
		Key:        attr.Key,
		Visibility: AttributeVisibility(attr.Visibility),
	}
	switch attr.ValueType {
	case cbinding.EOS_AT_Int64:
		a.Value = attr.AsInt64
	case cbinding.EOS_AT_Double:
		a.Value = attr.AsDouble
	case cbinding.EOS_AT_Boolean:
		a.Value = attr.AsBool
	case cbinding.EOS_AT_String:
		a.Value = attr.AsString
	}
	return a
}
