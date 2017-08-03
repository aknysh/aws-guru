package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	awslib "netguru/aws-guru/aws"
	"netguru/aws-guru/utils"
)

var iamAimLinkerCmd = &cobra.Command{
	Use:   "iam-scan",
	Short: "Links account to AWS-IAM-Manager",
	Long:  `Establises a trust relationship between two accounts and links account to AWS-IAM-Manager by creating Cross-Account role`,

	Run: func(cmd *cobra.Command, args []string) {
		iamLinkerRun()
	},
}

func init() {
	s3CloudfrontReportCmd.Flags().StringVarP(&path, "path", "p", "/", "Slave account name")
	s3CloudfrontReportCmd.Flags().StringVarP(&slaveName, "slaveName", "n", "", "Slave account name")
	s3CloudfrontReportCmd.Flags().StringVarP(&accountId, "id", "i", "696776776974", "Master account ID")
	s3CloudfrontReportCmd.Flags().StringVarP(&profile, "masterProfile", "m", "default", "Master account profile")
	s3CloudfrontReportCmd.Flags().StringVarP(&slaveProfile, "slaveProfile", "s", "", "Slave account profile")

	RootCmd.AddCommand(iamAimLinkerCmd)
}

func getTrustRelationshipPolicy(accountId string) string {
	policy := `{
  		"Version": "2012-10-17",
  		"Statement": [{
		 	"Effect": "Allow",
		 	"Principal": {
				"AWS": "arn:aws:iam::%s:root"
		 	},
     		"Action": "sts:AssumeRole"
    	}]
    }`

	return fmt.Sprintf(policy, accountId)
}

func getCrossAccountPolicyDocument(roleArn string) string {
	crossAccountPolicyDocument := `{
		"Version": "2012-10-17",
		"Statement": {
			"Effect": "Allow",
			"Action": "sts:AssumeRole",
			"Resource": "%s"
		}
	}`

	return fmt.Sprintf(crossAccountPolicyDocument, roleArn)
}


func iamLinkerRun() {
	if slaveProfile == "" {
		utils.ExitWithError(errors.New("--slaveProfile cannot be undefined!"))
	}

	if slaveName == "" {
		utils.ExitWithError(errors.New("--slaveName cannot be undefined!"))
	}

	administratorPolicyDocument := `{
  		"Version": "2012-10-17",
  		"Statement": [{
			"Effect": "Allow",
			"Action": "*",
			"Resource": "*"
    	}]
    }
	`

	slaveSess := awslib.CreateProfiledSession(&"us-east-1", &slaveProfile)
	masterSess := awslib.CreateProfiledSession(&"us-east-1", &profile)

	slaveIamSvc := awslib.CreateIAMContext(slaveSess)
	masterIamSvc := awslib.CreateIAMContext(masterSess)

	roleName := "AIM_FederationRole"

	_, err := awslib.CreateIAMRoleDetailed(roleName, path, "Netguru's role for Cross Account access",
		getTrustRelationshipPolicy(accountId), slaveIamSvc)
	if err != nil {
		utils.ExitWithError(err)
	}

	policy, err := awslib.CreateIAMPolicy("AIM_AdministratorAccess", administratorPolicyDocument,
		"AIM Administrator Policy", "/", slaveIamSvc)
	if err != nil {
		utils.ExitWithError(err)
	}

	err = awslib.AttachPolicyToRole(*policy.Policy.Arn, roleName, slaveIamSvc)
	if err != nil {
		utils.ExitWithError(err)
	}

	policy, err = awslib.CreateIAMPolicy(fmt.Sprintf("AIM_%s_AssumeRolePolicy"), getCrossAccountPolicyDocument(*role.Policy.Arn),
		"Policy allowing AWS resource to assume role of other AWS Account Role", "/", masterIamSvc)
	if err != nil {
		utils.ExitWithError(err)
	}


}
