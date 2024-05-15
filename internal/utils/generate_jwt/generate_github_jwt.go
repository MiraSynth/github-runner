package generate_jwt

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func Generate(clientId string, secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": clientId,
		"exp": time.Now().Add(600 * time.Second).Unix(),
		"iat": time.Now().Add(-10 * time.Second).Unix(),
	})

	secretKeyBytes := []byte(secretKey)
	parsedSecretKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretKeyBytes)
	if err != nil {
		return "", err
	}

	tokenString, err := claims.SignedString(parsedSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
