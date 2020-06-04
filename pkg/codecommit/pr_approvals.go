package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	//"github.com/carthewd/c3/internal/data"
)

// Approve PR using the CodeCommit API
func ApprovePR(c *codecommit.CodeCommit, pr string) error {
	getPRInput := &codecommit.GetPullRequestInput{
		PullRequestId: aws.String(pr),
	}

	resp, err := c.GetPullRequest(getPRInput)
	revisionId := *resp.PullRequest.RevisionId

	approvalInput := &codecommit.UpdatePullRequestApprovalStateInput{
		ApprovalState: aws.String("APPROVE"),
		PullRequestId: aws.String(pr),
		RevisionId: aws.String(revisionId),
	}

	_, err = c.UpdatePullRequestApprovalState(approvalInput)

	if err != nil {
		return err
	}
	return err
}

func RevokePR(c *codecommit.CodeCommit, pr string) error {
	getPRInput := &codecommit.GetPullRequestInput{
		PullRequestId: aws.String(pr),
	}

	resp, err := c.GetPullRequest(getPRInput)
	revisionId := *resp.PullRequest.RevisionId

	approvalInput := &codecommit.UpdatePullRequestApprovalStateInput{
		ApprovalState: aws.String("REVOKE"),
		PullRequestId: aws.String(pr),
		RevisionId: aws.String(revisionId),
	}

	_, err = c.UpdatePullRequestApprovalState(approvalInput)

	if err != nil {
		return err
	}
	return err
}