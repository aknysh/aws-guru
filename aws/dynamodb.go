package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func CreateDynamoDBContext(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func PutItem(tableName string, r interface{}, svc *dynamodb.DynamoDB) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(r)
	if err != nil {
		return nil, err
	}

	return svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
}
