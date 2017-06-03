package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
)

func CreateScheduledExpression(name, description, cronPattern string, svc *cloudwatchevents.CloudWatchEvents) error {
	params := &cloudwatchevents.PutRuleInput{
		Name:               aws.String(name),
		Description:        aws.String(description),
		ScheduleExpression: aws.String(fmt.Sprintf("cron(%s)", cronPattern)),
	}
	_, err := svc.PutRule(params)
	return err
}

func PutCloudwatchEventTargets(name string, targets []*cloudwatchevents.Target, svc *cloudwatchevents.CloudWatchEvents) error {
	params := &cloudwatchevents.PutTargetsInput{
		Rule: aws.String(name),
		Targets: targets,
	}
	_, err := svc.PutTargets(params)
	return err
}
