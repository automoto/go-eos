package sessions

import (
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type SessionModification struct {
	handle cbinding.EOS_HSessionModification
	worker *threadworker.Worker
}

func (m *SessionModification) SetBucketId(bucketId string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetBucketId(m.handle, bucketId)
	})
}

func (m *SessionModification) SetPermissionLevel(level SessionPermissionLevel) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetPermissionLevel(m.handle, cbinding.EOS_EOnlineSessionPermissionLevel(level))
	})
}

func (m *SessionModification) SetMaxPlayers(max uint32) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetMaxPlayers(m.handle, max)
	})
}

func (m *SessionModification) SetJoinInProgressAllowed(allowed bool) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetJoinInProgressAllowed(m.handle, allowed)
	})
}

func (m *SessionModification) SetInvitesAllowed(allowed bool) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetInvitesAllowed(m.handle, allowed)
	})
}

func (m *SessionModification) SetHostAddress(addr string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_SetHostAddress(m.handle, addr)
	})
}

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

func (m *SessionModification) RemoveAttribute(key string) error {
	return m.syncCall(func() cbinding.EOS_EResult {
		return cbinding.EOS_SessionModification_RemoveAttribute(m.handle, key)
	})
}

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
