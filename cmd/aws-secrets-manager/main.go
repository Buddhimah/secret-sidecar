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
	//os.Setenv("SECRET_NAME", "catsndogs")
	//os.Setenv("AWS_REGION", "us-east-1")
	fmt.Println("==============START ADDING KEYSTORE PASSWORD============")
	SecretName := os.Getenv("SECRET_NAME")
	AWSRegion := os.Getenv("AWS_REGION")
	IS_HOME := os.Getenv("IS_HOME")

	fmt.Println("==============IS HOME============")
	fmt.Println(IS_HOME)
	writeOutputDummy("TEST STRING","/home/")
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
		fmt.Println(secretString)

		// The object stored in the "birds" key is also stored as
		// a map[string]interface{} type, and its type is asserted from
		// the interface{} type
		keystorePass := result[SecretName]
		writeOutput(keystorePass.(string), IS_HOME)
		fmt.Println(keystorePass.(string))
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		writeOutput(decodedBinarySecret, IS_HOME)
		fmt.Println(decodedBinarySecret)
	}
}
func writeOutput(output string, path string) {
	fmt.Println("=====WRITE TO FILE=====")
	fmt.Println(path + "password-persist")
	fmt.Println("=====OUTPUT VALUE=====")
	fmt.Println(output)
	f, err := os.Create(path + "password-persist")
	if err != nil {
		return
	}
	defer f.Close()
	
	f.WriteString(output)
	readOutput(path + "password-persist")
}

func writeOutputDummy(output string, path string) {
	f, err := os.Create("/tmp/secret")
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(output)
}

func readOutput( path string) {
	fmt.Println("=====READ FROM FILE=====")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	fmt.Println(text)
}