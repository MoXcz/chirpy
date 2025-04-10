package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	issuerToken = "chirpy"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", fmt.Errorf("Could not hash password. Err: %w", err)
	}
	return string(hash), nil

}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuerToken,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	customClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &customClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != issuerToken {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization value")
	}

	splitAuthHeader := strings.Fields(authHeader)
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" { // expects "Bearer TOKEN_STRING" as value
		return "", errors.New("Invalid header value")
	}

	tokenString := splitAuthHeader[1]
	return tokenString, nil
}

// Token used to refresh JWT
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization value")
	}

	splitAuthHeader := strings.Fields(authHeader)
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "ApiKey" { // expects "ApiKey KEY" as value
		return "", errors.New("Invalid header value")
	}

	apiKey := splitAuthHeader[1]
	return apiKey, nil
}
