package types

import "fmt"

const (
	CodeSuccess              = 0
	CodeNoConnection         = 1
	CodeInvalidCredentials   = 2
	CodeInvalidUser          = 3
	CodeInvalidAuth          = 4
	CodeAccessDenied         = 5
	CodeMissingPermissions   = 6
	CodeTokenNotAccount      = 7
	CodeTooManyRequests      = 8
	CodeAlreadyPending       = 9
	CodeInvalidParameters    = 10
	CodeInvalidRequest       = 11
	CodeIncompatibleVersion  = 13
	CodeNotConfigured        = 14
	CodeAlreadyConfigured    = 15
	CodeNotImplemented       = 16
	CodeCanceled             = 17
	CodeNotFound             = 18
	CodeOperationWillRetry   = 19
	CodeNoChange             = 20
	CodeVersionMismatch      = 21
	CodeLimitExceeded        = 22
	CodeDisabled             = 23
	CodeTimedOut             = 27
	CodeInvalidState         = 38
	CodeNetworkDisconnected  = 41
	CodeUnexpectedError      = 0x7FFFFFFF
)

var codeNames = map[int]string{
	CodeSuccess:             "Success",
	CodeNoConnection:        "NoConnection",
	CodeInvalidCredentials:  "InvalidCredentials",
	CodeInvalidUser:         "InvalidUser",
	CodeInvalidAuth:         "InvalidAuth",
	CodeAccessDenied:        "AccessDenied",
	CodeMissingPermissions:  "MissingPermissions",
	CodeTokenNotAccount:     "Token_Not_Account",
	CodeTooManyRequests:     "TooManyRequests",
	CodeAlreadyPending:      "AlreadyPending",
	CodeInvalidParameters:   "InvalidParameters",
	CodeInvalidRequest:      "InvalidRequest",
	CodeIncompatibleVersion: "IncompatibleVersion",
	CodeNotConfigured:       "NotConfigured",
	CodeAlreadyConfigured:   "AlreadyConfigured",
	CodeNotImplemented:      "NotImplemented",
	CodeCanceled:            "Canceled",
	CodeNotFound:            "NotFound",
	CodeOperationWillRetry:  "OperationWillRetry",
	CodeNoChange:            "NoChange",
	CodeVersionMismatch:     "VersionMismatch",
	CodeLimitExceeded:       "LimitExceeded",
	CodeDisabled:            "Disabled",
	CodeTimedOut:            "TimedOut",
	CodeInvalidState:        "InvalidState",
	CodeNetworkDisconnected: "NetworkDisconnected",
	CodeUnexpectedError:     "UnexpectedError",
}

var (
	ErrNoConnection        = NewResult(CodeNoConnection)
	ErrInvalidCredentials  = NewResult(CodeInvalidCredentials)
	ErrInvalidUser         = NewResult(CodeInvalidUser)
	ErrInvalidAuth         = NewResult(CodeInvalidAuth)
	ErrAccessDenied        = NewResult(CodeAccessDenied)
	ErrTooManyRequests     = NewResult(CodeTooManyRequests)
	ErrAlreadyPending      = NewResult(CodeAlreadyPending)
	ErrInvalidParameters   = NewResult(CodeInvalidParameters)
	ErrNotConfigured       = NewResult(CodeNotConfigured)
	ErrAlreadyConfigured   = NewResult(CodeAlreadyConfigured)
	ErrCanceled            = NewResult(CodeCanceled)
	ErrNotFound            = NewResult(CodeNotFound)
	ErrTimedOut            = NewResult(CodeTimedOut)
	ErrNetworkDisconnected = NewResult(CodeNetworkDisconnected)
	ErrUnexpectedError     = NewResult(CodeUnexpectedError)
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
	return r.code == CodeSuccess
}
