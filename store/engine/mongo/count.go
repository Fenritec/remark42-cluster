// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"fmt"

	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type (
	countResp struct {
		Count int `bson:"count"`
	}
)

func (b *backend) Count(req engine.FindRequest) (count int, err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	if req.Locator.URL != "" { // comment's count for post
		cursor, err := b.db.Collection(commentCollection).Aggregate(ctx,
			mongo.Pipeline{
				bson.D{
					{Key: "$match", Value: bson.D{
						{Key: "user.siteid", Value: req.Locator.SiteID},
						{Key: "url", Value: req.Locator.URL},
					}},
				},
				bson.D{
					{Key: "$count", Value: "count"},
				},
			})

		if err != nil {
			b.options.logger.Error("Error getting comments", zap.Error(err))
			return 0, err
		}

		r := []countResp{}
		err = cursor.All(ctx, &r)
		if err != nil {
			b.options.logger.Error("Error getting comments count", zap.Error(err))
			return 0, err
		}

		if len(r) > 0 {
			count = r[0].Count
		}

		return count, nil
	}

	if req.UserID != "" { // comment's count for user
		cursor, err := b.db.Collection(commentCollection).Aggregate(ctx,
			mongo.Pipeline{
				bson.D{
					{Key: "$match", Value: bson.D{
						{Key: "user.siteid", Value: req.Locator.SiteID},
						{Key: "user.id", Value: req.UserID},
					}},
				},
				bson.D{
					{Key: "$count", Value: "count"},
				},
			})

		if err != nil {
			b.options.logger.Error("Error getting comments", zap.Error(err))
			return 0, err
		}

		r := []countResp{}
		err = cursor.All(ctx, &r)
		if err != nil {
			b.options.logger.Error("Error getting comments count", zap.Error(err))
			return 0, err
		}

		if len(r) > 0 {
			count = r[0].Count
		}

		return count, nil
	}

	return 0, fmt.Errorf("invalid count request %+v", req)
}
