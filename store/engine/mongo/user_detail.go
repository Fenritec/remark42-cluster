// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"fmt"

	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (b *backend) UserDetail(req engine.UserDetailRequest) (result []engine.UserDetailEntry, err error) {
	switch req.Detail {
	case engine.UserEmail, engine.UserTelegram:
		if req.UserID == "" {
			return nil, fmt.Errorf("userid cannot be empty in request for single detail")
		}

		if req.Update == "" { // read detail value, no update requested
			return b.getUserDetail(req)
		}

		return b.setUserDetail(req)
	case engine.AllUserDetails:
		// list of all details returned in case request is a read request
		// (Update is not set) and does not have UserID
		if req.Update == "" && req.UserID == "" { // read list of all details
			return b.listDetails(req.Locator)
		}
		return nil, fmt.Errorf("unsupported request with userdetail all")
	default:
		return nil, fmt.Errorf("unsupported detail %q", req.Detail)
	}
}

// getUserDetail returns UserDetailEntry with requested userDetail (omitting other details)
// as an only element of the slice.
func (b *backend) getUserDetail(req engine.UserDetailRequest) (result []engine.UserDetailEntry, err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	dbresult := b.db.Collection(userCollection).FindOne(ctx,
		bson.D{
			{Key: "_id", Value: req.UserID},
			{Key: "site", Value: req.Locator.SiteID},
		})

	user := User{}
	err = dbresult.Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			b.options.logger.Error("Error getting user details", zap.Error(err))
			return result, err
		}
		return []engine.UserDetailEntry{}, nil
	}

	switch req.Detail {
	case engine.UserEmail:
		user.UserDetailEntry.Telegram = ""
	case engine.UserTelegram:
		user.UserDetailEntry.Email = ""
	}
	result = []engine.UserDetailEntry{user.UserDetailEntry}

	return result, err
}

// setUserDetail sets requested userDetail, returning complete updated UserDetailEntry as an onlyIps
// element of the slice in case of success
func (b *backend) setUserDetail(req engine.UserDetailRequest) (result []engine.UserDetailEntry, err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	update := bson.M{
		"_id":  req.UserID,
		"site": req.Locator.SiteID,
	}
	entry := engine.UserDetailEntry{
		UserID: req.UserID,
	}

	switch req.Detail {
	case engine.UserEmail:
		update["Email"] = req.Update
		entry.Email = req.Update
	case engine.UserTelegram:
		update["Telegram"] = req.Update
		entry.Telegram = req.Update
	default:
		update["Email"] = req.Update
		update["Telegram"] = req.Update
		entry.Telegram = req.Update
		entry.Email = req.Update
	}

	status, err := b.db.Collection(userCollection).UpdateOne(ctx,
		bson.D{
			{Key: "ID", Value: req.UserID},
			{Key: "Site", Value: req.Locator.SiteID},
		},
		update,
		options.Update().SetUpsert(true))

	if err != nil {
		b.options.logger.Error("Error updating user details", zap.Error(err))
		return result, err
	}

	if status.UpsertedCount == 0 {
		return result, nil
	}

	return []engine.UserDetailEntry{entry}, err
}

// listDetails lists all available users details for given site
func (b *backend) listDetails(loc store.Locator) (result []engine.UserDetailEntry, err error) {
	ctx, cancel := b.createContext()
	defer cancel()

	dbresult, err := b.db.Collection(userCollection).Find(ctx,
		bson.M{"site": loc.SiteID},
	)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			b.options.logger.Error("Error getting all users details", zap.Error(err))
			return result, err
		}
		return result, nil
	}

	users := []User{}
	err = dbresult.All(ctx, &users)
	if err != nil {
		b.options.logger.Error("Error getting all users details", zap.Error(err))
		return result, nil
	}

	result = make([]engine.UserDetailEntry, len(users))
	for i, user := range users {
		result[i].UserID = user.ID
		result[i].Email = user.UserDetailEntry.Email
		result[i].Telegram = user.UserDetailEntry.Telegram
	}

	return result, err
}
