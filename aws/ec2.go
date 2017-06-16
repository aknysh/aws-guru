package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateEC2Context (sess *session.Session) *ec2.EC2 {
	return ec2.New(sess)
}

func ListVolumes(svc *ec2.EC2) ([]*ec2.Volume, error) {
	params := &ec2.DescribeVolumesInput{}
	resp, err := svc.DescribeVolumes(params)
	return resp.Volumes, err
}

func DescribeReservations(svc *ec2.EC2) ([]*ec2.Reservation, error) {
	params := &ec2.DescribeInstancesInput{}
	resp, err := svc.DescribeInstances(params)
	return resp.Reservations, err
}
