package data

import (
	"reflect"
)

// PullRequest describes a CodeCommit pull request
type PullRequest struct {
	ID           string
	Title        string
	Owner        string
	Status       string
	SourceBranch string
	DestBranch   string
}

// PullRequests is a list of PullRequest
type PullRequests struct {
	PRs []PullRequest
}

// PullRequestDiff describes fields for diff operations
type PullRequestDiff struct {
	DestCommit   string
	MergeCommit  string
	SourceBranch string
}

// NewPullRequest describes a new pull request to be created
type NewPullRequest struct {
	Title          string
	Description    string
	Repository     string
	SourceRef      string
	DestinationRef string
}

// Path describes fields used to identify a git filesystem object
type Path struct {
	Path     string
	PathType string
}

type MergeOptions struct {
	FF       bool
	Squash   bool
	ThreeWay bool
}

type MergeInput struct {
	Type               string
	PRID               string
	SourceCommit       string
	Repository         string
	ConflictDetail     string
	ConflictResolution string
	AuthorName         string
	AuthorEmail        string
}

// TableData Interface for dynamically creating tables with table_maker.go
type TableData interface {
	GetHeaders() []string
	GetRows() [][]string
}

// GetHeaders implements the TableData interface to generate tables for PullRequest objects
func (p PullRequests) GetHeaders() []string {
	if len(p.PRs) == 0 {
		var pr []PullRequest
		emptyPR := PullRequest{
			ID: "000",
		}
		pr = append(pr, emptyPR)
		pd := PullRequests{
			pr,
		}
		p = pd
	}

	val := reflect.ValueOf(p.PRs[0])
	var headers []string
	for i := 0; i < val.Type().NumField(); i++ {
		headers = append(headers, val.Type().Field(i).Name)
	}

	return headers
}

// GetRows implements the TableData interface to generate tables for PullRequest objects
func (p PullRequests) GetRows() [][]string {
	var allRows [][]string
	for _, pr := range p.PRs {
		var newRow []string
		v := reflect.ValueOf(pr)
		for i := 0; i < v.NumField(); i++ {
			newRow = append(newRow, v.Field(i).Interface().(string))
		}
		allRows = append(allRows, newRow)
	}

	return allRows
}
