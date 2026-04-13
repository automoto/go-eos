package sessions

import (
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// SessionModification wraps EOS_HSessionModification for building session changes before applying them with UpdateSession.
type SessionModification struct {
	handle cbinding.EOS_HSessionModification
	worker *threadworker.Worker
}

// SetBucketId sets the bucket ID on the modification. See EOS_SessionModification_SetBucketId.
func (m *SessionModification) SetBucketId(bucketId string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetBucketId(m.handle, bucketId)
	})
}

// SetPermissionLevel sets the permission level on the modification. See EOS_SessionModification_SetPermissionLevel.
func (m *SessionModification) SetPermissionLevel(level SessionPermissionLevel) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetPermissionLevel(m.handle, cbinding.EOS_EOnlineSessionPermissionLevel(level))
	})
}

// SetMaxPlayers sets the maximum player count on the modification. See EOS_SessionModification_SetMaxPlayers.
func (m *SessionModification) SetMaxPlayers(max uint32) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetMaxPlayers(m.handle, max)
	})
}

// SetJoinInProgressAllowed controls whether players can join a session that is already in progress. See EOS_SessionModification_SetJoinInProgressAllowed.
func (m *SessionModification) SetJoinInProgressAllowed(allowed bool) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetJoinInProgressAllowed(m.handle, allowed)
	})
}

// SetInvitesAllowed controls whether invites are permitted for this session. See EOS_SessionModification_SetInvitesAllowed.
func (m *SessionModification) SetInvitesAllowed(allowed bool) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetInvitesAllowed(m.handle, allowed)
	})
}

// SetHostAddress sets the host address on the modification. See EOS_SessionModification_SetHostAddress.
func (m *SessionModification) SetHostAddress(addr string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetHostAddress(m.handle, addr)
	})
}

// AddAttribute adds a typed key-value attribute to the modification. See EOS_SessionModification_AddAttribute.
func (m *SessionModification) AddAttribute(key string, value any, advType AdvertisementType) error {
	cAdvType := cbinding.EOS_ESessionAttributeAdvertisementType(advType)
	var result cbinding.EOS_EResult

	if err := m.worker.Submit(func() {
		switch v := value.(type) {
		case int64:
			result = cbinding.EOS_SessionModification_AddAttributeInt64(m.handle, key, v, cAdvType)
		case float64:
			result = cbinding.EOS_SessionModification_AddAttributeDouble(m.handle, key, v, cAdvType)
		case bool:
			result = cbinding.EOS_SessionModification_AddAttributeBool(m.handle, key, v, cAdvType)
		case string:
			result = cbinding.EOS_SessionModification_AddAttributeString(m.handle, key, v, cAdvType)
		default:
			result = cbinding.EOS_EResult(types.CodeInvalidParameters)
		}
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// RemoveAttribute removes the attribute with the given key from the modification. See EOS_SessionModification_RemoveAttribute.
func (m *SessionModification) RemoveAttribute(key string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_RemoveAttribute(m.handle, key)
	})
}

// Release frees the underlying SDK modification handle. See EOS_SessionModification_Release.
func (m *SessionModification) Release() {
	_ = m.worker.Submit(func() { cbinding.EOS_SessionModification_Release(m.handle) })
}

func (m *SessionModification) syncCall(fn func() cbinding.EOS_EResult) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() { result = fn() }); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}
