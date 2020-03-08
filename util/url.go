package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/pkg/gitconfig"
)

// CreatePullRequestURL generates CodeCommit console URLs
func CreatePullRequestURL(repo string, prID string) string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}

	url := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s/pull-requests/%s/details", region, repo, prID)

	return url
}

// CreatePathURL generates CodeCommit console URLs from filenames
func CreatePathURL(p data.Path) string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}

	if p.PathType == "file" {
		filePath, _ := gitconfig.GitCmd("ls-files", p.Path, "--full-name")
		repo, _ := gitconfig.GetOrigin()
		branch, _ := gitconfig.GitCmd("rev-parse", "--abbrev-ref", "HEAD")

		url := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s/browse/refs/heads/%s/--/%s", region, repo, strings.Replace(branch, "\n", "", 1), filePath)

		return url
	}
	return ""
}
