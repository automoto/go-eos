//go:build eosstub

package connect_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

func setupConnect(t *testing.T) (*connect.Connect, func()) {
	t.Helper()
	w := threadworker.New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	c := connect.New(cbinding.EOS_HConnect(1), w)
	return c, func() { cancel(); w.Stop() }
}

func Test_login_should_succeed(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	result, err := c.Login(context.Background(), connect.LoginOptions{
		CredentialType: types.ExternalCredentialEpicIDToken,
		Token:          "test-token",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.LocalUserId)
}

func Test_create_user_should_succeed(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	userId, err := c.CreateUser(context.Background(), 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, *userId)
}

func Test_link_account_should_succeed(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	err := c.LinkAccount(context.Background(), connect.LinkAccountOptions{})
	assert.NoError(t, err)
}

func Test_get_logged_in_users_count_should_return_zero_for_stub(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	assert.Equal(t, 0, c.GetLoggedInUsersCount())
}

func Test_get_logged_in_user_by_index_should_return_empty_for_stub(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	assert.Empty(t, c.GetLoggedInUserByIndex(0))
}

func Test_add_notify_auth_expiration_should_return_remove_func(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	remove := c.AddNotifyAuthExpiration(func(info connect.AuthExpirationInfo) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_add_notify_login_status_changed_should_return_remove_func(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	remove := c.AddNotifyLoginStatusChanged(func(info connect.LoginStatusChangedInfo) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_create_device_id_should_succeed(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	err := c.CreateDeviceId(context.Background(), "test-device")
	assert.NoError(t, err)
}

func Test_create_device_id_should_treat_duplicate_as_success(t *testing.T) {
	cbinding.StubCreateDeviceIdResultCode = cbinding.EOS_EResult_DuplicateNotAllowed
	defer func() { cbinding.StubCreateDeviceIdResultCode = cbinding.EOS_EResult_Success }()

	c, cleanup := setupConnect(t)
	defer cleanup()

	err := c.CreateDeviceId(context.Background(), "test-device")
	assert.NoError(t, err)
}

func Test_delete_device_id_should_succeed(t *testing.T) {
	c, cleanup := setupConnect(t)
	defer cleanup()

	err := c.DeleteDeviceId(context.Background())
	assert.NoError(t, err)
}
