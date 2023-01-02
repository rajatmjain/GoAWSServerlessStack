package main

import (
	"GoAWSServerlessStack/pkg/handlers"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/joho/godotenv"
)

var(
	dynaClient dynamodbiface.DynamoDBAPI
)

func main(){
	godotenv.Load(".env")
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	},)
	if err!=nil{
		log.Fatalf("AWS session error: %v",err)
		return
	}

	dynaClient := dynamodb.New(awsSession)
	res, err := lambda.Start(handler())
}

const tableName = "LambdaInGoUser"

func handler(req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse,error){
	switch req.HTTPMethod{
	case "GET":
		return handlers.GetUser(req,tableName,dynaClient)
	
	case "POST":
		return handlers.CreateUser(req,tableName,dynaClient)
	
	case "PUT":
		return handlers.UpdateUser(req,tableName,dynaClient)

	case "DELETE":
		return handlers.DeleteUser(req,tableName,dynaClient)
	
	default:
		return handlers.UnhandledMethod()
	}
}