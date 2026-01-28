package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	stsClient *sts.Client
}

func main() {
	h, err := NewHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create handler")
	}

	lambda.Start(h.handle)
}

func NewHandler() (*Handler, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg)

	return &Handler{
		stsClient: stsClient,
	}, nil
}

type UpstreamResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (h Handler) handle(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	token, err := h.stsClient.GetWebIdentityToken(ctx, &sts.GetWebIdentityTokenInput{
		Audience:         []string{"no-more-long-lived-credentials-lambda"},
		SigningAlgorithm: aws.String("RS256"),
		DurationSeconds:  aws.Int32(300),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Info().Str("token", aws.ToString(token.WebIdentityToken)).Msg("obtained web identity token")

	req, err := http.NewRequest("GET", "https://nomorelonglivedcredeoiul92go-no-more-long-lived-credentials.functions.fnc.fr-par.scw.cloud/", nil)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+aws.ToString(token.WebIdentityToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	defer resp.Body.Close()

	log.Info().Int("status_code", resp.StatusCode).Msg("invoked function with web identity token")

	bodyBytes, err := io.ReadAll(resp.Body)
	upstreamResponse := UpstreamResponse{
		Message:    string(bodyBytes),
		StatusCode: resp.StatusCode,
	}

	body, err := json.Marshal(upstreamResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil

}
