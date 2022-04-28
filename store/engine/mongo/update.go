// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"fmt"

	"github.com/umputun/remark42/backend/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (b *backend) Update(comment store.Comment) (err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	res, err := b.db.Collection(commentCollection).UpdateOne(ctx,
		bson.D{
			{Key: "_id", Value: comment.ID},
			{Key: "user.siteid", Value: comment.User.SiteID},
			{Key: "url", Value: comment.Locator.URL},
		},
		bson.M{
			"$set": mongoComment{
				Comment: comment,
				URL:     comment.Locator.URL,
			},
		},
		options.Update())

	if err != nil {
		b.options.logger.Error("Error creating comment", zap.Error(err))
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("No comment matched")
	}

	return err
}
