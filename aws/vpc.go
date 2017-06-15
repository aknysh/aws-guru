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

func AttachInternetGatewayToVPC(vpcId, ipgwId string, svc ec2.EC2) (error) {
	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(ipgwId),
		VpcId:             aws.String(vpcId),
	}

	_, err := svc.AttachInternetGateway(input)
	return err
}

func CreateRouteTable(vpcId string, svc ec2.EC2) (*ec2.CreateRouteTableOutput, error) {
	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
	}

	result, err := svc.CreateRouteTable(input)
	return result, err
}

func CreateRoute(cidr, gatewayId, routeTableId string, svc ec2.EC2) (error) {
	input := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String(cidr),
		GatewayId:            aws.String(gatewayId),
		RouteTableId:         aws.String(routeTableId),
	}

	_, err := svc.CreateRoute(input)
	return err
}

func AttachRouteTableToSubnet(routeTableId, subnetId string, svc ec2.EC2) (error) {
	input := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(routeTableId),
		SubnetId:     aws.String(subnetId),
	}

	_, err := svc.AssociateRouteTable(input)
	return err
}
