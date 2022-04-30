# Remark42-Cluster

Remark42-cluster is a docker / kubernetes side-car container app to allow full use of s3 / mongodb and stop relying on local storage.

It implements the rpc servers interfaces of remark42 for the image and engine stores.

## Warning

This is a work in progress. Some features are missing.

## Command line arguments and env arguments

```
Usage:
  remark42-cluster [OPTIONS]

Application Options:
      --port=          The port the server is listening to [$PORT]
      --debug          Show more verbose output in stderr [$DEBUG]
      --s3-access-key= Amazon S3 access key [$S3_ACCESS_KEY]
      --s3-secret=     Amazon S3 secret [$S3_SECRET]
      --s3-region=     Amazon S3 region [$S3_REGION]
      --s3-endpoint=   Amazon S3 endpoint URL (ie. OVH, Scaleway, etc.) [$S3_ENDPOINT_URL]
      --s3-bucket=     Amazon S3 bucket name [$S3_BUCKET_NAME]
      --db-URL=        Mongodb URL [$DB_URL]

Help Options:
  -h, --help           Show this help message
```

## Licence

The work is released under EUPL v1.2

## Versions

The current version supports remark42 v1.9.0
