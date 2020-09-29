package mocks

import (
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
)

type mockListPullRequestsOutput struct {
	NextToken      *string   `json: "nextToken"`
	PullRequestIds []*string `json: "pullRequestIds"`
}

type MockCodeCommitClient struct {
	codecommitiface.CodeCommitAPI
}

func (m *MockCodeCommitClient) ListPullRequests(input *codecommit.ListPullRequestsInput) (*codecommit.ListPullRequestsOutput, error) {
	file, _ := ioutil.ReadFile("../../tests/fixtures/pr_list.json")
	prListOutput := mockListPullRequestsOutput{}

	err := json.Unmarshal([]byte(file), &prListOutput)

	mockOutput := &codecommit.ListPullRequestsOutput{
		NextToken:      prListOutput.NextToken,
		PullRequestIds: prListOutput.PullRequestIds,
	}

	return mockOutput, err
}

func (m *MockCodeCommitClient) GetPullRequest(input *codecommit.GetPullRequestInput) (*codecommit.GetPullRequestOutput, error) {
	file, _ := ioutil.ReadFile("../../tests/fixtures/pr_details_new.json")

	var mockOutput []*codecommit.GetPullRequestOutput
	err := json.Unmarshal([]byte(file), &mockOutput)

	for i, v := range mockOutput {
		if *v.PullRequest.PullRequestId == *input.PullRequestId {
			return mockOutput[i], err
		}
	}
	return mockOutput[0], err
}

func (m *MockCodeCommitClient) CreatePullRequest(input *codecommit.CreatePullRequestInput) (*codecommit.CreatePullRequestOutput, error) {
	file, _ := ioutil.ReadFile("../../tests/fixtures/pr_details_new.json")

	var mockOutput []*codecommit.CreatePullRequestOutput
	err := json.Unmarshal([]byte(file), &mockOutput)

	return mockOutput[0], err
}
