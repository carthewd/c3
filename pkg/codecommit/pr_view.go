package codecommit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/codecommit"
	log "github.com/sirupsen/logrus"

	"c3/internal/data"
	"c3/util"
)

// ListPRs returns open pull requests for a CodeCommit repository
func ListPRs(c *codecommit.CodeCommit, repoName string, author string, status string) data.PRData {
	prInput := &codecommit.ListPullRequestsInput{
		RepositoryName:    aws.String(repoName),
		PullRequestStatus: aws.String(status),
	}

	result, err := c.ListPullRequests(prInput)
	if err != nil {
		log.Fatal("Could not list pull requests: ", err.Error())
	}

	wg := util.NewBoundWaitGroup(10)
	var allPRs []data.PR

	for _, r := range result.PullRequestIds {
		wg.Add(1)
		go func(c *codecommit.CodeCommit, pr string) {
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

	allPRData := data.PRData{
		PRs: allPRs,
	}

	return allPRData
}

// GetPRDetails describes a CodeCommit pull request object in detail
func GetPRDetails(c *codecommit.CodeCommit, pr string, author string) (data.PR, error) {
	newPR := data.PR{ID: pr}

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
