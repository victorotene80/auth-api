package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
)

var _ contracts.TokenGenerator = (*JWTGenerator)(nil)

type JWTGenerator struct {
	secret []byte
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{secret: []byte(secret)}
}

func (j *JWTGenerator) GenerateAccess(
	userID, sessionID, role string,
	duration time.Duration,
) (contracts.Token, error) {

	exp := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"sub":  userID,
		"sid":  sessionID,
		"role": role,
		"exp":  exp.Unix(),
		"typ":  "access",
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := t.SignedString(j.secret)
	if err != nil {
		return contracts.Token{}, err
	}

	return contracts.Token{
		Value:     tokenStr,
		ExpiresAt: exp,
	}, nil
}

func (j *JWTGenerator) GenerateRefresh(
	userID, sessionID string,
	duration time.Duration,
) (contracts.Token, error) {

	exp := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"sub": userID,
		"sid": sessionID,
		"exp": exp.Unix(),
		"typ": "refresh",
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := t.SignedString(j.secret)
	if err != nil {
		return contracts.Token{}, err
	}

	return contracts.Token{
		Value:     tokenStr,
		ExpiresAt: exp,
	}, nil
}

func (j *JWTGenerator) ValidateAccess(
	tokenStr string,
) (string, string, string, error) {

	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !tok.Valid {
		return "", "", "", err
	}

	claims := tok.Claims.(jwt.MapClaims)

	if claims["typ"] != "access" {
		return "", "", "", errors.New("not an access token")
	}

	return claims["sub"].(string),
		claims["sid"].(string),
		claims["role"].(string),
		nil
}

func (j *JWTGenerator) ValidateRefresh(
	tokenStr string,
) (string, string, error) {

	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !tok.Valid {
		return "", "", err
	}

	claims := tok.Claims.(jwt.MapClaims)

	if claims["typ"] != "refresh" {
		return "", "", errors.New("not a refresh token")
	}

	return claims["sub"].(string),
		claims["sid"].(string),
		nil
}
