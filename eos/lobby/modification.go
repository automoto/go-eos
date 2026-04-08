package lobby

import (
	"fmt"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type LobbyModification struct {
	handle cbinding.EOS_HLobbyModification
	worker *threadworker.Worker
}

func (m *LobbyModification) SetBucketId(bucketId string) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_SetBucketId(m.handle, bucketId)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) SetPermissionLevel(level PermissionLevel) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_SetPermissionLevel(m.handle, cbinding.EOS_ELobbyPermissionLevel(level))
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) SetMaxMembers(max uint32) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_SetMaxMembers(m.handle, max)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) SetInvitesAllowed(allowed bool) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_SetInvitesAllowed(m.handle, allowed)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) AddAttribute(key string, value any, visibility AttributeVisibility) error {
	vis := cbinding.EOS_ELobbyAttributeVisibility(visibility)
	var result cbinding.EOS_EResult

	if err := m.worker.Submit(func() {
		switch v := value.(type) {
		case int64:
			result = cbinding.EOS_LobbyModification_AddAttributeInt64(m.handle, key, v, vis)
		case float64:
			result = cbinding.EOS_LobbyModification_AddAttributeDouble(m.handle, key, v, vis)
		case bool:
			result = cbinding.EOS_LobbyModification_AddAttributeBool(m.handle, key, v, vis)
		case string:
			result = cbinding.EOS_LobbyModification_AddAttributeString(m.handle, key, v, vis)
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

func (m *LobbyModification) RemoveAttribute(key string) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_RemoveAttribute(m.handle, key)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) AddMemberAttribute(key string, value any, visibility AttributeVisibility) error {
	vis := cbinding.EOS_ELobbyAttributeVisibility(visibility)
	var result cbinding.EOS_EResult

	if err := m.worker.Submit(func() {
		switch v := value.(type) {
		case int64:
			result = cbinding.EOS_LobbyModification_AddMemberAttributeInt64(m.handle, key, v, vis)
		case float64:
			result = cbinding.EOS_LobbyModification_AddMemberAttributeDouble(m.handle, key, v, vis)
		case bool:
			result = cbinding.EOS_LobbyModification_AddMemberAttributeBool(m.handle, key, v, vis)
		case string:
			result = cbinding.EOS_LobbyModification_AddMemberAttributeString(m.handle, key, v, vis)
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

func (m *LobbyModification) RemoveMemberAttribute(key string) error {
	var result cbinding.EOS_EResult
	if err := m.worker.Submit(func() {
		result = cbinding.EOS_LobbyModification_RemoveMemberAttribute(m.handle, key)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (m *LobbyModification) Release() {
	_ = m.worker.Submit(func() { cbinding.EOS_LobbyModification_Release(m.handle) })
}

func init() { _ = fmt.Sprint }
