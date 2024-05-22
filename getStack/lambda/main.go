package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type Event struct {
	StackName string `json:"stackName"`
}

func HandleRequest(ctx context.Context, event Event) *cloudformation.DescribeStacksOutput {
	svc := cloudformation.New(session.New())
	input := &cloudformation.DescribeStacksInput{
		StackName: &event.StackName,
	}

	result, err := svc.DescribeStacks(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudformation.ErrCodeResourceScanNotFoundException:
				fmt.Println(cloudformation.ErrCodeResourceScanNotFoundException, aerr.Error())
			case cloudformation.ErrCodeResourceScanInProgressException:
				fmt.Println(cloudformation.ErrCodeResourceScanInProgressException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	return result

}

func main() {
	lambda.Start(HandleRequest)
}
