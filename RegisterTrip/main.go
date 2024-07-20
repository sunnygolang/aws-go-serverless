package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Trip struct {
	ID          string `json:"id"`
	Destination string `json:"destination"`
	OwnerEmail  string `json:"owner_email"`
	OwnerName   string `json:"owner_name"`
	StartsAt    string `json:"starts_at"`
	EndsAt      string `json:"ends_at"`
}

func RegisterTrip(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var trip Trip
	err := json.Unmarshal([]byte(request.Body), &trip)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	trip.ID = uuid.New().String()
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Trips"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(trip.ID),
			},
			"destination": {
				S: aws.String(trip.Destination),
			},
			"owner_email": {
				S: aws.String(trip.OwnerEmail),
			},
			"owner_name": {
				S: aws.String(trip.OwnerName),
			},
			"starts_at": {
				S: aws.String(trip.StartsAt),
			},
			"ends_at": {
				S: aws.String(trip.EndsAt),
			},
		},
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	body, err := json.Marshal(trip)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(RegisterTrip)
}
