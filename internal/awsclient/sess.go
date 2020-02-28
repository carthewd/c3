package awsclient

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession returns an AWS session
func NewSession(name string) *session.Session {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}

	if name == "" {
		return session.New(&aws.Config{Region: aws.String(region)})
	}
	return session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", name),
	})
}
