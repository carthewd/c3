package data

import (
	"reflect"
)

type PR struct {
	ID           string
	Title        string
	Owner        string
	Status       string
	SourceBranch string
	DestBranch   string
}

type PRData struct {
	PRs []PR
}

type TableData interface {
	GetHeaders() []string
	GetRows() [][]string
}

func (p PRData) GetHeaders() []string {
	if len(p.PRs) == 0 {
		var pr []PR
		emptyPR := PR{
			ID: "000",
		}
		pr = append(pr, emptyPR)
		pd := PRData{
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

func (p PRData) GetRows() [][]string {
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
