package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"context"
	"fmt"
	"os"
)

var svc dynamodbiface.DynamoDBAPI

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return
	}
	svc = dynamodb.New(sess)
	lambda.Start(handler)
}

// EcsTaskStateChangeEvent struct to hold info about event.Detail.
// Refer the doc for 'event' itself: https://github.com/aws/aws-lambda-go/blob/master/events/cloudwatch_events.go
type EcsTaskStateChangeEvent struct {
	TaskArn       string `json:taskArn`
	ClusterArn    string `json:clusterArn`
	LaunchType    string `json:launchType`
	DesiredStatus string `json:desiredStatus`
	LastStatus    string `json:lastStatus`
	Group         string `json:group`
	StartedAt     string `json:startedAt,omitempty`
	TaskDefArn    string `json:taskDefinitionArn`
}

// Task struct to hold info about ECS task
type Task struct {
	Arn       string // primary key
	Cluster   string
	Name      string
	StartedAt string
	Family    string
	TaskDef   string
}

func hasTaskStarted(desiredStatus, lastStatus string) bool {
	return desiredStatus == "RUNNING" && lastStatus == "RUNNING"
}

func hasTaskStopped(desiredStatus, lastStatus string) bool {
	return desiredStatus == "STOPPED" && lastStatus == "STOPPED"
}

func getClusterName(clusterArn string) string {
	arnObj, err := arn.Parse(clusterArn)
	if err != nil {
		return clusterArn
	}
	return arnObj.Resource
}

func getTaskName(taskArn string) string {
	arnObj, err := arn.Parse(taskArn)
	if err != nil {
		return taskArn
	}
	return arnObj.Resource
}

func getTaskDefName(taskDefArn string) string {
	arnObj, err := arn.Parse(taskDefArn)
	if err != nil {
		return taskDefArn
	}
	return arnObj.Resource
}

func taskEventToTask(taskEvent EcsTaskStateChangeEvent) Task {
	return Task{
		Arn:       taskEvent.TaskArn,
		Cluster:   getClusterName(taskEvent.ClusterArn),
		Name:      getTaskName(taskEvent.TaskArn),
		StartedAt: taskEvent.StartedAt,
		Family:    taskEvent.Group,
		TaskDef:   getTaskDefName(taskEvent.TaskDefArn),
	}
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	var taskEvent EcsTaskStateChangeEvent
	if err := json.Unmarshal(event.Detail, &taskEvent); err != nil {
		return fmt.Errorf("unable to parse event detail body: %v", err)
	}

	task := taskEventToTask(taskEvent)
	if hasTaskStarted(taskEvent.DesiredStatus, taskEvent.LastStatus) {
		av, err := dynamodbattribute.MarshalMap(task)
		if err != nil {
			return err
		}
		input := &dynamodb.PutItemInput{
			TableName:           aws.String(os.Getenv("DDB_RUNNING_TASKS_TABLE")),
			ConditionExpression: aws.String("attribute_not_exists(Arn)"),
			Item:                av,
		}
		_, err = svc.PutItem(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
					fmt.Printf("task event ignored since it has already been stored: %+v\n", taskEvent)
					return nil
				}
			}
			return fmt.Errorf("got error calling PutItem: %v", err)
		}
	} else if hasTaskStopped(taskEvent.DesiredStatus, taskEvent.LastStatus) {
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv("DDB_RUNNING_TASKS_TABLE")),
			Key: map[string]*dynamodb.AttributeValue{
				"Arn": {
					S: aws.String(task.Arn),
				},
			},
		}
		_, err := svc.DeleteItem(input)
		if err != nil {
			return fmt.Errorf("got error calling DeleteItem: %v", err)
		}
	} else {
		fmt.Printf("task event is not required to store in DDB: %+v\n", taskEvent)
	}
	return nil
}
