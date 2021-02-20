package codecommit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/util"
	log "github.com/sirupsen/logrus"
)

func ListRepoNames(c codecommitiface.CodeCommitAPI) []string {
	repoInput := &codecommit.ListRepositoriesInput{
		SortBy: aws.String("repositoryName"),
	}

	result, err := c.ListRepositories(repoInput)
	if err != nil {
		log.Fatal("Could not list repositories: ", err.Error())
	}

	var repos []string
	for _, r := range result.Repositories {
		repos = append(repos, *r.RepositoryName)
	}

	return repos
}

func GetRepoDetails(c codecommitiface.CodeCommitAPI, repos []string) data.AllRepoDetails {
	var allRepos []data.RepoDetails
	wg := util.NewBoundWaitGroup(10)

	for _, repo := range repos {
		wg.Add(1)
		go func(c codecommitiface.CodeCommitAPI, repo string) {
			defer wg.Done()

			repoInput := &codecommit.GetRepositoryInput{
				RepositoryName: aws.String(repo),
			}
			repoDetails, err := c.GetRepository(repoInput)
			if err != nil {
				log.Fatal("Could not get repository details: ", err.Error())
			}

			repoD := data.RepoDetails{
				Name: *repoDetails.RepositoryMetadata.RepositoryName,
				//DefaultBranch: *repoDetails.RepositoryMetadata.DefaultBranch,
				CloneSSH:     *repoDetails.RepositoryMetadata.CloneUrlSsh,
				CloneHTTP:    *repoDetails.RepositoryMetadata.CloneUrlHttp,
				LastModified: repoDetails.RepositoryMetadata.LastModifiedDate.String(),
				//Description:  *repoDetails.RepositoryMetadata.RepositoryDescription,
			}

			allRepos = append(allRepos, repoD)
		}(c, repo)
	}
	wg.Wait()

	allRepoData := data.AllRepoDetails{
		Repos: allRepos,
	}

	return allRepoData
}
