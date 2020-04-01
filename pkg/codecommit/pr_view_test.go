package codecommit

import (
	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/internal/mocks"
	"sort"
	"strconv"
	"testing"
)

func Test_PRView(t *testing.T) {
	for _, test := range []struct {
		Type   string
		Args   []string
		Output data.PullRequests
	}{
		{
			Type: "list",
			Args: []string{"test", "", "OPEN"},
			Output: data.PullRequests{
				PRs:
				[]data.PullRequest{
					data.PullRequest{
						ID:           "2",
						Title:        "Pronunciation difficulty analyzer",
						Owner:        "Jane_Doe",
						Status:       "OPEN",
						SourceBranch: "jane-branch",
						DestBranch:   "master",
					},
					data.PullRequest{
						ID:           "12",
						Title:        "Build Beowulf cluster",
						Owner:        "John_Smith",
						Status:       "OPEN",
						SourceBranch: "john-branch",
						DestBranch:   "master",
					},
					data.PullRequest{
						ID:           "16",
						Title:        "Ph'nglui mglw'nafh Cthulhu R'lyeh wgah'nagl fhtagn",
						Owner:        "Cthuluhu",
						Status:       "OPEN",
						SourceBranch: "cthulhu-branch",
						DestBranch:   "master",
					},
				},
			},
		},
		{
			Type: "details",
			Args: []string{"16", ""},
			Output: data.PullRequests{
				PRs:
				[]data.PullRequest{
					data.PullRequest{
						ID:           "16",
						Title:        "Ph'nglui mglw'nafh Cthulhu R'lyeh wgah'nagl fhtagn",
						Owner:        "Cthulhu",
						Status:       "OPEN",
						SourceBranch: "cthulhu-branch",
						DestBranch:   "master",
					},
				},
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			switch test.Type {
			case "list":
				sort.SliceStable(test.Output.PRs, func(i, j int) bool {
					leftId, _ := strconv.Atoi(test.Output.PRs[i].ID)
					rightId, _ := strconv.Atoi(test.Output.PRs[j].ID)

					return leftId < rightId
				})

				svc := &mocks.MockCodeCommitClient{}

				actual := ListPRs(svc, test.Args[0], test.Args[1], test.Args[2])
				sort.SliceStable(actual.PRs, func(i, j int) bool {
					leftId, _ := strconv.Atoi(actual.PRs[i].ID)
					rightId, _ := strconv.Atoi(actual.PRs[j].ID)

					return leftId < rightId
				})

				if actual.PRs[0] != test.Output.PRs[0] {
					t.Errorf("expected %s but got %s", test.Output, actual)
				}

			case "details":
				svc := &mocks.MockCodeCommitClient{}

				actual, err := GetPRDetails(svc, test.Args[0], test.Args[1])
				if err != nil {
					t.Errorf("unexpected error occured: %s", err)
				}

				if actual != test.Output.PRs[0] {
					t.Errorf("expected %s but got %s", test.Output.PRs[0], actual)
				}
			}
		})
	}

	for _, test := range []struct {
		Args   string
		Output data.PullRequestDiff
	}{
		{
			Args: "16",
			Output: data.PullRequestDiff{
				DestCommit:   "5d036259EXAMPLE",
				MergeCommit:  "317f8570EXAMPLE",
				SourceBranch: "cthulhu-branch",
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			svc := &mocks.MockCodeCommitClient{}

			actual, err := GetPRCommits(svc, test.Args)
			if err != nil {
				t.Errorf("unexpected error occured: %s", err)
			}

			if actual != test.Output {
				t.Errorf("expected %s but got %s", test.Output, actual)
			}

		})
	}
}
