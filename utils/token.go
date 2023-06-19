package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type TokenDetails struct {
	Token *string
	TokenUuid string
	UserID string
	ExpiresAt *int64
}

func GenerateToken(userID string, expireTimeDuration time.Duration, privateKey string) (*TokenDetails, error) {
	td := &TokenDetails{
		ExpiresAt: new(int64),
		Token: new(string),
		TokenUuid: uuid.NewV4().String(),
		UserID: userID,
	}
	*td.ExpiresAt = time.Now().Add(expireTimeDuration).Unix()

	// obtain the private key from the environment variable
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return nil, err
	}

	// create a new token with the user id as the subject and the expiry time as the expiry
	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims {
		"token_uuid": td.TokenUuid,
		"sub": td.UserID,
		"exp": td.ExpiresAt,
		"iat": time.Now().Unix(),
	}).SignedString(key)

	if err != nil {
		return nil, err
	}

	return td, nil
}


func VerifyToken(tokenString string, publicKey string) (*TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error while decoding public key: %v", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("error while parsing public key: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &TokenDetails{
			TokenUuid: claims["token_uuid"].(string),
			UserID: claims["sub"].(string),
		}, nil

	} else {
		return nil, fmt.Errorf("invalid token")
	}

}
