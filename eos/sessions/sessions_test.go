//go:build eosstub

package sessions_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/sessions"
)

func setupSessions(t *testing.T) (*sessions.Sessions, func()) {
	t.Helper()
	w := threadworker.New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	s := sessions.New(cbinding.EOS_HSessions(1), w)
	return s, func() { cancel(); w.Stop() }
}

func Test_create_session_modification_should_return_handle(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	mod, err := s.CreateSessionModification(sessions.CreateSessionOptions{
		SessionName: "test-session",
		BucketId:    "test:mode=dm",
		MaxPlayers:  8,
	})
	assert.NoError(t, err)
	assert.NotNil(t, mod)
	mod.Release()
}

func Test_destroy_session_should_succeed(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	err := s.DestroySession(context.Background(), "test-session")
	assert.NoError(t, err)
}

func Test_start_session_should_succeed(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	err := s.StartSession(context.Background(), "test-session")
	assert.NoError(t, err)
}

func Test_end_session_should_succeed(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	err := s.EndSession(context.Background(), "test-session")
	assert.NoError(t, err)
}

func Test_modification_add_attribute_should_succeed(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	mod, err := s.CreateSessionModification(sessions.CreateSessionOptions{
		SessionName: "test",
		BucketId:    "test:mode=dm",
		MaxPlayers:  4,
	})
	assert.NoError(t, err)
	defer mod.Release()

	assert.NoError(t, mod.AddAttribute("score", int64(100), sessions.AdvertisementAdvertise))
	assert.NoError(t, mod.AddAttribute("name", "test", sessions.AdvertisementDontAdvertise))
}

func Test_create_session_search_should_return_handle(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	search, err := s.CreateSessionSearch(10)
	assert.NoError(t, err)
	assert.NotNil(t, search)
	search.Release()
}

func Test_add_notify_invite_received_should_return_remove_func(t *testing.T) {
	s, cleanup := setupSessions(t)
	defer cleanup()

	remove := s.AddNotifySessionInviteReceived(func(info sessions.InviteReceivedInfo) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}
