package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func CreateCloudfrontContext(sess *session.Session) *cloudfront.CloudFront {
	return cloudfront.New(sess)
}

func ListDistributions(svc *cloudfront.CloudFront) (*cloudfront.DistributionList, error) {
	params := &cloudfront.ListDistributionsInput{}
	resp, err := svc.ListDistributions(params)
	return resp.DistributionList, err
}
