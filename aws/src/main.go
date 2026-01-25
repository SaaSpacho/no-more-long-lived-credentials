package main

import (
	"context"
	"net/http"

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

func (h Handler) handle(ctx context.Context) error {
	token, err := h.stsClient.GetWebIdentityToken(ctx, &sts.GetWebIdentityTokenInput{
		Audience:         []string{"no-more-long-lived-credentials-lambda"},
		SigningAlgorithm: aws.String("RS256"),
		DurationSeconds:  aws.Int32(300),
	})
	if err != nil {
		return err
	}

	log.Info().Str("token", aws.ToString(token.WebIdentityToken)).Msg("obtained web identity token")

	req, err := http.NewRequest("GET", "https://nomorelonglivedcredeoiul92go-no-more-long-lived-credentials.functions.fnc.fr-par.scw.cloud/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+aws.ToString(token.WebIdentityToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Info().Int("status_code", resp.StatusCode).Msg("invoked function with web identity token")

	return nil
}
