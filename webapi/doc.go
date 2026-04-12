// Package webapi provides a pure-Go client for the EOS Web API.
//
// This package has no Cgo dependency and builds with CGO_ENABLED=0.
// It covers OAuth2 authentication, leaderboard queries, and entitlement checks.
//
// All methods accept a context.Context for cancellation and deadlines.
// The Client is safe for concurrent use from multiple goroutines.
package webapi
