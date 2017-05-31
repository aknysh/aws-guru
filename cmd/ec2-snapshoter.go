package cmd

import (
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/ec2"
"fmt"
"os"
"github.com/aws/aws-sdk-go/service/cloudwatchevents"
"github.com/aws/aws-sdk-go/aws"
"errors"
"github.com/aws/aws-sdk-go/service/iam"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func listVolumes(svc *ec2.EC2) ([]*ec2.Volume, error) {
	params := &ec2.DescribeVolumesInput{}
	resp, err := svc.DescribeVolumes(params);
	return resp.Volumes, err
}

func createScheduledExpression(name, cronPattern string, svc *cloudwatchevents.CloudWatchEvents) error {
	params := &cloudwatchevents.PutRuleInput{
		Name:               aws.String(name),
		Description:        aws.String("Invokes Cloudwatch event to take EC2 Volume Snapshot"),
		ScheduleExpression: aws.String(cronPattern),
	}
	_, err := svc.PutRule(params)
	return err
}

func getTargetArnName(region, accountId, stackName string) string {
	return fmt.Sprintf("arn:aws:automation:%s:%s:action/EBSCreateSnapshot/EBSCreateSnapshot_%s",
		region, accountId, stackName)
}

func getInputArnName(region, accountId, ebsVolume string) string {
	return fmt.Sprintf("\"arn:aws:ec2:%s:%s:volume/%s\"",
		region, accountId, ebsVolume)
}

func prepareCloudWatchEventTargets(region, accountId, stackName string, volumes []*ec2.Volume) []*cloudwatchevents.Target {
	targets := make([]*cloudwatchevents.Target, len(volumes))

	for i, volume := range volumes {
		targets[i] = &cloudwatchevents.Target{
			Arn: aws.String(getTargetArnName(region, accountId, stackName)),
			Id:  aws.String("1231312312"),
			Input: aws.String(getInputArnName(region, accountId, *volume.VolumeId)),
		}
	}

	return targets
}

func createIAMRole(name string, svc *iam.IAM) error {
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
		Path:                     aws.String("/aws-guru"),
	}
	_, err := svc.CreateRole(params)
	return err
}

func createIAMPolicy(name string, svc *iam.IAM) error {

}

func attachIAMPolicy(name string, svc *iam.IAM) {

}

func putCloudwatchEventTargets(name string, targets []*cloudwatchevents.Target, svc *cloudwatchevents.CloudWatchEvents) error {
	params := &cloudwatchevents.PutTargetsInput{
		Rule: aws.String(name),
		Targets: targets,
	}
	_, err := svc.PutTargets(params)
	return err
}

func main() {
	cron_name := os.Getenv("CRON_NAME"); if cron_name == "" {
		exitWithError(errors.New("CRON_NAME cannot be null!"))
	}

	cron_pattern := os.Getenv("CRON_PATTERN"); if cron_name == "" {
		exitWithError(errors.New("CRON_PATTERN cannot be null!"))
	}
	region := os.Getenv("REGION"); if cron_name == "" {
		exitWithError(errors.New("REGION cannot be null!"))
	}

	account_id := os.Getenv("ACCOUNT_ID"); if cron_name == "" {
		exitWithError(errors.New("ACCOUNT_ID cannot be null!"))
	}


	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	ec2Svc := ec2.New(sess)
	cloudwatchEventsSvc := cloudwatchevents.New(sess)

	fmt.Println("Listing volumes...")

	volumes, err := listVolumes(ec2Svc); if err != nil {
		exitWithError(err)
	}

	fmt.Println("Creating scheduled expression...")

	err = createScheduledExpression(cron_name, cron_pattern, cloudwatchEventsSvc); if err != nil {
		exitWithError(err)
	}

	fmt.Println("Attaching target volumes to scheduled expression...")

	err = putCloudwatchEventTargets(
		cron_name,
		prepareCloudWatchEventTargets(region, account_id, "stack", volumes),
		cloudwatchEventsSvc,
	); if err != nil {
		exitWithError(err)
	}

	fmt.Println("Done!")
}
