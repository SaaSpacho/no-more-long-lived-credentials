package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/rs/zerolog/log"
)

var allowedIssuers = map[string]bool{"https://a1759d82-9b03-4cd4-8483-8660e54a25b2.tokens.sts.global.api.aws": true}

var keyFunc keyfunc.Keyfunc

func init() {
	var jwksEndpoints []string
	for issuer := range allowedIssuers {
		jwksEndpoints = append(jwksEndpoints, path.Join(issuer, ".well-known/jwks.json"))
	}

	f, err := keyfunc.NewDefaultCtx(context.Background(), jwksEndpoints)
	if err != nil {
		panic(fmt.Sprintf("failed to create JWKS keyfunc: %v", err))
	}
	keyFunc = f
}

func Handle(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}

	_, err = io.WriteString(w, "Access Granted")

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}

func checkAuth(r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("missing authorization token")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	issuer, err := validateIssuer(token)
	if err != nil {
		return fmt.Errorf("failed to validate issuer: %w", err)
	}
	log.Info().Str("issuer", issuer).Msg("validated issuer")

	parse, err := jwt.Parse(token, keyFunc.Keyfunc, jwt.WithIssuer(issuer))
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	log.Info().Msg("parsed token successfully")

	if !parse.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func validateIssuer(token string) (string, error) {
	parsedToken, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to assert claims as MapClaims")
	}

	iss, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("issuer (iss) claim is missing or not a string")
	}

	_, ok = allowedIssuers[iss]
	if !ok {
		return "", fmt.Errorf("issuer is not allowed: %s", iss)
	}

	return iss, nil
}
