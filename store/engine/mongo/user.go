// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
)

const (
	userCollection = "users"
)

type (
	User struct {
		store.User             `bson:",inline"`
		engine.UserDetailEntry `bson:",inline"`
	}
)
