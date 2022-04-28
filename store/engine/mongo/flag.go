// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"time"

	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	flagCollection = "flags"
)

type (
	flag struct {
		Key   string
		Site  string
		Flag  engine.Flag
		Until *time.Time
	}
)

func (b *backend) Flag(req engine.FlagRequest) (status bool, err error) {
	if req.Update == engine.FlagNonSet { // read flag value, no update requested
		return b.checkFlag(req), nil
	}

	// write flag value
	return b.setFlag(req)
}

func (b *backend) checkFlag(req engine.FlagRequest) (val bool) {
	key := req.Locator.URL
	if req.UserID != "" {
		key = req.UserID
	}

	ctx, cancel := b.createContext()
	defer cancel()

	dbresult := b.db.Collection(flagCollection).FindOne(ctx,
		bson.D{
			{Key: "Key", Value: key},
			{Key: "Site", Value: req.Locator.SiteID},
		})

	f := flag{}
	err := dbresult.Decode(&f)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			b.options.logger.Error("Error getting flag details", zap.Error(err))
			return false
		}
		return false
	}

	if req.Flag == engine.Blocked && f.Until != nil {
		return time.Now().Before(*f.Until)
	}

	return false
}

func (b *backend) setFlag(req engine.FlagRequest) (res bool, err error) {
	key := req.Locator.URL
	if req.UserID != "" {
		key = req.UserID
	}

	ctx, cancel := b.createContext()
	defer cancel()

	match := bson.D{
		{Key: "Key", Value: key},
		{Key: "Site", Value: req.Locator.SiteID},
	}

	switch req.Update {
	case engine.FlagTrue:
		f := flag{
			Key:  key,
			Site: req.Locator.SiteID,
			Flag: req.Flag,
		}

		if req.Flag == engine.Blocked {
			*f.Until = time.Now().AddDate(100, 0, 0) // permanent is 100 year
			if req.TTL > 0 {
				*f.Until = time.Now().Add(req.TTL)
			}
		}

		status, err := b.db.Collection(flagCollection).UpdateOne(ctx,
			match,
			f,
			options.Update().SetUpsert(true))

		if err != nil {
			b.options.logger.Error("Error updating flag details", zap.Error(err))
			return res, err
		}

		if status.UpsertedCount == 0 {
			return res, nil
		}

		res = true

		return res, err
	case engine.FlagFalse:
		dr, err := b.db.Collection(flagCollection).DeleteOne(ctx, match)
		if err != nil {
			b.options.logger.Error("Error deleting flag details", zap.Error(err))
			return res, err
		}
		if dr.DeletedCount == 0 {
			b.options.logger.Warn("No flag found to delete")
		}
	}

	return res, err
}
