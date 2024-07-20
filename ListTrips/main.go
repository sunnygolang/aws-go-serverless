package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Trip struct {
	ID          string `json:"id"`
	Destination string `json:"destination"`
	OwnerEmail  string `json:"owner_email"`
	OwnerName   string `json:"owner_name"`
	StartsAt    string `json:"starts_at"`
	EndsAt      string `json:"ends_at"`
}

func ListTrips(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("Trips"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	var trips []Trip
	for _, item := range result.Items {
		trips = append(trips, Trip{
			ID:          *item["id"].S,
			Destination: *item["destination"].S,
			OwnerEmail:  *item["owner_email"].S,
			OwnerName:   *item["owner_name"].S,
			StartsAt:    *item["starts_at"].S,
			EndsAt:      *item["ends_at"].S,
		})
	}

	body, err := json.Marshal(trips)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(ListTrips)
}
