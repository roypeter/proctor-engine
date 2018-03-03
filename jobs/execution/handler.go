package execution

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojektech/proctor-engine/jobs/metadata"
	"github.com/gojektech/proctor-engine/jobs/secrets"
	"github.com/gojektech/proctor-engine/kubernetes"
	"github.com/gojektech/proctor-engine/logger"
	"github.com/gojektech/proctor-engine/utility"
)

type executioner struct {
	kubeClient    kubernetes.Client
	metadataStore metadata.Store
	secretsStore  secrets.Store
}

type Executioner interface {
	Handle() http.HandlerFunc
}

func NewExecutioner(kubeClient kubernetes.Client, metadataStore metadata.Store, secretsStore secrets.Store) Executioner {
	return &executioner{
		kubeClient:    kubeClient,
		metadataStore: metadataStore,
		secretsStore:  secretsStore,
	}
}

func (executioner *executioner) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var job Job
		err := json.NewDecoder(req.Body).Decode(&job)
		defer req.Body.Close()
		if err != nil {
			logger.Error("Error parsing request body", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(utility.ClientError))
			return
		}

		jobMetadata, err := executioner.metadataStore.GetJobMetadata(job.Name)
		if err != nil {
			logger.Error("Error finding job to image", job.Name, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		jobSecrets, err := executioner.secretsStore.GetJobSecrets(job.Name)
		if err != nil {
			logger.Error("Error retrieving secrets for job", job.Name, err.Error())

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(utility.ServerError))
			return
		}

		envVars := utility.MergeMaps(job.Args, jobSecrets)
		imageName := jobMetadata.ImageName
		executedJobName, err := executioner.kubeClient.ExecuteJob(imageName, envVars)
		if err != nil {
			logger.Error("Error executing job: %v", job, imageName, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("{ \"name\":\"%s\" }", executedJobName)))

	}
}
