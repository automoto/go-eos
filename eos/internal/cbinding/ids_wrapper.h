#ifndef IDS_WRAPPER_H
#define IDS_WRAPPER_H

#include <stdint.h>

int eos_epic_account_id_to_string(uintptr_t id, char* buf, int32_t* bufLen);
uintptr_t eos_epic_account_id_from_string(const char* s);
int eos_epic_account_id_is_valid(uintptr_t id);
int eos_product_user_id_to_string(uintptr_t id, char* buf, int32_t* bufLen);
uintptr_t eos_product_user_id_from_string(const char* s);
int eos_product_user_id_is_valid(uintptr_t id);

#endif
