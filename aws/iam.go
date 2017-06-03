package aws

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws"
)

func CreateIAMRole(name, path string, svc *iam.IAM) error {
	assumeRolePolicyDocument := `{
		"Version" : "2012-10-17",
			"Statement": [ {
		"Effect": "Allow",
		"Principal": {
		"Service": [ "ec2.amazonaws.com" ]
		},
		"Action": [ "sts:AssumeRole" ]
		} ]
	}`

	params := &iam.CreateRoleInput{
	AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
	RoleName:                 aws.String(name),
	Path:                     aws.String(path),
	}
	_, err := svc.CreateRole(params)
	return err
}

func CreateIAMPolicy(name, policyDocument string, svc *iam.IAM) (string, error) {
	return "", nil
}

func AttachPolicyToRole(roleArn, policyName string, svc *iam.IAM) error {
	return nil
}
