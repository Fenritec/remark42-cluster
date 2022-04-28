// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"time"

	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (b *backend) Find(req engine.FindRequest) (comments []store.Comment, err error) {
	comments = []store.Comment{}

	switch {
	case req.Locator.SiteID != "" && req.Locator.URL != "": // find post comments, i.e. for site and url
		ctx, cancel := b.createContext()
		defer cancel()

		cursor, err := b.db.Collection(commentCollection).Find(ctx, bson.D{
			{Key: "user.siteid", Value: req.Locator.SiteID},
			{Key: "url", Value: req.Locator.URL},
		})

		if err != nil {
			b.options.logger.Error("Error getting comments", zap.Error(err))
			return nil, err
		}

		err = cursor.All(ctx, &comments)
		if err != nil {
			b.options.logger.Error("Error getting comments", zap.Error(err))
			return nil, err
		}
	case req.Locator.SiteID != "" && req.Locator.URL == "" && req.UserID == "": // find last comments for site
		comments, err = b.lastComments(req.Locator.SiteID, req.Limit, req.Since)
	case req.Locator.SiteID != "" && req.UserID != "": // find comments for user
		comments, err = b.userComments(req.Locator.SiteID, req.UserID, req.Limit, req.Skip)
	}

	if err != nil {
		return nil, err
	}
	return engine.SortComments(comments, req.Sort), nil
}

// Last returns up to max last comments for given siteID
func (b *backend) lastComments(siteID string, max int, since time.Time) (comments []store.Comment, err error) {
	comments = []store.Comment{}

	if max > 1000 || max == 0 {
		max = 1000
	}

	ctx, cancel := b.createContext()
	defer cancel()

	cursor, err := b.db.Collection(commentCollection).Find(ctx, bson.D{
		{Key: "user.siteid", Value: siteID},
		{Key: "time", Value: bson.M{"$gte": since}},
	}, options.Find().SetLimit(int64(max)).SetSort(bson.M{"time": -1}))

	if err != nil {
		b.options.logger.Error("Error getting comments", zap.Error(err))
		return nil, err
	}

	err = cursor.All(ctx, &comments)
	if err != nil {
		b.options.logger.Error("Error getting comments", zap.Error(err))
		return nil, err
	}

	return comments, err
}

// userComments extracts all comments for given site and given userID
func (b *backend) userComments(siteID, userID string, limit, skip int) (comments []store.Comment, err error) {
	comments = []store.Comment{}

	if limit == 0 || limit > 500 {
		limit = 500
	}

	ctx, cancel := b.createContext()
	defer cancel()

	cursor, err := b.db.Collection(commentCollection).Find(ctx, bson.D{
		{Key: "user.siteid", Value: siteID},
		{Key: "user.id", Value: userID},
	}, options.Find().SetLimit(int64(limit)).SetSort(bson.M{"time": -1}))

	if err != nil {
		b.options.logger.Error("Error getting comments", zap.Error(err))
		return nil, err
	}

	err = cursor.All(ctx, &comments)
	if err != nil {
		b.options.logger.Error("Error getting comments", zap.Error(err))
		return nil, err
	}

	return comments, err
}
