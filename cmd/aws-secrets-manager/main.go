package main

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"fmt"
	"log"
	"os"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {

	log.Print("aws.secret.manager: Initializing the password retrieving from AWS")
	SecretName := os.Args[1]
	AWSRegion := os.Args[2]
	PASSWORD_HOME := os.Args[3]
	FILE_NAME := "password-tmp"
	FILE_PATH := PASSWORD_HOME + FILE_NAME
	sess, err := session.NewSession()
	svc := secretsmanager.New(sess, &aws.Config{
		Region: aws.String(AWSRegion),
	})
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(SecretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		var result map[string]interface{}
		json.Unmarshal([]byte(secretString), &result)

		keystorePass := result[SecretName]
		if keystorePass == nil {
			log.Print("aws.secret.manager: Fetched password is null.")
			return
		}

		writeOutput(keystorePass.(string), FILE_PATH)
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Print("aws.secret.manager: Base64 Decode Error:", err)
			return
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		writeOutput(decodedBinarySecret, FILE_PATH)
	}
}
func writeOutput(output string, path string) {

	log.Print("aws.secret.manager: Writing the retrieved password to a file.")
	f, err := os.Create(path)
	if err != nil {
		log.Print("aws.secret.manager: Error while writing the password to file.")
		return
	}
	defer f.Close()
	
	f.WriteString(output)
	readOutput(path)
}


func readOutput( path string) {

	log.Print("aws.secret.manager: Reading from the file to validating the write process.")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}


	if content == nil {
		log.Print("aws.secret.manager: File does not contain the expected content.")
		return
	}

	log.Print("aws.secret.manager: Changing the ownership of the file from root user to WSO2 user.")
	// Change permissions Linux.
	os.Chmod(path, 0777)

	// Change file ownership.
	os.Chown(path, 802, 802)
}