package awsclient

import (
	"github.com/aws/aws-sdk-go/service/codecommit"
)

// NewClient initializes a new AWS Code Commit client
func NewClient() *codecommit.CodeCommit {
	return codecommit.New(NewSession(""))
}
