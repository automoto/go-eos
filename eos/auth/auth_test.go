//go:build eosstub

package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/auth"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

func setupAuth(t *testing.T) (*auth.Auth, func()) {
	t.Helper()
	w := threadworker.New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	a := auth.New(cbinding.EOS_HAuth(1), w)
	return a, func() { cancel(); w.Stop() }
}

func Test_login_should_succeed_with_developer_credentials(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	result, err := a.Login(context.Background(), auth.LoginOptions{
		CredentialType: types.LoginCredentialDeveloper,
		ID:             "localhost:6547",
		Token:          "test-cred",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.LocalUserId)
	assert.NotEmpty(t, result.SelectedAccountId)
}

func Test_logout_should_succeed(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	err := a.Logout(context.Background(), "test-account-id")
	assert.NoError(t, err)
}

func Test_get_logged_in_accounts_count_should_return_zero_for_stub(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	count := a.GetLoggedInAccountsCount()
	assert.Equal(t, 0, count)
}

func Test_get_logged_in_account_by_index_should_return_empty_for_stub(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	id := a.GetLoggedInAccountByIndex(0)
	assert.Empty(t, id)
}

func Test_copy_user_auth_token_should_return_not_found_for_stub(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	_, err := a.CopyUserAuthToken("test-account-id")
	assert.Error(t, err)
}

func Test_delete_persistent_auth_should_succeed(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	err := a.DeletePersistentAuth(context.Background())
	assert.NoError(t, err)
}

func Test_add_notify_login_status_changed_should_return_remove_func(t *testing.T) {
	a, cleanup := setupAuth(t)
	defer cleanup()

	remove := a.AddNotifyLoginStatusChanged(func(info auth.LoginStatusChangedInfo) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}
