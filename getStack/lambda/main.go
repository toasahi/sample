package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type Response struct {
	Email string `json:"email"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := cloudformation.New(session.New())
	// name := event.PathParameters["stackName"]
	test2 := event.MultiValueQueryStringParameters["test"][1]
	// fmt.Print(request.QueryStringParameters["stackName"])
	// name := request.QueryStringParameters["stackName"]
	input := &cloudformation.DescribeStacksInput{
		StackName: &test2,
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
			fmt.Println(err.Error())
		}
	}

	parameter := result.Stacks[0].Parameters

	email := ""
	for _, param := range parameter {
		if aws.StringValue(param.ParameterKey) == "Email" {
			email = *param.ParameterValue
		}
	}

	response := Response{
		Email: email,
	}

	jsonBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
