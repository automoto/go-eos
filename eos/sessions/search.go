package sessions

import (
	"context"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type SessionSearch struct {
	handle cbinding.EOS_HSessionSearch
	worker *threadworker.Worker
}

func (s *SessionSearch) SetParameter(key string, value any, op ComparisonOp) error {
	cOp := cbinding.EOS_EComparisonOp(op)
	var result cbinding.EOS_EResult

	if err := s.worker.Submit(func() {
		switch v := value.(type) {
		case int64:
			result = cbinding.EOS_SessionSearch_SetParameterInt64(s.handle, key, v, cOp)
		case float64:
			result = cbinding.EOS_SessionSearch_SetParameterDouble(s.handle, key, v, cOp)
		case bool:
			result = cbinding.EOS_SessionSearch_SetParameterBool(s.handle, key, v, cOp)
		case string:
			result = cbinding.EOS_SessionSearch_SetParameterString(s.handle, key, v, cOp)
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

func (s *SessionSearch) SetSessionId(sessionId string) error {
	var result cbinding.EOS_EResult
	if err := s.worker.Submit(func() {
		result = cbinding.EOS_SessionSearch_SetSessionId(s.handle, sessionId)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

func (s *SessionSearch) Find(ctx context.Context, localUserId types.ProductUserId) ([]*SessionDetails, error) {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	if err := s.worker.Submit(func() {
		cbinding.EOS_SessionSearch_Find(s.handle, cUserId, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}
	info := result.Data.(*cbinding.EOS_SessionSearch_FindCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(info.ResultCode))
	}

	var count uint32
	if err := s.worker.Submit(func() {
		count = cbinding.EOS_SessionSearch_GetSearchResultCount(s.handle)
	}); err != nil {
		return nil, err
	}

	results := make([]*SessionDetails, 0, count)
	for i := uint32(0); i < count; i++ {
		var details cbinding.EOS_HSessionDetails
		var copyResult cbinding.EOS_EResult
		idx := i
		if err := s.worker.Submit(func() {
			details, copyResult = cbinding.EOS_SessionSearch_CopySearchResultByIndex(s.handle, idx)
		}); err != nil {
			return results, err
		}
		if copyResult == cbinding.EOS_EResult_Success {
			results = append(results, &SessionDetails{handle: details, worker: s.worker})
		}
	}
	return results, nil
}

func (s *SessionSearch) Release() {
	_ = s.worker.Submit(func() { cbinding.EOS_SessionSearch_Release(s.handle) })
}

type SessionDetails struct {
	handle cbinding.EOS_HSessionDetails
	worker *threadworker.Worker
}

func (d *SessionDetails) Handle() cbinding.EOS_HSessionDetails {
	return d.handle
}

func (d *SessionDetails) CopyInfo() (*SessionInfo, error) {
	var info *cbinding.EOS_SessionDetails_Info
	var result cbinding.EOS_EResult

	var ownerStr string
	if err := d.worker.Submit(func() {
		info, result = cbinding.EOS_SessionDetails_CopyInfo(d.handle)
		if result == cbinding.EOS_EResult_Success {
			ownerStr = string(cbinding.EOS_ProductUserId_ToString(info.OwnerUserId))
		}
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &SessionInfo{
		SessionId:                info.SessionId,
		HostAddress:              info.HostAddress,
		NumOpenPublicConnections: info.NumOpenPublicConnections,
		OwnerUserId:              types.ProductUserId(ownerStr),
		BucketId:                 info.BucketId,
		NumPublicConnections:     info.NumPublicConnections,
		AllowJoinInProgress:      info.AllowJoinInProgress,
		PermissionLevel:          SessionPermissionLevel(info.PermissionLevel),
		InvitesAllowed:           info.InvitesAllowed,
	}, nil
}

func (d *SessionDetails) GetAttributeCount() int {
	var count uint32
	if err := d.worker.Submit(func() {
		count = cbinding.EOS_SessionDetails_GetSessionAttributeCount(d.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

func (d *SessionDetails) CopyAttributeByIndex(index int) (*Attribute, error) {
	var attr *cbinding.EOS_Sessions_Attribute
	var result cbinding.EOS_EResult

	if err := d.worker.Submit(func() {
		attr, result = cbinding.EOS_SessionDetails_CopySessionAttributeByIndex(d.handle, uint32(index))
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return attrFromCBinding(attr), nil
}

func (d *SessionDetails) Release() {
	_ = d.worker.Submit(func() { cbinding.EOS_SessionDetails_Release(d.handle) })
}

func attrFromCBinding(attr *cbinding.EOS_Sessions_Attribute) *Attribute {
	a := &Attribute{
		Key:               attr.Key,
		AdvertisementType: AdvertisementType(attr.AdvertisementType),
	}
	switch attr.ValueType {
	case cbinding.EOS_AT_Int64:
		a.Value = attr.AsInt64
	case cbinding.EOS_AT_Double:
		a.Value = attr.AsDouble
	case cbinding.EOS_AT_Boolean:
		a.Value = attr.AsBool
	case cbinding.EOS_AT_String:
		a.Value = attr.AsString
	}
	return a
}
