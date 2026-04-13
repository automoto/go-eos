# Migration Reference: EOS C/C++ SDK to go-eos

Quick-reference mapping from the EOS C/C++ SDK to go-eos equivalents.

## Platform

| C/C++ SDK | go-eos |
|---|---|
| `EOS_Platform_Create` | `platform.Initialize(cfg)` |
| `EOS_Platform_Tick` | Handled automatically (internal 16ms tick loop) |
| `EOS_Platform_Release` | `p.Shutdown()` |

## Auth

| C/C++ SDK | go-eos |
|---|---|
| `EOS_Auth_Login` | `p.Auth().Login(ctx, auth.LoginOptions{...})` |
| `EOS_Auth_Logout` | `p.Auth().Logout(ctx, localUserId)` |
| `EOS_Auth_DeletePersistentAuth` | `p.Auth().DeletePersistentAuth(ctx)` |
| `EOS_Auth_CopyIdToken` | `p.Auth().CopyIdToken(accountId)` |
| `EOS_Auth_CopyUserAuthToken` | `p.Auth().CopyUserAuthToken(localUserId)` |
| `EOS_Auth_GetLoggedInAccountsCount` | `p.Auth().GetLoggedInAccountsCount()` |
| `EOS_Auth_GetLoggedInAccountByIndex` | `p.Auth().GetLoggedInAccountByIndex(index)` |
| `EOS_Auth_AddNotifyLoginStatusChanged` | `p.Auth().AddNotifyLoginStatusChanged(fn)` returns `RemoveNotifyFunc` |

## Connect

| C/C++ SDK | go-eos |
|---|---|
| `EOS_Connect_Login` | `p.Connect().Login(ctx, connect.LoginOptions{...})` |
| `EOS_Connect_CreateUser` | `p.Connect().CreateUser(ctx, continuanceToken)` |
| `EOS_Connect_LinkAccount` | `p.Connect().LinkAccount(ctx, opts)` |
| `EOS_Connect_CreateDeviceId` | `p.Connect().CreateDeviceId(ctx, deviceModel)` |
| `EOS_Connect_DeleteDeviceId` | `p.Connect().DeleteDeviceId(ctx)` |
| `EOS_Connect_AddNotifyAuthExpiration` | `p.Connect().AddNotifyAuthExpiration(fn)` |
| `EOS_Connect_AddNotifyLoginStatusChanged` | `p.Connect().AddNotifyLoginStatusChanged(fn)` |

## Lobby

| C/C++ SDK | go-eos |
|---|---|
| `EOS_Lobby_CreateLobby` | `p.Lobby().CreateLobby(ctx, lobby.CreateLobbyOptions{...})` |
| `EOS_Lobby_JoinLobby` | `p.Lobby().JoinLobby(ctx, userId, details)` |
| `EOS_Lobby_LeaveLobby` | `p.Lobby().LeaveLobby(ctx, userId, lobbyId)` |
| `EOS_Lobby_SendInvite` | `p.Lobby().SendInvite(ctx, userId, lobbyId)` |
| `EOS_Lobby_KickMember` | `p.Lobby().KickMember(ctx, lobbyId, memberId)` |
| `EOS_LobbySearch_Find` | `search.Find(ctx, userId)` |
| `EOS_Lobby_UpdateLobbyModification` | `p.Lobby().UpdateLobbyModification(userId, lobbyId)` |
| `EOS_Lobby_UpdateLobby` | `p.Lobby().UpdateLobby(ctx, mod)` |
| `EOS_Lobby_AddNotifyLobbyUpdateReceived` | `p.Lobby().AddNotifyLobbyUpdateReceived(fn)` |
| `EOS_Lobby_AddNotifyLobbyMemberStatusReceived` | `p.Lobby().AddNotifyLobbyMemberStatusReceived(fn)` |

## Sessions

| C/C++ SDK | go-eos |
|---|---|
| `EOS_Sessions_CreateSessionModification` | `p.Sessions().CreateSession(ctx, opts)` |
| `EOS_Sessions_UpdateSession` | `p.Sessions().UpdateSessionAttributes(ctx, sessionId, attrs)` |
| `EOS_Sessions_DestroySession` | `p.Sessions().EndSession(ctx, sessionId)` |

## P2P

| C/C++ SDK | go-eos |
|---|---|
| `EOS_P2P_SendPacket` | `p.P2P().SendPacket(ctx, opts)` |
| `EOS_P2P_ReceivePacket` | Handled via `p.P2P().AddNotifyIncomingPacketReceived(fn)` |
| `EOS_P2P_AcceptConnection` | `p.P2P().OpenConnection(ctx, opts)` |
| `EOS_P2P_CloseConnection` | `p.P2P().CloseConnection(ctx, opts)` |
| `EOS_P2P_QueryNATType` | `p.P2P().GetNAT(ctx)` |

## Key Differences

| Concern | C/C++ SDK | go-eos |
|---|---|---|
| **Callback model** | Function pointers passed to each SDK call | Context-based blocking calls return results directly; notification callbacks return a `RemoveNotifyFunc` |
| **Error handling** | `EOS_EResult` enum | `error` (with `types.Result` for EOS-specific codes); use `errors.Is`/`errors.As` |
| **Thread safety** | Manual `EOS_Platform_Tick` calls required | Automatic internal tick loop |
| **Memory management** | Manual `EOS_*_Release` calls for some types | Cleanup via `defer` and GC |
| **Notification pattern** | `EOS_*_AddNotify*` returns `NotificationId`, then call `EOS_*_RemoveNotify*` | `AddNotify*` returns a `RemoveNotifyFunc` closure -- call it directly or `defer` it |
