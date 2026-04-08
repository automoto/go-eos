//go:build eosstub

package lobby_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/lobby"
)

func setupLobby(t *testing.T) (*lobby.Lobby, func()) {
	t.Helper()
	w := threadworker.New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	l := lobby.New(cbinding.EOS_HLobby(1), w)
	return l, func() { cancel(); w.Stop() }
}

func Test_create_lobby_should_return_lobby_id(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	lobbyId, err := l.CreateLobby(context.Background(), lobby.CreateLobbyOptions{
		MaxMembers:      4,
		PermissionLevel: lobby.PermissionPublicAdvertised,
		BucketId:        "test:mode=dm",
		AllowInvites:    true,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, lobbyId)
}

func Test_destroy_lobby_should_succeed(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	err := l.DestroyLobby(context.Background(), "", "test-lobby")
	assert.NoError(t, err)
}

func Test_leave_lobby_should_succeed(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	err := l.LeaveLobby(context.Background(), "", "test-lobby")
	assert.NoError(t, err)
}

func Test_update_lobby_modification_should_return_handle(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	mod, err := l.UpdateLobbyModification("", "test-lobby")
	assert.NoError(t, err)
	assert.NotNil(t, mod)
	mod.Release()
}

func Test_modification_set_attribute_should_succeed(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	mod, err := l.UpdateLobbyModification("", "test-lobby")
	assert.NoError(t, err)
	defer mod.Release()

	assert.NoError(t, mod.AddAttribute("score", int64(100), lobby.VisibilityPublic))
	assert.NoError(t, mod.AddAttribute("ratio", float64(1.5), lobby.VisibilityPublic))
	assert.NoError(t, mod.AddAttribute("active", true, lobby.VisibilityPublic))
	assert.NoError(t, mod.AddAttribute("name", "test", lobby.VisibilityPrivate))
}

func Test_create_lobby_search_should_return_handle(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	search, err := l.CreateLobbySearch(10)
	assert.NoError(t, err)
	assert.NotNil(t, search)
	search.Release()
}

func Test_search_find_should_return_results(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	search, err := l.CreateLobbySearch(10)
	assert.NoError(t, err)
	defer search.Release()

	results, err := search.Find(context.Background(), "")
	assert.NoError(t, err)
	assert.NotNil(t, results)
}

func Test_add_notify_member_status_should_return_remove_func(t *testing.T) {
	l, cleanup := setupLobby(t)
	defer cleanup()

	remove := l.AddNotifyLobbyMemberStatusReceived(func(info lobby.MemberStatusInfo) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}
