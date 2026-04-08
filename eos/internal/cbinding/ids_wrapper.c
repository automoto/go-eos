// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_common.h"
#include <stdint.h>

int eos_epic_account_id_to_string(uintptr_t id, char* buf, int32_t* bufLen) {
	return (int)EOS_EpicAccountId_ToString((EOS_EpicAccountId)id, buf, bufLen);
}

uintptr_t eos_epic_account_id_from_string(const char* s) {
	return (uintptr_t)EOS_EpicAccountId_FromString(s);
}

int eos_epic_account_id_is_valid(uintptr_t id) {
	return (int)EOS_EpicAccountId_IsValid((EOS_EpicAccountId)id);
}

int eos_product_user_id_to_string(uintptr_t id, char* buf, int32_t* bufLen) {
	return (int)EOS_ProductUserId_ToString((EOS_ProductUserId)id, buf, bufLen);
}

uintptr_t eos_product_user_id_from_string(const char* s) {
	return (uintptr_t)EOS_ProductUserId_FromString(s);
}

int eos_product_user_id_is_valid(uintptr_t id) {
	return (int)EOS_ProductUserId_IsValid((EOS_ProductUserId)id);
}

#endif /* EOS_CGO */
