# go-eos Developer Notes

## macOS: Main Thread Requirement

The EOS SDK on macOS uses Apple's networking stack (`CFRunLoop`, `NSURLSession`) which dispatches HTTP completions through the **main thread's run loop**. All SDK calls — init, tick, and shutdown — must happen on the main OS thread.

Use `RunOnMainThread` instead of `Run`:

```go
func init() {
    runtime.LockOSThread() // locks main goroutine to main OS thread
}

func main() {
    // RunOnMainThread drives the tick loop on the calling (main) goroutine.
    // The user callback runs on a separate goroutine.
    if err := platform.RunOnMainThread(ctx, cfg, run); err != nil {
        log.Fatal(err)
    }
}
```

`Run` uses a background worker goroutine which gets a non-main OS thread. HTTP requests will silently hang on macOS. On Linux/Windows either function works — `RunOnMainThread` is harmless there.

**Recommendation:** Always use `RunOnMainThread` for cross-platform code.

## Connect Login: Use ID Token, Not Access Token

For `EOS_Connect_Login` with `ExternalCredentialEpicIDToken`, you must pass a **JWT ID token** (from `CopyIdToken`), not a bearer access token (from `CopyUserAuthToken`).

```go
// Correct: ID token (JWT) for Connect login
idToken, err := p.Auth().CopyIdToken(loginResult.LocalUserId)
connectResult, err := p.Connect().Login(ctx, connect.LoginOptions{
    CredentialType: types.ExternalCredentialEpicIDToken,
    Token:          idToken,
})

// Wrong: access token (bearer) — will fail with "Token type mismatch"
token, err := p.Auth().CopyUserAuthToken(loginResult.LocalUserId)
p.Connect().Login(ctx, connect.LoginOptions{
    CredentialType: types.ExternalCredentialEpicIDToken,
    Token:          token.AccessToken, // bearer, not id_token
})
```

## DevAuth Tool Setup

The Developer Auth Tool is an Electron app bundled with the SDK at `static/EOS-SDK-*/SDK/Tools/`.

1. Extract and launch (see `docs/test.md` for full steps)
2. Note the **port** shown in the UI — set `EOS_DEV_AUTH_HOST=localhost:<port>`
3. Log in and create a credential (e.g. `player1`) — set `EOS_DEV_AUTH_CREDENTIAL=player1`
4. The tool must stay running during testing

Common pitfall: the default port in documentation is `6547` but the tool may listen on a different port. Verify with `lsof -iTCP -sTCP:LISTEN -P | grep DevAuth`.

## SDK Init Ordering

The EOS SDK requires this exact call order (per SDK samples):

1. `EOS_Initialize` — must be first
2. `EOS_Logging_SetCallback` + `EOS_Logging_SetLogLevel` — after init
3. `EOS_Platform_Create` — after logging

Calling logging functions before `EOS_Initialize` is undefined behavior.
