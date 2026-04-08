//go:build !eosstub

package cbinding

/*
#include "ids_wrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

const idMaxLength = 33 // EOS_EPICACCOUNTID_MAX_LENGTH (32) + null terminator

func EOS_EpicAccountId_ToString(id EOS_EpicAccountId) string {
	if id == 0 {
		return ""
	}
	var buf [idMaxLength]C.char
	bufLen := C.int32_t(idMaxLength)
	result := C.eos_epic_account_id_to_string(C.uintptr_t(id), &buf[0], &bufLen)
	if EOS_EResult(result) != EOS_EResult_Success {
		return ""
	}
	return C.GoString(&buf[0])
}

func EOS_EpicAccountId_FromString(s string) EOS_EpicAccountId {
	if s == "" {
		return 0
	}
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return EOS_EpicAccountId(C.eos_epic_account_id_from_string(cs))
}

func EOS_EpicAccountId_IsValid(id EOS_EpicAccountId) bool {
	return C.eos_epic_account_id_is_valid(C.uintptr_t(id)) != 0
}

func EOS_ProductUserId_ToString(id EOS_ProductUserId) string {
	if id == 0 {
		return ""
	}
	var buf [idMaxLength]C.char
	bufLen := C.int32_t(idMaxLength)
	result := C.eos_product_user_id_to_string(C.uintptr_t(id), &buf[0], &bufLen)
	if EOS_EResult(result) != EOS_EResult_Success {
		return ""
	}
	return C.GoString(&buf[0])
}

func EOS_ProductUserId_FromString(s string) EOS_ProductUserId {
	if s == "" {
		return 0
	}
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return EOS_ProductUserId(C.eos_product_user_id_from_string(cs))
}

func EOS_ProductUserId_IsValid(id EOS_ProductUserId) bool {
	return C.eos_product_user_id_is_valid(C.uintptr_t(id)) != 0
}
