// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"go.uber.org/zap"

	"github.com/Fenritec/remark42-cluster/store/image"
)

const (
	packageName = "storage/s3"
)

//Backend is the storage handle
type Backend struct {
	client *s3.Client
	bucket string

	logger *zap.Logger
}

//StorageOption are S3 required storage options
type StorageOption struct {
	AccessKey   string
	Secret      string
	EndpointURL string
	Region      string
	BucketName  string

	Logger *zap.Logger
}

//Connect return a storage handler
func Connect(opts *StorageOption) (image.IBackend, error) {
	cfg := aws.NewConfig()
	cfg.Region = opts.Region
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     opts.AccessKey,
			SecretAccessKey: opts.Secret,
		}, nil
	})
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           opts.EndpointURL,
			SigningRegion: opts.Region,
		}, nil
	})

	b := &Backend{
		bucket: opts.BucketName,
		client: s3.NewFromConfig(*cfg),

		logger: opts.Logger,
	}

	_, err := b.Info()
	return b, err
}

//Disconnect stop the S3 session
func (b *Backend) Disconnect() error {
	return nil
}
