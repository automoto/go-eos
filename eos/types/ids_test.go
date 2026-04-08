package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_epic_account_id_is_valid(t *testing.T) {
	tests := []struct {
		name  string
		id    EpicAccountId
		valid bool
	}{
		{"non-empty is valid", EpicAccountId("abc123"), true},
		{"empty is invalid", EpicAccountId(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.id.IsValid())
		})
	}
}

func Test_epic_account_id_string_should_return_value(t *testing.T) {
	assert.Equal(t, "abc123", EpicAccountId("abc123").String())
}

func Test_product_user_id_is_valid(t *testing.T) {
	tests := []struct {
		name  string
		id    ProductUserId
		valid bool
	}{
		{"non-empty is valid", ProductUserId("user456"), true},
		{"empty is invalid", ProductUserId(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.id.IsValid())
		})
	}
}

func Test_product_user_id_string_should_return_value(t *testing.T) {
	assert.Equal(t, "user456", ProductUserId("user456").String())
}
