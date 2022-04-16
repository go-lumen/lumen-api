package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// SigfoxAPIKey struct holds informations about Sigfox API Key
type SigfoxAPIKey struct {
	Name          string   `json:"name"`
	DeviceTypeIDs []string `json:"deviceTypeIDs"`
	SigfoxIDs     []string `json:"sigfoxIDs"`
	SigfoxKey     string   `json:"sigfoxKey"`
	SigfoxSecret  string   `json:"sigfoxSecret"`
}

// SigfoxAPIKeys struct is an array of SigfoxAPIKey
type SigfoxAPIKeys struct {
	Keys []SigfoxAPIKey `json:"data"`
}

// RetrieveSigfoxAPIKey allows to retrieve a SigfoxAPIKey from a Sigfox Device ID
func RetrieveSigfoxAPIKey(apiKeys SigfoxAPIKeys, sigfoxID, fleetSigfoxDeviceTypeID string) SigfoxAPIKey {
	var ret SigfoxAPIKey
	//apiKeys := ExtractSigfoxAPIKeyFromFile(GetString(c, "API_KEYS_FILENAME"))

	for _, key := range apiKeys.Keys {
		for _, sigID := range key.SigfoxIDs {
			if sigID == sigfoxID {
				return key
			}
		}
		for _, sfxDevTypID := range key.DeviceTypeIDs {
			if sfxDevTypID == fleetSigfoxDeviceTypeID {
				return key
			}
		}
	}

	return ret
}

// GetSigfoxAPIKeysFromAWS allows parsing secret file from AWS to retrieve SigfoxAPIKeys
func GetSigfoxAPIKeysFromAWS(awsAPIId, awsAPIKey, fileName string) SigfoxAPIKeys {
	var data SigfoxAPIKeys

	svc := secretsmanager.New(session.New(), &aws.Config{Credentials: credentials.NewStaticCredentials(awsAPIId, awsAPIKey, ""), Region: aws.String("eu-west-3")})
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(fileName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		logrus.Errorln("GetSecretValue Error:", err)
	}
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		err = json.Unmarshal([]byte(secretString), &data)
		if err != nil {
			logrus.Errorln(err)
		}
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			logrus.Errorln("Base64 Decode Error:", err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		err = json.Unmarshal([]byte(decodedBinarySecret), &data)
		if err != nil {
			logrus.Errorln(err)
		}
	}

	return data
}

// ExtractSigfoxAPIKeyFromFile allows parsing JSON file to retrieve SigfoxAPIKeys
func ExtractSigfoxAPIKeyFromFile(fileName string) SigfoxAPIKeys {
	var data SigfoxAPIKeys

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	return data
}
