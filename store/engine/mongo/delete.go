// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package mongo

import (
	"fmt"

	"github.com/umputun/remark42/backend/app/store"
	"github.com/umputun/remark42/backend/app/store/engine"
)

func (b *backend) Delete(req engine.DeleteRequest) error {
	switch {
	case req.UserDetail != "": // delete user detail
		return b.deleteUserDetail(req.Locator, req.UserID, req.UserDetail)
	case req.Locator.URL != "" && req.CommentID != "" && req.UserDetail == "": // delete comment
		return b.deleteComment(req.Locator, req.CommentID, req.DeleteMode)
	case req.Locator.SiteID != "" && req.UserID != "" && req.CommentID == "" && req.UserDetail == "": // delete user
		return b.deleteUser(req.Locator.SiteID, req.UserID, req.DeleteMode)
	case req.Locator.SiteID != "" && req.Locator.URL == "" && req.CommentID == "" && req.UserID == "" && req.UserDetail == "": // delete site
		return b.deleteAll(req.Locator.SiteID)
	}

	return fmt.Errorf("invalid delete request %+v", req)
}

func (b *backend) deleteUserDetail(loc store.Locator, userID string, userDetail engine.UserDetail) error {
	req := engine.UserDetailRequest{
		UserID:  userID,
		Locator: loc,
		Detail:  userDetail,
		Update:  "",
	}

	_, err := b.setUserDetail(req)
	return err
}

func (b *backend) deleteComment(locator store.Locator, commentID string, mode store.DeleteMode) error {
	comment, err := b.Get(engine.GetRequest{
		CommentID: commentID,
		Locator:   locator,
	})
	if err != nil {
		return err
	}

	comment.SetDeleted(mode)

	return b.Update(comment)
}

func (b *backend) deleteUser(siteID, userID string, mode store.DeleteMode) error {
	return fmt.Errorf("TODO Not Implemented")
	// // get list of all comments outside of transaction loop
	// posts, err := b.Info(InfoRequest{Locator: store.Locator{SiteID: siteID}})
	// if err != nil {
	// 	return err
	// }
	//
	// type commentInfo struct {
	// 	locator   store.Locator
	// 	commentID string
	// }
	//
	// // get list of commentID for all user's comment
	// comments := []commentInfo{}
	// for _, postInfo := range posts {
	// 	postInfo := postInfo
	// 	err = bdb.View(func(tx *bolt.Tx) error {
	// 		postsBkt := tx.Bucket([]byte(postsBucketName))
	// 		postBkt := postsBkt.Bucket([]byte(postInfo.URL))
	// 		err = postBkt.ForEach(func(postURL []byte, commentVal []byte) error {
	// 			comment := store.Comment{}
	// 			if err = json.Unmarshal(commentVal, &comment); err != nil {
	// 				return fmt.Errorf("failed to unmarshal: %w", err)
	// 			}
	// 			if comment.User.ID == userID {
	// 				comments = append(comments, commentInfo{locator: comment.Locator, commentID: comment.ID})
	// 			}
	// 			return nil
	// 		})
	// 		if err != nil {
	// 			return fmt.Errorf("failed to collect list of comments for deletion from %s: %w", postInfo.URL, err)
	// 		}
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// log.Printf("[DEBUG] comments for removal=%d", len(comments))
	//
	// // delete collected comments
	// for _, ci := range comments {
	// 	if e := b.deleteComment(bdb, ci.locator, ci.commentID, mode); e != nil {
	// 		return fmt.Errorf("failed to delete comment %+v: %w", ci, err)
	// 	}
	// }
	//
	// // delete user bucket in hard mode
	// if mode == store.HardDelete {
	// 	err = bdb.Update(func(tx *bolt.Tx) error {
	// 		usersBkt := tx.Bucket([]byte(userBucketName))
	// 		if usersBkt != nil {
	// 			if e := usersBkt.DeleteBucket([]byte(userID)); e != nil {
	// 				return fmt.Errorf("failed to delete user bucket for %s: %w", userID, err)
	// 			}
	// 		}
	// 		return nil
	// 	})
	//
	// 	if err != nil {
	// 		return fmt.Errorf("can't delete user meta: %w", err)
	// 	}
	// }
	//
	// if len(comments) == 0 {
	// 	return fmt.Errorf("unknown user %s", userID)
	// }
	//
	// return b.deleteUserDetail(bdb, userID, AllUserDetails)
}

func (b *backend) deleteAll(siteID string) error {
	return fmt.Errorf("TODO Not Implemented")
	// // delete all buckets except blocked users
	// toDelete := []string{postsBucketName, lastBucketName, userBucketName, userDetailsBucketName, infoBucketName}
	//
	// // delete top-level buckets
	// err := bdb.Update(func(tx *bolt.Tx) error {
	// 	for _, bktName := range toDelete {
	// 		if e := tx.DeleteBucket([]byte(bktName)); e != nil {
	// 			return fmt.Errorf("failed to delete top level bucket %s: %w", bktName, e)
	// 		}
	// 		if _, e := tx.CreateBucketIfNotExists([]byte(bktName)); e != nil {
	// 			return fmt.Errorf("failed to create top level bucket %s: %w", bktName, e)
	// 		}
	// 	}
	// 	return nil
	// })
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to delete top level buckets from site %s: %w", siteID, err)
	// }
	// return nil
}
