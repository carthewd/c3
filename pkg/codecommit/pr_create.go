package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/carthewd/c3/internal/data"
)

func CreatePR(c *codecommit.CodeCommit, newPR data.NewPullRequest) (string, error) {
	prTargets := make([]*codecommit.Target, 1)
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
