// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"github.com/umputun/remark42/backend/app/store"
)

const (
	commentCollection = "comments"
)

type (
	mongoComment struct {
		store.Comment `bson:",inline"`
		URL           string
	}
)
