package execution

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojektech/proctor-engine/jobs/metadata"
	"github.com/gojektech/proctor-engine/jobs/secrets"
	"github.com/gojektech/proctor-engine/kubernetes"
	"github.com/gojektech/proctor-engine/utility"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExecutionerTestSuite struct {
	suite.Suite
	mockKubeClient    kubernetes.MockClient
	mockMetadataStore *metadata.MockStore
	mockSecretsStore  *secrets.MockStore
	testExecutioner   Executioner
}

func (suite *ExecutionerTestSuite) SetupTest() {
	suite.mockKubeClient = kubernetes.MockClient{}
	suite.mockMetadataStore = &metadata.MockStore{}
	suite.mockSecretsStore = &secrets.MockStore{}
	suite.testExecutioner = NewExecutioner(&suite.mockKubeClient, suite.mockMetadataStore, suite.mockSecretsStore)
}

func (suite *ExecutionerTestSuite) TestSuccessfulJobExecution() {
	t := suite.T()

	jobName := "sample-job-name"
	jobArgs := map[string]string{
		"argOne": "sample-arg",
		"argTwo": "another-arg",
	}
	job := Job{
		Name: jobName,
		Args: jobArgs,
	}

	requestBody, err := json.Marshal(job)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/execute", bytes.NewReader(requestBody))
	responseRecorder := httptest.NewRecorder()

	jobMetadata := metadata.Metadata{
		ImageName: "img",
	}
	suite.mockMetadataStore.On("GetJobMetadata", jobName).Return(&jobMetadata, nil).Once()

	jobSecrets := map[string]string{
		"secretOne": "sample-secrets",
		"secretTwo": "another-secret",
	}
	suite.mockSecretsStore.On("GetJobSecrets", jobName).Return(jobSecrets, nil).Once()

	executedJobName := "proctor-ipsum-lorem"
	envVarsForImage := utility.MergeMaps(jobArgs, jobSecrets)
	suite.mockKubeClient.On("ExecuteJob", jobMetadata.ImageName, envVarsForImage).Return(executedJobName, nil).Once()

	suite.testExecutioner.Handle()(responseRecorder, req)

	suite.mockMetadataStore.AssertExpectations(t)
	suite.mockSecretsStore.AssertExpectations(t)
	suite.mockKubeClient.AssertExpectations(t)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, fmt.Sprintf("{ \"name\":\"%s\" }", executedJobName), responseRecorder.Body.String())
}

func (suite *ExecutionerTestSuite) TestJobExecutionOnMalformedRequest() {
	t := suite.T()

	jobExecutionRequest := fmt.Sprintf("{ some-malformed-request }")
	req := httptest.NewRequest("POST", "/execute", bytes.NewReader([]byte(jobExecutionRequest)))
	responseRecorder := httptest.NewRecorder()

	suite.testExecutioner.Handle()(responseRecorder, req)

	suite.mockMetadataStore.AssertNotCalled(t, "GetJobMetadata", mock.Anything)
	suite.mockSecretsStore.AssertNotCalled(t, "GetJobSecrets", mock.Anything)
	suite.mockKubeClient.AssertNotCalled(t, "ExecuteJob", mock.Anything, mock.Anything)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, utility.ClientError, responseRecorder.Body.String())
}

func (suite *ExecutionerTestSuite) TestJobExecutionOnImageLookupFailuer() {
	t := suite.T()

	jobName := "sample-job-name"
	job := Job{
		Name: jobName,
	}
	requestBody, err := json.Marshal(job)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/execute", bytes.NewReader(requestBody))
	responseRecorder := httptest.NewRecorder()

	suite.mockMetadataStore.On("GetJobMetadata", jobName).Return(&metadata.Metadata{}, errors.New("No image found for job name")).Once()

	suite.testExecutioner.Handle()(responseRecorder, req)

	suite.mockMetadataStore.AssertExpectations(t)
	suite.mockSecretsStore.AssertNotCalled(t, "GetJobSecrets", mock.Anything)
	suite.mockKubeClient.AssertNotCalled(t, "ExecuteJob", mock.Anything, mock.Anything)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, utility.ServerError, responseRecorder.Body.String())
}

func (suite *ExecutionerTestSuite) TestJobExecutionOnSecretsFetchFailuer() {
	t := suite.T()

	jobName := "sample-job-name"
	job := Job{
		Name: jobName,
	}
	requestBody, err := json.Marshal(job)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/execute", bytes.NewReader(requestBody))
	responseRecorder := httptest.NewRecorder()

	suite.mockMetadataStore.On("GetJobMetadata", jobName).Return(&metadata.Metadata{}, nil).Once()

	emptyMap := make(map[string]string)
	suite.mockSecretsStore.On("GetJobSecrets", jobName).Return(emptyMap, errors.New("secrets fetch error")).Once()

	suite.testExecutioner.Handle()(responseRecorder, req)

	suite.mockMetadataStore.AssertExpectations(t)
	suite.mockSecretsStore.AssertExpectations(t)
	suite.mockKubeClient.AssertNotCalled(t, "ExecuteJob", mock.Anything, mock.Anything, mock.Anything)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
	assert.Equal(t, utility.ServerError, responseRecorder.Body.String())
}

func (suite *ExecutionerTestSuite) TestJobExecutionOnExecutionFailure() {
	t := suite.T()

	jobName := "sample-job-name"
	emptyMap := make(map[string]string)
	job := Job{
		Name: jobName,
		Args: emptyMap,
	}

	requestBody, err := json.Marshal(job)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/execute", bytes.NewReader(requestBody))
	responseRecorder := httptest.NewRecorder()

	jobMetadata := metadata.Metadata{
		ImageName: "img",
	}
	suite.mockMetadataStore.On("GetJobMetadata", jobName).Return(&jobMetadata, nil).Once()

	suite.mockSecretsStore.On("GetJobSecrets", jobName).Return(emptyMap, nil).Once()

	suite.mockKubeClient.On("ExecuteJob", jobMetadata.ImageName, emptyMap).Return("", errors.New("Kube client job execution error")).Once()

	suite.testExecutioner.Handle()(responseRecorder, req)

	suite.mockMetadataStore.AssertExpectations(t)
	suite.mockSecretsStore.AssertExpectations(t)
	suite.mockKubeClient.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, utility.ServerError, responseRecorder.Body.String())
}

func TestExecutionerTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionerTestSuite))
}
