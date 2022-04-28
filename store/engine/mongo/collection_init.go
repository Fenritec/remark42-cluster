// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	collIndexes = map[string]indexInfo{
		"auth_id_index":  indexInfo{key: "auth_id", unique: true},
		"items_id_index": indexInfo{key: "items._id", unique: false},
	}
)

func (b *backend) createCollection(name string) (exists bool) {
	exists = false
	ctx, cancel := b.createContext()
	defer cancel()

	collections, err := b.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		b.options.logger.Error("ListCollectionName", zap.String("coll", name), zap.Error(err))
		return
	}

	//Checking if the collection is already created
	for _, collection := range collections {
		if collection == name {
			exists = true
			break
		}
	}

	if !exists {
		err := b.db.CreateCollection(ctx, name)
		if err != nil {
			b.options.logger.Error("CreateCollection", zap.String("coll", name), zap.Error(err))
			return
		}
	}

	return true
}

func (b *backend) createCollectionIndexes(name string, indexes map[string]indexInfo) {
	ctx, cancel := b.createContext()
	defer cancel()

	indexesCursor, err := b.db.Collection(name).Indexes().List(ctx)

	if err != nil {
		b.options.logger.Error("createSharesCollectionIndexes", zap.String("coll", name), zap.Error(err))
		return
	}

	var results []bson.M
	if err = indexesCursor.All(ctx, &results); err != nil {
		b.options.logger.Error("Getting existing indexes", zap.String("coll", name), zap.Error(err))
		return
	}

	for indexName, indexCol := range indexes {
		exists := false
		for _, existingIndex := range results {
			if existingIndex["name"] == indexName {
				exists = true
				break
			}
		}

		if !exists {
			uniqueOptions := options.Index()
			uniqueOptions.Unique = &indexCol.unique
			uniqueOptions.Name = &indexName
			if _, err := b.db.Collection(name).Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: indexCol.key, Value: 1}},
				Options: uniqueOptions,
			}); err != nil {
				b.options.logger.Error("createSharesCollectionIndexes", zap.String("coll", name), zap.Error(err))
			}
		}
	}
}
