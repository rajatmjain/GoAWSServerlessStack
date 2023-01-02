package handlers

import (
	"GoAWSServerlessStack/pkg/user"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "method not allowed"
type ErrorBody struct{
	ErrorMessage *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest,tablename string,dynaClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse,error){
	email := req.QueryStringParameters["email"]
	if len((email))>0{
		res, err := user.FetchUser(email,tablename,dynaClient)
		if err!=nil{
			return apiResponse(http.StatusBadRequest,ErrorBody{
				aws.String(err.Error()),
			})
		}
		return apiResponse(http.StatusOK,res)
	}
	res, err := user.FetchUsers(tablename,dynaClient)
	if err!= nil{
		return apiResponse(http.StatusBadRequest,ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK,res)
}

func CreateUser(req events.APIGatewayProxyRequest,tablename string,dynaClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse,error){

}

func UpdateUser()(*events.APIGatewayProxyResponse,error){

}

func DeleteUser(req events.APIGatewayProxyRequest,tablename string,dynaClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse,error){

}

func UnhandledMethod()(*events.APIGatewayProxyResponse,error){
	return apiResponse(http.StatusMethodNotAllowed,ErrorMethodNotAllowed)

}

