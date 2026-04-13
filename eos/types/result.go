package types

import "fmt"

// Result code constants correspond to EOS_EResult values from the EOS C SDK.
const (
	// CodeSuccess indicates the operation completed successfully.
	CodeSuccess = 0
	// CodeNoConnection indicates no network connection is available.
	CodeNoConnection = 1
	// CodeInvalidCredentials indicates the supplied credentials are invalid.
	CodeInvalidCredentials = 2
	// CodeInvalidUser indicates the specified user does not exist.
	CodeInvalidUser = 3
	// CodeInvalidAuth indicates an authentication failure.
	CodeInvalidAuth = 4
	// CodeAccessDenied indicates the caller lacks permission for the operation.
	CodeAccessDenied = 5
	// CodeMissingPermissions indicates required permissions are missing.
	CodeMissingPermissions = 6
	// CodeTokenNotAccount indicates the token is not associated with an account.
	CodeTokenNotAccount = 7
	// CodeTooManyRequests indicates the request was rate-limited.
	CodeTooManyRequests = 8
	// CodeAlreadyPending indicates an identical operation is already in progress.
	CodeAlreadyPending = 9
	// CodeInvalidParameters indicates one or more parameters are invalid.
	CodeInvalidParameters = 10
	// CodeInvalidRequest indicates the request is malformed.
	CodeInvalidRequest = 11
	// CodeIncompatibleVersion indicates an SDK version mismatch.
	CodeIncompatibleVersion = 13
	// CodeNotConfigured indicates the required feature has not been configured.
	CodeNotConfigured = 14
	// CodeAlreadyConfigured indicates the feature is already configured.
	CodeAlreadyConfigured = 15
	// CodeNotImplemented indicates the operation is not implemented.
	CodeNotImplemented = 16
	// CodeCanceled indicates the operation was canceled.
	CodeCanceled = 17
	// CodeNotFound indicates the requested resource was not found.
	CodeNotFound = 18
	// CodeOperationWillRetry indicates the operation failed but will be retried.
	CodeOperationWillRetry = 19
	// CodeNoChange indicates the operation completed but nothing changed.
	CodeNoChange = 20
	// CodeVersionMismatch indicates a version mismatch between components.
	CodeVersionMismatch = 21
	// CodeLimitExceeded indicates a limit was exceeded.
	CodeLimitExceeded = 22
	// CodeDisabled indicates the feature is disabled.
	CodeDisabled = 23
	// CodeDuplicateNotAllowed indicates a duplicate entry is not permitted.
	CodeDuplicateNotAllowed = 24
	// CodeTimedOut indicates the operation timed out.
	CodeTimedOut = 27
	// CodeInvalidState indicates the operation is invalid for the current state.
	CodeInvalidState = 38
	// CodeNetworkDisconnected indicates the network connection was lost.
	CodeNetworkDisconnected = 41
	// CodeUnexpectedError indicates an unexpected internal error.
	CodeUnexpectedError = 0x7FFFFFFF
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
	CodeDuplicateNotAllowed: "DuplicateNotAllowed",
	CodeTimedOut:            "TimedOut",
	CodeInvalidState:        "InvalidState",
	CodeNetworkDisconnected: "NetworkDisconnected",
	CodeUnexpectedError:     "UnexpectedError",
}

// Sentinel errors for common EOS result codes. Use errors.Is to compare.
var (
	// ErrNoConnection is returned when no network connection is available.
	ErrNoConnection = NewResult(CodeNoConnection)
	// ErrInvalidCredentials is returned when the supplied credentials are invalid.
	ErrInvalidCredentials = NewResult(CodeInvalidCredentials)
	// ErrInvalidUser is returned when the specified user does not exist.
	ErrInvalidUser = NewResult(CodeInvalidUser)
	// ErrInvalidAuth is returned when authentication fails.
	ErrInvalidAuth = NewResult(CodeInvalidAuth)
	// ErrAccessDenied is returned when the caller lacks permission.
	ErrAccessDenied = NewResult(CodeAccessDenied)
	// ErrTooManyRequests is returned when the request is rate-limited.
	ErrTooManyRequests = NewResult(CodeTooManyRequests)
	// ErrAlreadyPending is returned when an identical operation is already in progress.
	ErrAlreadyPending = NewResult(CodeAlreadyPending)
	// ErrInvalidParameters is returned when one or more parameters are invalid.
	ErrInvalidParameters = NewResult(CodeInvalidParameters)
	// ErrNotConfigured is returned when a required feature has not been configured.
	ErrNotConfigured = NewResult(CodeNotConfigured)
	// ErrAlreadyConfigured is returned when the feature is already configured.
	ErrAlreadyConfigured = NewResult(CodeAlreadyConfigured)
	// ErrCanceled is returned when the operation was canceled.
	ErrCanceled = NewResult(CodeCanceled)
	// ErrNotFound is returned when the requested resource was not found.
	ErrNotFound = NewResult(CodeNotFound)
	// ErrTimedOut is returned when the operation exceeds its time limit.
	ErrTimedOut = NewResult(CodeTimedOut)
	// ErrNetworkDisconnected is returned when the network connection is lost.
	ErrNetworkDisconnected = NewResult(CodeNetworkDisconnected)
	// ErrUnexpectedError is returned when an unexpected internal error occurs.
	ErrUnexpectedError = NewResult(CodeUnexpectedError)
)

// Result wraps an EOS SDK result code (EOS_EResult) and implements the error interface.
type Result struct {
	code int
}

// NewResult creates a Result from a numeric EOS result code.
func NewResult(code int) *Result {
	return &Result{code: code}
}

// Code returns the numeric EOS result code.
func (r *Result) Code() int {
	return r.code
}

// Error returns a human-readable string for the result code.
func (r *Result) Error() string {
	name, ok := codeNames[r.code]
	if !ok {
		name = "Unknown"
	}
	return fmt.Sprintf("eos: %s (%d)", name, r.code)
}

// Is reports whether target is a *Result with the same code, enabling errors.Is support.
func (r *Result) Is(target error) bool {
	t, ok := target.(*Result)
	if !ok {
		return false
	}
	return r.code == t.code
}

// IsSuccess reports whether the result code indicates success.
func IsSuccess(r *Result) bool {
	return r.code == CodeSuccess
}
