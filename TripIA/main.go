package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Trip struct {
	Destination string `json:"destination"`
	Month       string `json:"month"`
	Days        string `json:"days"`
}

type Request struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"maxTokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type Response struct {
	Completions []Completion `json:"completions"`
}

type Completion struct {
	Data Data `json:"data"`
}

type Data struct {
	Text string `json:"text"`
}

func TripIA(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var trip Trip
	err := json.Unmarshal([]byte(request.Body), &trip)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	prompt := fmt.Sprintf("Vou viajar para %s em %s, ficarei %s dias. Me dê dicas de viagem no seguinte formato: Estação do ano no período, Temperatura média nessa época do ano, Que tipo de roupa levar, Pontos Turísticos para Visitar. Responda em até 140 caracteres", trip.Destination, trip.Month, trip.Days)

	generation, err := RunBedrock(prompt)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(*generation),
	}, nil
}

func RunBedrock(prompt string) (*string, error) {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	brc := bedrockruntime.NewFromConfig(cfg)

	payload := Request{
		Prompt:      prompt,
		MaxTokens:   140,
		Temperature: 0.5,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	output, err := brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String("ai21.j2-mid-v1"),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("*/*"),
	})
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(output.Body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Completions[0].Data.Text, nil
}

func main() {
	lambda.Start(TripIA)
}
