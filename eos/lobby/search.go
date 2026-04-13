package lobby

import (
	"context"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// LobbySearch is a handle used to configure and execute lobby search queries.
type LobbySearch struct {
	handle cbinding.EOS_HLobbySearch
	worker *threadworker.Worker
}

// SetParameter adds a search filter parameter. Wraps EOS_LobbySearch_SetParameter.
func (s *LobbySearch) SetParameter(key string, value any, op ComparisonOp) error {
	cOp := cbinding.EOS_EComparisonOp(op)
	var result cbinding.EOS_EResult

	if err := s.worker.Submit(func() {
		switch v := value.(type) {
		case int64:
			result = cbinding.EOS_LobbySearch_SetParameterInt64(s.handle, key, v, cOp)
		case float64:
			result = cbinding.EOS_LobbySearch_SetParameterDouble(s.handle, key, v, cOp)
		case bool:
			result = cbinding.EOS_LobbySearch_SetParameterBool(s.handle, key, v, cOp)
		case string:
			result = cbinding.EOS_LobbySearch_SetParameterString(s.handle, key, v, cOp)
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

// SetLobbyId filters the search to a specific lobby ID. Wraps EOS_LobbySearch_SetLobbyId.
func (s *LobbySearch) SetLobbyId(lobbyId string) error {
	var result cbinding.EOS_EResult
	if err := s.worker.Submit(func() {
		result = cbinding.EOS_LobbySearch_SetLobbyId(s.handle, lobbyId)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// SetMaxResults limits the number of search results returned. Wraps EOS_LobbySearch_SetMaxResults.
func (s *LobbySearch) SetMaxResults(maxResults uint32) error {
	var result cbinding.EOS_EResult
	if err := s.worker.Submit(func() {
		result = cbinding.EOS_LobbySearch_SetMaxResults(s.handle, maxResults)
	}); err != nil {
		return err
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// Find executes the lobby search and returns matching lobby details. Wraps EOS_LobbySearch_Find.
func (s *LobbySearch) Find(ctx context.Context, localUserId types.ProductUserId) ([]*LobbyDetails, error) {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	if err := s.worker.Submit(func() {
		cbinding.EOS_LobbySearch_Find(s.handle, cUserId, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}
	info := result.Data.(*cbinding.EOS_LobbySearch_FindCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(info.ResultCode))
	}

	var count uint32
	if err := s.worker.Submit(func() {
		count = cbinding.EOS_LobbySearch_GetSearchResultCount(s.handle)
	}); err != nil {
		return nil, err
	}

	results := make([]*LobbyDetails, 0, count)
	for i := uint32(0); i < count; i++ {
		var details cbinding.EOS_HLobbyDetails
		var copyResult cbinding.EOS_EResult
		idx := i
		if err := s.worker.Submit(func() {
			details, copyResult = cbinding.EOS_LobbySearch_CopySearchResultByIndex(s.handle, idx)
		}); err != nil {
			return results, err
		}
		if copyResult == cbinding.EOS_EResult_Success {
			results = append(results, &LobbyDetails{handle: details, worker: s.worker})
		}
	}
	return results, nil
}

// Release frees the underlying EOS lobby search handle.
func (s *LobbySearch) Release() {
	_ = s.worker.Submit(func() { cbinding.EOS_LobbySearch_Release(s.handle) })
}
