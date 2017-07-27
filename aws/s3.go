package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateS3Context(sess *session.Session) *s3.S3 {
	return s3.New(sess)
}

func ListBuckets(svc *s3.S3) ([]*s3.Bucket, error) {
	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(params)
	return resp.Buckets, err
}
