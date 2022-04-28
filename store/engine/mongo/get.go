// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (b *backend) Get(req engine.GetRequest) (comment store.Comment, err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	dbresult := b.db.Collection(commentCollection).FindOne(ctx, bson.D{
		{Key: "_id", Value: req.CommentID},
		{Key: "user.siteid", Value: req.Locator.SiteID},
		{Key: "url", Value: req.Locator.URL},
	})

	err = dbresult.Decode(&comment)
	if err != nil {
		b.options.logger.Error("Error getting comment", zap.Error(err))
	}

	return comment, err
}
