// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package image

import (
	"time"

	"github.com/umputun/remark42/backend/app/store/image"
)

//IBackend represents an image backend
type IBackend interface {
	Disconnect() error

	Save(id string, img []byte) error
	Commit(id string) error
	ResetCleanupTimer(id string) error
	Load(id string) ([]byte, error)
	Cleanup(ttl time.Duration) error
	Info() (image.StoreInfo, error)
}
