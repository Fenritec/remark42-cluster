// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package rpc

import (
	"encoding/json"
	"time"

	"github.com/go-pkgz/jrpc"
)

func (s *server) ImageSave(id uint64, params json.RawMessage) (rr jrpc.Response) {
	var args []json.RawMessage
	if err := json.Unmarshal(params, &args); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if len(args) != 2 {
		return jrpc.Response{Error: "Not enought args"}
	}

	var paramID string
	if err := json.Unmarshal(args[0], &paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	var img []byte
	if err := json.Unmarshal(args[1], &img); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.storage.Save(paramID, img); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) ImageCommit(id uint64, params json.RawMessage) (rr jrpc.Response) {
	var paramID string
	if err := json.Unmarshal(params, &paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.storage.Commit(paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) ImageResetCleanupTimer(id uint64, params json.RawMessage) (rr jrpc.Response) {
	var paramID string
	if err := json.Unmarshal(params, &paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.storage.ResetCleanupTimer(paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) ImageLoad(id uint64, params json.RawMessage) (rr jrpc.Response) {
	var paramID string
	if err := json.Unmarshal(params, &paramID); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	data, err := s.storage.Load(paramID)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, data, nil)
}

func (s *server) ImageCleanup(id uint64, params json.RawMessage) (rr jrpc.Response) {
	var ttl time.Duration
	if err := json.Unmarshal(params, &ttl); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.storage.Cleanup(ttl); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) ImageInfo(id uint64, params json.RawMessage) (rr jrpc.Response) {
	info, err := s.storage.Info()
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, info, nil)
}
