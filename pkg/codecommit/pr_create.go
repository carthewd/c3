package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/carthewd/c3/internal/data"
)

// CreatePR uses the CodeCommit API to create a new pull request from a NewPullRequest struct
func CreatePR(c codecommitiface.CodeCommitAPI, newPR data.NewPullRequest) (string, error) {
	var prTargets []*codecommit.Target
	prTarget := &codecommit.Target{
		RepositoryName:  aws.String(newPR.Repository),
		SourceReference: aws.String(newPR.SourceRef),
	}

	prTargets = append(prTargets, prTarget)

	newPRInput := &codecommit.CreatePullRequestInput{
		Description: aws.String(newPR.Description),
		Title:       aws.String(newPR.Title),
		Targets:     prTargets,
	}

	result, err := c.CreatePullRequest(newPRInput)
	if err != nil {
		return "", err
	}

	return *result.PullRequest.PullRequestId, err
}
