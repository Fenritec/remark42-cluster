// © Fenritec S.A.S. France 2022 released under EUPL v1.2

package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type (
	//Values contains all the config values
	Values interface {
		Port() int
		Debug() bool

		S3AccessKey() string
		S3Secret() string
		S3Region() string
		S3Endpoint() string
		S3Bucket() string

		DbURL() string
	}

	values struct {
		VPort  int  `long:"port" env:"PORT" required:"true" description:"The port the server is listening to"`
		VDebug bool `long:"debug" env:"DEBUG" description:"Show more verbose output in stderr"`

		VS3AccessKey string `long:"s3-access-key" env:"S3_ACCESS_KEY" required:"true" description:"Amazon S3 access key"`
		VS3Secret    string `long:"s3-secret" env:"S3_SECRET" required:"true" description:"Amazon S3 secret"`
		VS3Region    string `long:"s3-region" env:"S3_REGION" required:"true" description:"Amazon S3 region"`
		VS3Endpoint  string `long:"s3-endpoint" env:"S3_ENDPOINT_URL" required:"true" description:"Amazon S3 endpoint URL (ie. OVH, Scaleway, etc.)"`
		VS3Bucket    string `long:"s3-bucket" env:"S3_BUCKET_NAME" required:"true" description:"Amazon S3 bucket name"`

		VDbURL string `long:"db-URL" env:"DB_URL" required:"true" description:"Mongodb URL"`
	}
)

func (v *values) Port() int {
	return v.VPort
}

func (v *values) Debug() bool {
	return v.VDebug
}

func (v *values) S3AccessKey() string {
	return v.VS3AccessKey
}

func (v *values) S3Secret() string {
	return v.VS3Secret
}

func (v *values) S3Region() string {
	return v.VS3Region
}

func (v *values) S3Endpoint() string {
	return v.VS3Endpoint
}

func (v *values) S3Bucket() string {
	return v.VS3Bucket
}

func (v *values) DbURL() string {
	return v.VDbURL
}

//Init returns the config values needed for the app
func Init() Values {
	ret := &values{}

	p := flags.NewParser(ret, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	return ret
}
