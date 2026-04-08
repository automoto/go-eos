package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_result_error_should_return_formatted_string(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{"known code", 16, "eos: NotFound (16)"},
		{"success code", 0, "eos: Success (0)"},
		{"unknown code", 99999, "eos: Unknown (99999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResult(tt.code)
			assert.Equal(t, tt.expected, r.Error())
		})
	}
}

func Test_result_code_should_return_code_value(t *testing.T) {
	r := NewResult(4)
	assert.Equal(t, 4, r.Code())
}

func Test_errors_is_should_match_same_code(t *testing.T) {
	r := NewResult(16)
	assert.ErrorIs(t, r, ErrNotFound)
}

func Test_errors_is_should_not_match_different_code(t *testing.T) {
	r := NewResult(16)
	assert.NotErrorIs(t, r, ErrInvalidAuth)
}

func Test_errors_as_should_extract_result(t *testing.T) {
	var err error = NewResult(30)
	var r *Result

	assert.True(t, errors.As(err, &r))
	assert.Equal(t, 30, r.Code())
}

func Test_is_success(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected bool
	}{
		{"success code", 0, true},
		{"error code", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResult(tt.code)
			assert.Equal(t, tt.expected, IsSuccess(r))
		})
	}
}
