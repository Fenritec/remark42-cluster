// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package engine

import (
	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
)

type (
	//IBackend represents an engine backend
	IBackend interface {
		Create(comment store.Comment) (commentID string, err error)
		Get(req engine.GetRequest) (comment store.Comment, err error)
		Update(comment store.Comment) error
		Find(req engine.FindRequest) (comments []store.Comment, err error)
		Info(req engine.InfoRequest) (info []store.PostInfo, err error)
		Flag(req engine.FlagRequest) (status bool, err error)
		ListFlags(req engine.FlagRequest) (list []interface{}, err error)
		UserDetail(req engine.UserDetailRequest) (result []engine.UserDetailEntry, err error)
		Count(req engine.FindRequest) (count int, err error)
		Delete(req engine.DeleteRequest) error
		Close() error

		Disconnect() error
	}
)
