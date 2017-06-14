package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
)

func CreateVPCContext (sess *session.Session) *ec2.EC2 {
	return ec2.New(sess)
}

func CreateVPC(cidr string, svc ec2.EC2) (*ec2.CreateVpcOutput, error){
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidr),
	}

	result, err := svc.CreateVpc(input)
	return result, err
}

func CreateSubnet(cidr, vpcId string, svc ec2.EC2) (*ec2.CreateSubnetOutput, error) {
	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidr),
		VpcId: aws.String(vpcId),
	}

	result, err := svc.CreateSubnet(input)
	return result, err
}

func CreateInternetGateway(svc ec2.EC2) (*ec2.CreateInternetGatewayOutput, error) {
	input := &ec2.CreateInternetGatewayInput{}

	result, err := svc.CreateInternetGateway(input)
	return result, err
}