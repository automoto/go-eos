package lobby

import (
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// LobbyDetails provides read-only access to a lobby's information and attributes.
type LobbyDetails struct {
	handle cbinding.EOS_HLobbyDetails
	worker *threadworker.Worker
}

// Info returns the lobby's metadata. Wraps EOS_LobbyDetails_CopyInfo.
func (d *LobbyDetails) Info() (*LobbyInfo, error) {
	var info *cbinding.EOS_LobbyDetails_Info
	var result cbinding.EOS_EResult

	var ownerStr string
	if err := d.worker.Submit(func() {
		info, result = cbinding.EOS_LobbyDetails_CopyInfo(d.handle)
		if result == cbinding.EOS_EResult_Success {
			ownerStr = string(cbinding.EOS_ProductUserId_ToString(info.LobbyOwnerUserId))
		}
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &LobbyInfo{
		LobbyId:          info.LobbyId,
		LobbyOwnerUserId: types.ProductUserId(ownerStr),
		PermissionLevel:  PermissionLevel(info.PermissionLevel),
		AvailableSlots:   info.AvailableSlots,
		MaxMembers:       info.MaxMembers,
		AllowInvites:     info.AllowInvites,
		BucketId:         info.BucketId,
	}, nil
}

// GetOwner returns the lobby owner's product user ID. Wraps EOS_LobbyDetails_GetLobbyOwner.
func (d *LobbyDetails) GetOwner() types.ProductUserId {
	var result string
	if err := d.worker.Submit(func() {
		owner := cbinding.EOS_LobbyDetails_GetLobbyOwner(d.handle)
		result = string(cbinding.EOS_ProductUserId_ToString(owner))
	}); err != nil {
		return ""
	}
	return types.ProductUserId(result)
}

// GetMemberCount returns the number of members in the lobby. Wraps EOS_LobbyDetails_GetMemberCount.
func (d *LobbyDetails) GetMemberCount() int {
	var count uint32
	if err := d.worker.Submit(func() {
		count = cbinding.EOS_LobbyDetails_GetMemberCount(d.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

// GetMemberByIndex returns the product user ID of the member at the given index. Wraps EOS_LobbyDetails_GetMemberByIndex.
func (d *LobbyDetails) GetMemberByIndex(index int) types.ProductUserId {
	var result string
	if err := d.worker.Submit(func() {
		member := cbinding.EOS_LobbyDetails_GetMemberByIndex(d.handle, uint32(index))
		result = string(cbinding.EOS_ProductUserId_ToString(member))
	}); err != nil {
		return ""
	}
	return types.ProductUserId(result)
}

// GetAttributeCount returns the number of attributes on the lobby. Wraps EOS_LobbyDetails_GetAttributeCount.
func (d *LobbyDetails) GetAttributeCount() int {
	var count uint32
	if err := d.worker.Submit(func() {
		count = cbinding.EOS_LobbyDetails_GetAttributeCount(d.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

// CopyAttributeByIndex returns the lobby attribute at the given index. Wraps EOS_LobbyDetails_CopyAttributeByIndex.
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

// CopyAttributeByKey returns the lobby attribute with the given key. Wraps EOS_LobbyDetails_CopyAttributeByKey.
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

// Release frees the underlying EOS lobby details handle.
func (d *LobbyDetails) Release() {
	_ = d.worker.Submit(func() { cbinding.EOS_LobbyDetails_Release(d.handle) })
}
