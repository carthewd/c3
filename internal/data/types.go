package data

import (
	"reflect"
)

type PullRequest struct {
	ID           string
	Title        string
	Owner        string
	Status       string
	SourceBranch string
	DestBranch   string
}

type PullRequests struct {
	PRs []PullRequest
}

type TableData interface {
	GetHeaders() []string
	GetRows() [][]string
}

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
