// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"fmt"

	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (b *backend) Create(comment store.Comment) (commentID string, err error) {
	if b.checkFlag(engine.FlagRequest{Locator: comment.Locator, Flag: engine.ReadOnly}) {
		return "", fmt.Errorf("post %s is read-only", comment.Locator.URL)
	}

	ctx, cancel := b.createContext()
	defer cancel()

	_, err = b.db.Collection(commentCollection).UpdateOne(ctx,
		bson.D{
			{Key: "_id", Value: comment.ID},
			{Key: "user.siteid", Value: comment.User.SiteID},
			{Key: "url", Value: comment.Locator.URL},
		},
		bson.M{
			"$setOnInsert": mongoComment{
				Comment: comment,
				URL:     comment.Locator.URL,
			},
		},
		options.Update().SetUpsert(true))

	if err != nil {
		b.options.logger.Error("Error creating comment", zap.Error(err))
		return commentID, err
	}

	return comment.ID, err
}
