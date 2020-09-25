package helpers

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-lumen/lumen-api/utils"
	"time"
)

// GetRSAPrivateKey retrieves the RSA private key
func GetRSAPrivateKey(encodedKey []byte) (*rsa.PrivateKey, error) {
	//Decode base64 key
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(encodedKey)))
	_, err := base64.StdEncoding.Decode(base64Text, []byte(encodedKey))
	utils.CheckErr(err)

	//Parse RSA key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(base64Text)
	if err != nil {
		utils.CheckErr(err)
		return nil, err
	}

	return privateKey, nil
}

// GenerateToken Generates an access token
func GenerateToken(encodedKey []byte, userID string, audience string, expiration int64) (*string, error) {
	privateKey, err := GetRSAPrivateKey(encodedKey)
	if err != nil {
		utils.CheckErr(err)
		return nil, err
	}

	access := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": userID,
		"aud": audience,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * time.Duration(expiration)).Unix(),
	})

	accessString, err := access.SignedString(privateKey)
	if err != nil {
		utils.CheckErr(err)
		return nil, err
	}

	return &accessString, nil
}

// ValidateJwtToken Validates a JWT token
func ValidateJwtToken(token string, encodedKey []byte, audience string) (jwt.MapClaims, error) {
	privateKey, err := GetRSAPrivateKey(encodedKey)
	if err != nil {
		utils.CheckErr(err)
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
