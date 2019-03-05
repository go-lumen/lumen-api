package helpers

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetRSAPrivateKey(encodedKey []byte) (*rsa.PrivateKey, error) {
	//Decode base64 key
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(encodedKey)))
	base64.StdEncoding.Decode(base64Text, []byte(encodedKey))

	//Parse RSA key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(base64Text)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GenerateAccessToken(encodedKey []byte, subject string) (*string, error) {
	privateKey, err := GetRSAPrivateKey(encodedKey)
	if err != nil {
		return nil, err
	}

	access := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": subject,
		"aud": "access",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * time.Duration(8760)).Unix(),
	})

	accessString, err := access.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return &accessString, nil
}

func ValidateJwtToken(token string, encodedKey []byte, audience string) (jwt.MapClaims, error) {
	privateKey, err := GetRSAPrivateKey(encodedKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey

	rawToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return &publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	//validate algorithm
	if rawToken.Header["alg"] != jwt.SigningMethodRS256.Alg() {
		return nil, err
	}

	//validate signature
	if !rawToken.Valid {
		return nil, errors.New("token in invalid (wrong signature)")
	}

	claims, ok := rawToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse the claims")
	}

	//validate exp
	tokenExp := claims["exp"].(float64)
	if tokenExp < float64(time.Now().Unix()) {
		return nil, errors.New("token in invalid (expired)")
	}

	//validate aud
	tokenAud := claims["aud"].(string)
	if tokenAud != audience {
		return nil, errors.New("token in invalid (wrong audience)")
	}

	return claims, nil
}
