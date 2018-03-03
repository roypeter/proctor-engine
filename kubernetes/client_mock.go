package kubernetes

import (
	"io"

	"github.com/gojektech/proctor-engine/utility"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) ExecuteJob(jobName string, envMap map[string]string) (string, error) {
	args := m.Called(jobName, envMap)
	return args.String(0), args.Error(1)
}

func (m *MockClient) StreamJobLogs(jobName string) (io.ReadCloser, error) {
	args := m.Called(jobName)
	return args.Get(0).(*utility.Buffer), args.Error(1)
}
