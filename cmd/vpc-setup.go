package cmd

import "github.com/spf13/cobra"
import (
	awslib "aws-guru/aws"
	"aws-guru/utils"
	"fmt"
)

var vpcSetupCmd = &cobra.Command{
	Use:   "vpc-setup",
	Short: "Provisions VPC and Subnets suitable for future development",
	Long: `Provisions VPC and Subnets suitable for future development`,

	Run: func(cmd *cobra.Command, args []string) {
		vpcSetupRun()
	},
}

var vpcCidr string
var privateSubnetCidr string
var publicSubnetCidr string

func init() {
	vpcSetupCmd.Flags().StringVarP(&vpcCidr, "vpc-cidr", "v", "10.0.0.0/16", "VPC CIDR")
	vpcSetupCmd.Flags().StringVarP(&privateSubnetCidr, "private-cidr", "r", "10.0.0.0/24", "Private Subnet CIDR")
	vpcSetupCmd.Flags().StringVarP(&publicSubnetCidr, "public-cidr", "p", "10.0.1.0/24", "Public Subnet CIDR")

	RootCmd.AddCommand(vpcSetupCmd)
}

func vpcSetupRun() {
	sess := awslib.CreateSession(&region)
	vpcSvc := awslib.CreateVPCContext(sess)

	fmt.Println("Creating VPC...")

	vpcResult, err := awslib.CreateVPC(vpcCidr, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Creating private subnet...")

	_, err = awslib.CreateSubnet(privateSubnetCidr, *vpcResult.Vpc.VpcId, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Creating public subnet...")

	publicSubnetResult, err := awslib.CreateSubnet(publicSubnetCidr, *vpcResult.Vpc.VpcId, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Creating Internet Gateway...")

	ipgwResult, err := awslib.CreateInternetGateway(*vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Attaching Internet Gateway...")

	err = awslib.AttachInternetGatewayToVPC(*vpcResult.Vpc.VpcId, *ipgwResult.InternetGateway.InternetGatewayId,
		*vpcSvc)
	if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Creating Route Table...")

	routeTableResult, err := awslib.CreateRouteTable(*vpcResult.Vpc.VpcId, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Creating Route...")

	err = awslib.CreateRoute("0.0.0.0/0", *ipgwResult.InternetGateway.InternetGatewayId,
		*routeTableResult.RouteTable.RouteTableId, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Attaching Route Table to Subnet...")

	err = awslib.AttachRouteTableToSubnet(*routeTableResult.RouteTable.RouteTableId,
		*publicSubnetResult.Subnet.SubnetId, *vpcSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Done!")
}
