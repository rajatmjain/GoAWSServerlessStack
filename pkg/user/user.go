package user

import (
	"GoAWSServerlessStack/pkg/validators"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToMarshalRecord = "failed to marshal data"
	ErrorFailedToFetchRecords = "failed to fetch records"
	ErrorInvalidUserData = "invalid user data"
	ErrorInvalidEmail = "invalid email"
	ErrorFailedToDeleteRecord = "failed to delete record"
	ErrorFailedToDynamoPutRecord = "failed to dynamo put record"
	ErrorUserAlreadyExists = "user already exists"
	ErrorUserDoesntExists = "user doesn't exists"

)

type User struct{
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchUser(email string,tablename string,dynaClient dynamodbiface.DynamoDBAPI)(*User,error){
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	res, err := dynaClient.GetItem(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(res.Item,item)
	if err != nil{
		return nil,errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item,nil
}

func FetchUsers(tablename string,dynaClient dynamodbiface.DynamoDBAPI)(*[]User,error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}
	res, err := dynaClient.Scan(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToFetchRecords)
	}
	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items,items)
	if err!=nil{
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items,nil

}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error){
	var u User
	if err := json.Unmarshal([]byte(req.Body),&u); err!=nil{
		return nil,errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	res, _ := FetchUser(u.Email,tablename,dynaClient)
	if res != nil && len(res.Email) != 0{
		return nil,errors.New(ErrorUserAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err!=nil{
		return nil, errors.New(ErrorFailedToMarshalRecord)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tablename),

	}
	_,err = dynaClient.PutItem(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToDynamoPutRecord)
	}

	return &u,nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error){
	var u User
	if err := json.Unmarshal([]byte(req.Body),&u); err!=nil{
		return nil,errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	res, _ := FetchUser(u.Email,tablename,dynaClient)
	if res == nil{
		return nil,errors.New(ErrorUserDoesntExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err!=nil{
		return nil, errors.New(ErrorFailedToMarshalRecord)
	}
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tablename),
	}
	_,err = dynaClient.PutItem(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedToDynamoPutRecord)
	}
	return &u,nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(error){
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S:aws.String(email),
			},
		},
		TableName: &tablename,
	}
	_, err := dynaClient.DeleteItem(input)
	if err!=nil{
		return errors.New(ErrorFailedToDeleteRecord)
	}
	return nil
}