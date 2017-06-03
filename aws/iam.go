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


func CreateIAMRole(name, path string, svc *iam.IAM) (string, error) {
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
	output, err := svc.CreateRole(params)
	return *output.Role.Arn, err
}

func CreateIAMPolicy(name, policyDocument, description, path string, svc *iam.IAM) (string, error) {
	params := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(policyDocument),
		PolicyName:     aws.String(name),
		Description:    aws.String(description),
		Path:           aws.String(path),
	}
	output, err := svc.CreatePolicy(params)
	return *output.Policy.Arn, err
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

	return roleArn, AttachPolicyToRole(policyArn, roleName, svc)
}