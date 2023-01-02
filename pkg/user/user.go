package user

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecords = "failed to fetch records"
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
	err = dynamodbattribute.UnmarshalMap(res.Items[],items)
	return items,nil

}

func CreateUser()(){

}

func UpdateUser()(){

}

func DeleteUser()(){
	
}