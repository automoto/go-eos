//go:build eosstub

package cbinding

import "fmt"

func EOS_EpicAccountId_ToString(id EOS_EpicAccountId) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%016x", uint64(id))
}

func EOS_EpicAccountId_FromString(s string) EOS_EpicAccountId {
	if s == "" {
		return 0
	}
	return EOS_EpicAccountId(1)
}

func EOS_EpicAccountId_IsValid(id EOS_EpicAccountId) bool {
	return id != 0
}

func EOS_ProductUserId_ToString(id EOS_ProductUserId) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%016x", uint64(id))
}

func EOS_ProductUserId_FromString(s string) EOS_ProductUserId {
	if s == "" {
		return 0
	}
	return EOS_ProductUserId(1)
}

func EOS_ProductUserId_IsValid(id EOS_ProductUserId) bool {
	return id != 0
}
