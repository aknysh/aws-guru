package cmd

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/aws"
	awslib "aws-guru/aws"
	"aws-guru/utils"
	"strings"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
)

var ec2snapshoterCmd = &cobra.Command{
	Use:   "ec2-snapshoter",
	Short: "Setup automatic EC2 snapshots using Cloudwatch Events",
	Long: `EC2 Snapshoter configures a scheduled expression (Cloudwatch Event) which will take snapshot of your EC2 volumes every X hours.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

var cronPattern string
var cronName string
var region string
var accountId string

func init() {
	ec2snapshoterCmd.Flags().StringVarP(&cronPattern, "cron-pattern", "c", "0 10 * * ? *", "scheduled expression cron pattern (UTC time)")
	ec2snapshoterCmd.Flags().StringVarP(&cronName, "cron-name", "n", "ec2-snapshoter", "name of the scheduled expression")
	ec2snapshoterCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "region where backing up should be setup")
	ec2snapshoterCmd.Flags().StringVarP(&accountId, "account-id", "i", "", "AWS Account ID")
	RootCmd.AddCommand(ec2snapshoterCmd)
}

func getAutomationTargetArnName(region, accountId, stackName string) string {
	return fmt.Sprintf("arn:aws:automation:%s:%s:action/EBSCreateSnapshot/EBSCreateSnapshot_%s",
		region, accountId, stackName)
}

func getEC2VolumeInputArnName(region, accountId, ebsVolume string) string {
	return fmt.Sprintf("\"arn:aws:ec2:%s:%s:volume/%s\"",
		region, accountId, ebsVolume)
}

func prepareCloudWatchEventTargets(region, accountId, stackName string, volumes []*ec2.Volume) []*cloudwatchevents.Target {
	targets := make([]*cloudwatchevents.Target, len(volumes))

	for i, volume := range volumes {
		targets[i] = &cloudwatchevents.Target{
			Arn: aws.String(getAutomationTargetArnName(region, accountId, stackName)),
			Id:  aws.String(fmt.Sprintf("id_%d", i)),
			Input: aws.String(getEC2VolumeInputArnName(region, accountId, *volume.VolumeId)),
		}
	}

	return targets
}


func run(cmd *cobra.Command, args []string) {
	fmt.Printf(strings.Join(args, " "))

	sess := awslib.CreateSession(&region)
	ec2Svc := awslib.CreateEC2Context(sess)
	cloudwatchEventsSvc := cloudwatchevents.New(sess)

	fmt.Println("Listing volumes...")

	volumes, err := awslib.ListVolumes(ec2Svc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Printf("Found %d volumes to backup.\n", len(volumes))

	if len(volumes) == 0 {
		utils.ExitWithError(errors.New("Couldn't find volumes to backup"))
	}

	fmt.Println("Creating scheduled expression...")

	err = awslib.CreateScheduledExpression(cronName, "", cronPattern, cloudwatchEventsSvc); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Attaching target volumes to scheduled expression...")

	err = awslib.PutCloudwatchEventTargets(
		cronName,
		prepareCloudWatchEventTargets(region, accountId, "stack", volumes),
		cloudwatchEventsSvc,
	); if err != nil {
		utils.ExitWithError(err)
	}

	fmt.Println("Done!")
}
