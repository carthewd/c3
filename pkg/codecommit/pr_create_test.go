package codecommit

import (
	"testing"

	"github.com/carthewd/c3/internal/data"
	"github.com/carthewd/c3/internal/mocks"
)

func Test_PRCreate(t *testing.T) {
	for _, test := range []struct {
		Type   string
		Args   data.NewPullRequest
		Output string
	}{
		{
			Type: "create",
			Args: data.NewPullRequest{
				Title:          "Pronunciation difficulty analyzer",
				Description:    "A code review of the new feature I just added to the service.",
				Repository:     "MyDemoRepo",
				SourceRef:      "refs/heads/jane-branch",
				DestinationRef: "refs/heads/master",
			},
			Output: "2",
		},
	} {
		t.Run("", func(t *testing.T) {
			switch test.Type {
			case "create":
				svc := &mocks.MockCodeCommitClient{}

				actual, err := CreatePR(svc, test.Args)
				if err != nil {
					t.Errorf("unexpected error occured: %s", err)
				}

				if actual != test.Output {
					t.Errorf("expected %s but got %s", test.Output, actual)
				}
			}
		})
	}
}
