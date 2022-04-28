// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package main

import (
	"go.uber.org/zap"

	"github.com/Fenritec/remark42-cluster/config"
	"github.com/Fenritec/remark42-cluster/rpc"
	"github.com/Fenritec/remark42-cluster/store/engine/mongo"
	"github.com/Fenritec/remark42-cluster/store/image/s3"
)

func main() {
	config := config.Init()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	imageStore, err := s3.Connect(&s3.StorageOption{
		AccessKey:   config.S3AccessKey(),
		Secret:      config.S3Secret(),
		EndpointURL: config.S3Endpoint(),
		Region:      config.S3Region(),
		BucketName:  config.S3Bucket(),

		Logger: logger,
	})

	if err != nil {
		panic(err)
	}

	engineStoreOpts := mongo.DefaultOptions(logger)
	engineStore, err := mongo.Connect(config.DbURL(), &engineStoreOpts)
	if err != nil {
		panic(err)
	}

	rpc.New(&rpc.Options{
		Port:    config.Port(),
		Storage: imageStore,
		Engine:  engineStore,
		Logger:  logger,
	})

	select {}
}
