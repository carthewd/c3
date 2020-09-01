package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/carthewd/c3/internal/data"
)

func MergeOptions(c *codecommit.CodeCommit, pr data.PullRequest, repoName string) (data.MergeOptions, error) {
	mergeOptInput := &codecommit.GetMergeOptionsInput{
		RepositoryName:             aws.String(repoName),
		DestinationCommitSpecifier: aws.String(pr.DestBranch),
		SourceCommitSpecifier:      aws.String(pr.SourceBranch),
	}

	mergeOpts, err := c.GetMergeOptions(mergeOptInput)

	availMergeOpt := data.MergeOptions{}
	for _, option := range mergeOpts.MergeOptions {
		if *option == "FAST_FORWARD_MERGE" {
			availMergeOpt.FF = true
		} else if *option == "SQUASH_MERGE" {
			availMergeOpt.Squash = true
		} else if *option == "THREE_WAY_MERGE" {
			availMergeOpt.ThreeWay = true
		}
	}

	return availMergeOpt, err
}

func Merge(c *codecommit.CodeCommit, m data.MergeInput) error {
	switch m.Type {
	case "FF":
		fastForwardInput := &codecommit.MergePullRequestByFastForwardInput{
			PullRequestId:  aws.String(m.PRID),
			RepositoryName: aws.String(m.Repository),
			SourceCommitId: aws.String(m.SourceCommit),
		}
		mergeFF(c, fastForwardInput)
	case "Squash":
		squashMergeInput := &codecommit.MergePullRequestBySquashInput{
			PullRequestId:  aws.String(m.PRID),
			RepositoryName: aws.String(m.Repository),
			SourceCommitId: aws.String(m.SourceCommit),
			AuthorName:     aws.String(m.AuthorName),
			Email:          aws.String(m.AuthorEmail),
		}
		mergeSquash(c, squashMergeInput)
	case "ThreeWay":
		threeWayMergeInput := &codecommit.MergePullRequestBySquashInput{
			PullRequestId:  aws.String(m.PRID),
			RepositoryName: aws.String(m.Repository),
			SourceCommitId: aws.String(m.SourceCommit),
			AuthorName:     aws.String(m.AuthorName),
			Email:          aws.String(m.AuthorEmail),
		}
		mergeThreeWay(c, threeWayMergeInput)
	}

	return nil
}

func mergeFF(c *codecommit.CodeCommit, m *codecommit.MergePullRequestByFastForwardInput) (*codecommit.MergePullRequestByFastForwardOutput, error) {
	mergeResult, err := c.MergePullRequestByFastForward(m)
	if err != nil {
		return mergeResult, err
	}

	return mergeResult, err
}

func mergeSquash(c *codecommit.CodeCommit, m *codecommit.MergePullRequestBySquashInput) (*codecommit.MergePullRequestBySquashOutput, error) {
	mergeResult, err := c.MergePullRequestBySquash(m)
	if err != nil {
		return mergeResult, err
	}

	return mergeResult, err
}

func mergeThreeWay(c *codecommit.CodeCommit, m *codecommit.MergePullRequestByThreeWayInput) (*codecommit.MergePullRequestByThreeWayOutput, error) {
	mergeResult, err := c.MergePullRequestByThreeWay(m)
	if err != nil {
		return mergeResult, err
	}

	return mergeResult, err
}
