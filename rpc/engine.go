package rpc

import (
	"encoding/json"

	"github.com/go-pkgz/jrpc"
	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
)

func (s *server) EngineCreate(id uint64, params json.RawMessage) (rr jrpc.Response) {
	comment := store.Comment{}
	if err := json.Unmarshal(params, &comment); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	commentId, err := s.engine.Create(comment)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, commentId, nil)
}

func (s *server) EngineGet(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.GetRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	comment, err := s.engine.Get(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, comment, nil)
}

func (s *server) EngineUpdate(id uint64, params json.RawMessage) (rr jrpc.Response) {
	comment := store.Comment{}
	if err := json.Unmarshal(params, &comment); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.engine.Update(comment); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) EngineFind(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.FindRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	comments, err := s.engine.Find(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, comments, nil)
}

func (s *server) EngineInfo(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.InfoRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	info, err := s.engine.Info(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, info, nil)
}

func (s *server) EngineFlag(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.FlagRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	status, err := s.engine.Flag(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, status, nil)
}

func (s *server) EngineListFlags(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.FlagRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	list, err := s.engine.ListFlags(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, list, nil)
}

func (s *server) EngineUserDetail(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.UserDetailRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	result, err := s.engine.UserDetail(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, result, nil)
}

func (s *server) EngineCount(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.FindRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	count, err := s.engine.Count(req)
	if err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, count, nil)
}

func (s *server) EngineDelete(id uint64, params json.RawMessage) (rr jrpc.Response) {
	req := engine.DeleteRequest{}
	if err := json.Unmarshal(params, &req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}

	if err := s.engine.Delete(req); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}

func (s *server) EngineClose(id uint64, params json.RawMessage) (rr jrpc.Response) {
	if err := s.engine.Close(); err != nil {
		return jrpc.Response{Error: err.Error()}
	}
	return jrpc.EncodeResponse(id, nil, nil)
}
