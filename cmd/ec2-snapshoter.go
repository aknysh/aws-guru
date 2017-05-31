package cmd

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/aws"
	"errors"
	awslib "aws-guru/aws"
	"aws-guru/utils"
)


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


func main() {
	cronName := os.Getenv("CRON_NAME"); if cronName == "" {
		utils.ExitWithError(errors.New("CRON_NAME cannot be null!"))
	}

	cronPattern := os.Getenv("CRON_PATTERN"); if cronPattern == "" {
		utils.ExitWithError(errors.New("CRON_PATTERN cannot be null!"))
	}
	region := os.Getenv("REGION"); if region == "" {
		utils.ExitWithError(errors.New("REGION cannot be null!"))
	}

	accountId := os.Getenv("ACCOUNT_ID"); if accountId == "" {
		utils.ExitWithError(errors.New("ACCOUNT_ID cannot be null!"))
	}

	sess := awslib.CreateSession(&region)
	ec2Svc := awslib.CreateEC2Context(sess)
	cloudwatchEventsSvc := cloudwatchevents.New(sess)

	fmt.Println("Listing volumes...")

	volumes, err := awslib.ListVolumes(ec2Svc); if err != nil {
		utils.ExitWithError(err)
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
