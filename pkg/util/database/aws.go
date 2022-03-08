package database

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var AwsSession *session.Session

func NewAwsSession() (error, bool) {
	AccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion := os.Getenv("AWS_REGION")

	// AccessKeyID := "AKIAUM2TD2VTAWVCAQXZ"
	// SecretAccessKey := "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"
	// MyRegion := "ap-southeast-1"

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		return err, false
	} else {
		AwsSession = sess
	}
	return nil, true
}

func GetAWSSession() *session.Session {
	return AwsSession
}
