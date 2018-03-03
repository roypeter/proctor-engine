package secrets

import (
	"encoding/json"
	"net/http"

	"github.com/gojektech/proctor-engine/logger"
	"github.com/gojektech/proctor-engine/utility"
)

type secretsHandler struct {
	secretsStore Store
}

type SecretsHandler interface {
	HandleSubmission() http.HandlerFunc
}

func NewSecretsHandler(secretsStore Store) SecretsHandler {
	return &secretsHandler{
		secretsStore: secretsStore,
	}
}

func (secretsHandler *secretsHandler) HandleSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var secret Secret
		err := json.NewDecoder(req.Body).Decode(&secret)
		defer req.Body.Close()
		if err != nil {
			logger.Error("Error parsing request body", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(utility.ClientError))
			return
		}

		err = secretsHandler.secretsStore.CreateOrUpdateJobSecret(secret)
		if err != nil {
			logger.Error("Error updating secrets", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
