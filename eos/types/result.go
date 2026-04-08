package types

import "fmt"

const (
	codeSuccess            = 0
	codeNoConnection       = 1
	codeInvalidCredentials = 2
	codeInvalidUser        = 3
	codeInvalidAuth        = 4
	codeNotFound           = 16
	codeAlreadyPending     = 28
	codeTooManyRequests    = 30
	codeUnexpectedError    = 1001
)

var codeNames = map[int]string{
	codeSuccess:            "Success",
	codeNoConnection:       "NoConnection",
	codeInvalidCredentials: "InvalidCredentials",
	codeInvalidUser:        "InvalidUser",
	codeInvalidAuth:        "InvalidAuth",
	codeNotFound:           "NotFound",
	codeAlreadyPending:     "AlreadyPending",
	codeTooManyRequests:    "TooManyRequests",
	codeUnexpectedError:    "UnexpectedError",
}

var (
	ErrNoConnection       = NewResult(codeNoConnection)
	ErrInvalidCredentials = NewResult(codeInvalidCredentials)
	ErrInvalidUser        = NewResult(codeInvalidUser)
	ErrInvalidAuth        = NewResult(codeInvalidAuth)
	ErrNotFound           = NewResult(codeNotFound)
	ErrAlreadyPending     = NewResult(codeAlreadyPending)
	ErrTooManyRequests    = NewResult(codeTooManyRequests)
	ErrUnexpectedError    = NewResult(codeUnexpectedError)
)

type Result struct {
	code int
}

func NewResult(code int) *Result {
	return &Result{code: code}
}

func (r *Result) Code() int {
	return r.code
}

func (r *Result) Error() string {
	name, ok := codeNames[r.code]
	if !ok {
		name = "Unknown"
	}
	return fmt.Sprintf("eos: %s (%d)", name, r.code)
}

func (r *Result) Is(target error) bool {
	t, ok := target.(*Result)
	if !ok {
		return false
	}
	return r.code == t.code
}

func IsSuccess(r *Result) bool {
	return r.code == codeSuccess
}
