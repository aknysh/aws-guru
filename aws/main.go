package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateSession(region *string) *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*region),
	}))

	return sess
}

func CreateProfiledSession(region *string, profile *string) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(*region)},
		Profile: *profile,
	}))

	return sess
}
