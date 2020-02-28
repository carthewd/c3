package gitconfig

import (
	"os/exec"
	"strings"
)

type GitURL struct {
	cloneURL  string
	cloneType string
}

func GitCmd(c ...string) (string, error) {
	o, err := exec.Command("git", c...).Output()

	output := string(o[:])
	return output, err
}

func GetGitURL() (GitURL, error) {
	args := []string{"config", "remote.origin.url"}

	url, err := GitCmd(args...)
	if err != nil {
		url := GitURL{}
		return url, err
	}

	var cloneType string
	if strings.Contains(url, "https://") {
		cloneType = "https"
	} else if strings.Contains(url, "ssh://") {
		cloneType = "ssh"
	}

	parsedURL := GitURL{
		cloneURL:  url,
		cloneType: cloneType,
	}

	return parsedURL, err
}
