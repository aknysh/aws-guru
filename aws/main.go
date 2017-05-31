package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

func CreateSession(region *string) *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*region),
	}))

	return sess
}