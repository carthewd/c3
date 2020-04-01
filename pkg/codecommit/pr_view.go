package codecommit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"

	log "github.com/sirupsen/logrus"

	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/util"
)

// ListPRs returns open pull requests for a CodeCommit repository
func ListPRs(c codecommitiface.CodeCommitAPI, repoName string, author string, status string) data.PullRequests {
	prInput := &codecommit.ListPullRequestsInput{
		RepositoryName:    aws.String(repoName),
		PullRequestStatus: aws.String(status),
	}

	result, err := c.ListPullRequests(prInput)
	if err != nil {
		log.Fatal("Could not list pull requests: ", err.Error())
	}

	wg := util.NewBoundWaitGroup(10)
	var allPRs []data.PullRequest

	for _, r := range result.PullRequestIds {
		wg.Add(1)
		go func(c codecommitiface.CodeCommitAPI, pr string) {
			defer wg.Done()
			newPR, err := GetPRDetails(c, pr, author)

			if err == nil {
				if newPR.ID != "" {
					allPRs = append(allPRs, newPR)
				}
			}
		}(c, *r)
	}
	wg.Wait()

	allPRData := data.PullRequests{
		PRs: allPRs,
	}

	return allPRData
}

// GetPRDetails describes a CodeCommit pull request object in detail
func GetPRDetails(c codecommitiface.CodeCommitAPI, pr string, author string) (data.PullRequest, error) {
	newPR := data.PullRequest{ID: pr}

	prDetailInput := &codecommit.GetPullRequestInput{
		PullRequestId: aws.String(pr),
	}

	detailResult, err := c.GetPullRequest(prDetailInput)
	if err != nil {
		fmt.Println(err)
		return newPR, err
	}

	newPR.Title = *detailResult.PullRequest.Title

	prAuthor, _ := arn.Parse(*detailResult.PullRequest.AuthorArn)
	var re = regexp.MustCompile(`.*\/`)
	auth := re.ReplaceAllString(prAuthor.Resource, `$2`)

	newPR.Owner = auth
	newPR.Status = *detailResult.PullRequest.PullRequestStatus
	targets := (*detailResult.PullRequest).PullRequestTargets

	newPR.SourceBranch = strings.Replace(*targets[0].SourceReference, "refs/heads/", "", -1)
	newPR.DestBranch = strings.Replace(*targets[0].DestinationReference, "refs/heads/", "", -1)

	if author != "" && strings.Contains(prAuthor.Resource, author) {
		return newPR, nil
	} else if author != "" && !strings.Contains(prAuthor.Resource, author) {
		err := errors.New("No matching author")
		return newPR, err
	}

	return newPR, nil
}

// GetPRCommits uses the CodeCommit API to get relevant commit hashes for a diff
func GetPRCommits(c codecommitiface.CodeCommitAPI, pr string) (data.PullRequestDiff, error) {
	prInput := &codecommit.GetPullRequestInput{
		PullRequestId: aws.String(pr),
	}

	result, err := c.GetPullRequest(prInput)
	if err != nil {
		log.Error(err)
		return data.PullRequestDiff{}, err
	}

	prTargets := (*result.PullRequest).PullRequestTargets

	prCommits := data.PullRequestDiff{
		DestCommit:  *prTargets[0].DestinationCommit,
		MergeCommit: *prTargets[0].SourceCommit,
		SourceBranch: strings.Replace(*prTargets[0].SourceReference, "refs/heads/", "", -1),
	}

	return prCommits, err
}
