package util

import (
	"fmt"
	"os"
)

func CreatePullRequestURL(repo string, prID string) string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}

	url := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s/pull-requests/%s/details", region, repo, prID)

	return url
}

func CreatePathURL() {
	// https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/old-mutual_platform_terraform/browse/refs/heads/master/--/terraform/modules/eks/nodes.tf
}
