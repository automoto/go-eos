package lobby

import "github.com/mydev/go-eos/eos/internal/cbinding"

type AttributeVisibility int

const (
	VisibilityPublic  AttributeVisibility = 0
	VisibilityPrivate AttributeVisibility = 1
)

type ComparisonOp int

const (
	ComparisonEqual              ComparisonOp = 0
	ComparisonNotEqual           ComparisonOp = 1
	ComparisonGreaterThan        ComparisonOp = 2
	ComparisonGreaterThanOrEqual ComparisonOp = 3
	ComparisonLessThan           ComparisonOp = 4
	ComparisonLessThanOrEqual    ComparisonOp = 5
	ComparisonDistance           ComparisonOp = 6
	ComparisonAnyOf              ComparisonOp = 7
	ComparisonNotAnyOf           ComparisonOp = 8
	ComparisonOneOf              ComparisonOp = 9
	ComparisonNotOneOf           ComparisonOp = 10
	ComparisonContains           ComparisonOp = 11
)

type PermissionLevel int

const (
	PermissionPublicAdvertised PermissionLevel = 0
	PermissionJoinViaPresence  PermissionLevel = 1
	PermissionInviteOnly       PermissionLevel = 2
)

type MemberStatus int

const (
	MemberJoined       MemberStatus = 0
	MemberLeft         MemberStatus = 1
	MemberDisconnected MemberStatus = 2
	MemberKicked       MemberStatus = 3
	MemberPromoted     MemberStatus = 4
	MemberClosed       MemberStatus = 5
)

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
