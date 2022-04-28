// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"errors"

	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
)

func (b *backend) Info(req engine.InfoRequest) (info []store.PostInfo, err error) {
	return nil, errors.New("TODO not implemented")
}

func (b *backend) ListFlags(req engine.FlagRequest) (list []interface{}, err error) {
	return nil, errors.New("TODO not implemented")
}

func (b *backend) Close() error {
	return nil
}
