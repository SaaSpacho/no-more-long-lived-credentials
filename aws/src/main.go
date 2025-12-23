package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Handler struct {
	stsClient *sts.Client
}

func main() {
	h, err := NewHandler()
	if err != nil {
		log.Fatal(err)
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

	log.Println(*token.WebIdentityToken)
	return nil
}
