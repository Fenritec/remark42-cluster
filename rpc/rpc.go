// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package rpc

import (
	"fmt"
	"os"

	"github.com/go-pkgz/jrpc"
	"go.uber.org/zap"

	"github.com/Fenritec/remark42-cluster/store/engine"
	"github.com/Fenritec/remark42-cluster/store/image"
)

type (
	Server interface {
	}

	server struct {
		rpc     *jrpc.Server
		storage image.IBackend
		engine  engine.IBackend
		logger  *zap.Logger
	}

	Options struct {
		Port    int
		Storage image.IBackend
		Logger  *zap.Logger
		Engine  engine.IBackend
	}
)

func New(opts *Options) Server {
	logger := jrpc.LoggerFunc(func(format string, args ...interface{}) {
		opts.Logger.Info("RPC",
			zap.String("event", fmt.Sprintf(format, args...)),
		)
	})

	ret := server{
		rpc: &jrpc.Server{
			API:    "/",
			Logger: logger,
		},
		storage: opts.Storage,
		logger:  opts.Logger,
		engine:  opts.Engine,
	}

	ret.rpc.Group("image", jrpc.HandlersGroup{
		"save_with_id":        ret.ImageSave,
		"reset_cleanup_timer": ret.ImageResetCleanupTimer,
		"load":                ret.ImageLoad,
		"commit":              ret.ImageCommit,
		"cleanup":             ret.ImageCleanup,
		"info":                ret.ImageInfo,
	})

	ret.rpc.Group("store", jrpc.HandlersGroup{
		"create":      ret.EngineCreate,
		"get":         ret.EngineGet,
		"update":      ret.EngineUpdate,
		"find":        ret.EngineFind,
		"info":        ret.EngineInfo,
		"flag":        ret.EngineFlag,
		"list_flags":  ret.EngineListFlags,
		"user_detail": ret.EngineUserDetail,
		"count":       ret.EngineCount,
		"delete":      ret.EngineDelete,
		"close":       ret.EngineClose,
	})

	go func() {
		err := ret.rpc.Run(opts.Port)
		if err != nil {
			ret.logger.Error("RPC error occured",
				zap.Error(err),
			)
			os.Exit(1)
		}
	}()

	return ret
}
