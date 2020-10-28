package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
)

// ClosePR uses the CodeCommit API to close a pull request from a PR ID
func ClosePR(c codecommitiface.CodeCommitAPI, pr string) (string, error) {
	closePRInput := &codecommit.UpdatePullRequestStatusInput{
		PullRequestId:     aws.String(pr),
		PullRequestStatus: aws.String("CLOSED"),
	}

	result, err := c.UpdatePullRequestStatus(closePRInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case codecommit.ErrCodeInvalidPullRequestStatusUpdateException:
				return "", err
			case codecommit.ErrCodePullRequestDoesNotExistException:
				return "", err
			}
		}
		return "", err
	}

	return *result.PullRequest.PullRequestStatus, err
}
