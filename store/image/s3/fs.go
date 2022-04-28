// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package s3

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/umputun/remark42/backend/app/store/image"

	"go.uber.org/zap"
)

const (
	stagingLocation   = "staging_"
	permanentLocation = "permanent_"

	tagMTime = "mtime"
)

//Save saves a file
func (b *Backend) Save(id string, img []byte) error {
	_, err := b.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(stagingLocation + id),
		Body:   bytes.NewReader(img),
	})
	return err
}

//Commit sets tags location to permanent
func (b *Backend) Commit(id string) error {
	_, err := b.client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(b.bucket),
		Key:        aws.String(permanentLocation + id),
		CopySource: aws.String(b.bucket + "/" + stagingLocation + id),
	})

	if err != nil {
		return err
	}

	_, err = b.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(stagingLocation + id),
	})

	return err
}

// ResetCleanupTimer resets cleanup timer for the image
func (b *Backend) ResetCleanupTimer(id string) error {
	_, err := b.client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(b.bucket),
		Key:        aws.String(stagingLocation + id),
		CopySource: aws.String(b.bucket + "/" + stagingLocation + id),
	})

	return err
}

//Load downloads a file
func (b *Backend) Load(id string) ([]byte, error) {
	obj, err := b.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(permanentLocation + id),
	})
	if err != nil {

		obj, err = b.client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(b.bucket),
			Key:    aws.String(stagingLocation + id),
		})

		if err != nil {
			return nil, err
		}

	}

	return ioutil.ReadAll(obj.Body)
}

//Cleanup removes old files
func (b *Backend) Cleanup(ttl time.Duration) error {
	toDelete := []types.Object{}

	objs, err := b.getAllObjects(stagingLocation)
	if err != nil {
		return err
	}

	ttl += 100 * time.Millisecond

	limit := time.Now().Add(-ttl)
	for _, obj := range objs {
		if obj.LastModified.Before(limit) {
			toDelete = append(toDelete, obj)
		}
	}

	for _, obj := range toDelete {
		_, err := b.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(b.bucket),
			Key:    obj.Key,
		})
		if err != nil {
			b.logger.Info("Failed to delete object",
				zap.String("key", *obj.Key),
				zap.Error(err),
			)
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// Info returns meta information about storage
func (b *Backend) Info() (image.StoreInfo, error) {
	var ts time.Time

	objs, err := b.getAllObjects(stagingLocation)
	if err != nil {
		return image.StoreInfo{}, err
	}

	for _, obj := range objs {
		if ts.IsZero() || obj.LastModified.Before(ts) {
			ts = *obj.LastModified
		}
	}

	return image.StoreInfo{FirstStagingImageTS: ts}, nil
}

func (b *Backend) getAllObjects(prefix string) ([]types.Object, error) {
	ret := []types.Object{}
	var nextContinuationToken *string
	for {
		status, err := b.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(b.bucket),
			ContinuationToken: nextContinuationToken,
			Prefix:            aws.String(prefix),
		})
		if err != nil {
			return nil, err
		}

		ret = append(ret, status.Contents...)

		if !status.IsTruncated {
			return ret, nil
		}
		nextContinuationToken = status.NextContinuationToken
	}
}
