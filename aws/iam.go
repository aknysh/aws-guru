package aws

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateIAMContext (sess *session.Session) *iam.IAM {
	return iam.New(sess)
}

func CreateIAMRole(name, path string, svc *iam.IAM) (*iam.CreateRoleOutput, error) {
	assumeRolePolicyDocument := `{
		"Version" : "2012-10-17",
		"Statement": [ {
			"Effect": "Allow",
			"Principal": {
			"Service": [ "events.amazonaws.com" ]
		},
		"Action": [ "sts:AssumeRole" ]
		} ]
	}`

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
		RoleName:                 aws.String(name),
		Path:                     aws.String(path),
	}
	return svc.CreateRole(params)
}

func CreateIAMPolicy(name, policyDocument, description, path string, svc *iam.IAM) (*iam.CreatePolicyOutput, error) {
	params := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(policyDocument),
		PolicyName:     aws.String(name),
		Description:    aws.String(description),
		Path:           aws.String(path),
	}
	return svc.CreatePolicy(params)
}

func AttachPolicyToRole(policyArn, roleName string, svc *iam.IAM) error {
	params := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	}
	_, err := svc.AttachRolePolicy(params)
	return err
}

func CreateRoleWithAttachedPolicy(name, path, policyDocument, description string, svc *iam.IAM) (string, error) {
	policyName := fmt.Sprintf("%s_policy", name)
	roleName := fmt.Sprintf("%s_role", name)

	roleArn, err := CreateIAMRole(roleName, path, svc)
	if err != nil {
		return "", err
	}

	policyArn, err := CreateIAMPolicy(policyName, policyDocument, description, path, svc)
	if err != nil {
		return "", err
	}

	return *roleArn.Role.Arn, AttachPolicyToRole(*policyArn.Policy.Arn, roleName, svc)
}