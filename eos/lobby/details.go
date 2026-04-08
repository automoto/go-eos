package lobby

import (
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type LobbyDetails struct {
	handle cbinding.EOS_HLobbyDetails
	worker *threadworker.Worker
}

func (d *LobbyDetails) Info() (*LobbyInfo, error) {
	var info *cbinding.EOS_LobbyDetails_Info
	var result cbinding.EOS_EResult

	if err := d.worker.Submit(func() {
		info, result = cbinding.EOS_LobbyDetails_CopyInfo(d.handle)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &LobbyInfo{
		LobbyId:          info.LobbyId,
		LobbyOwnerUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LobbyOwnerUserId)),
		PermissionLevel:  PermissionLevel(info.PermissionLevel),
		AvailableSlots:   info.AvailableSlots,
		MaxMembers:       info.MaxMembers,
		AllowInvites:     info.AllowInvites,
		BucketId:         info.BucketId,
	}, nil
}

func (d *LobbyDetails) GetOwner() types.ProductUserId {
	var owner cbinding.EOS_ProductUserId
	if err := d.worker.Submit(func() {
		owner = cbinding.EOS_LobbyDetails_GetLobbyOwner(d.handle)
	}); err != nil {
		return ""
	}
	return types.ProductUserId(cbinding.EOS_ProductUserId_ToString(owner))
}

func (d *LobbyDetails) GetMemberCount() int {
	var count uint32
	if err := d.worker.Submit(func() {
		count = cbinding.EOS_LobbyDetails_GetMemberCount(d.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

func (d *LobbyDetails) GetMemberByIndex(index int) types.ProductUserId {
	var member cbinding.EOS_ProductUserId
	if err := d.worker.Submit(func() {
		member = cbinding.EOS_LobbyDetails_GetMemberByIndex(d.handle, uint32(index))
	}); err != nil {
		return ""
	}
	return types.ProductUserId(cbinding.EOS_ProductUserId_ToString(member))
}

func (d *LobbyDetails) GetAttributeCount() int {
	var count uint32
	if err := d.worker.Submit(func() {
		count = cbinding.EOS_LobbyDetails_GetAttributeCount(d.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

func (d *LobbyDetails) CopyAttributeByIndex(index int) (*Attribute, error) {
	var attr *cbinding.EOS_Lobby_Attribute
	var result cbinding.EOS_EResult

	if err := d.worker.Submit(func() {
		attr, result = cbinding.EOS_LobbyDetails_CopyAttributeByIndex(d.handle, uint32(index))
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	a := attributeFromCBinding(attr)
	return &a, nil
}

func (d *LobbyDetails) CopyAttributeByKey(key string) (*Attribute, error) {
	var attr *cbinding.EOS_Lobby_Attribute
	var result cbinding.EOS_EResult

	if err := d.worker.Submit(func() {
		attr, result = cbinding.EOS_LobbyDetails_CopyAttributeByKey(d.handle, key)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	a := attributeFromCBinding(attr)
	return &a, nil
}

func (d *LobbyDetails) Release() {
	_ = d.worker.Submit(func() { cbinding.EOS_LobbyDetails_Release(d.handle) })
}
