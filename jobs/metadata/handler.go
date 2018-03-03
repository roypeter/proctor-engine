package metadata

import (
	"encoding/json"
	"net/http"

	"github.com/gojektech/proctor-engine/logger"
	"github.com/gojektech/proctor-engine/utility"
)

type metadataHandler struct {
	store Store
}

type MetadataHandler interface {
	HandleSubmission() http.HandlerFunc
	HandleBulkDisplay() http.HandlerFunc
}

func NewMetadataHandler(store Store) MetadataHandler {
	return &metadataHandler{
		store: store,
	}
}

func (metadataHandler *metadataHandler) HandleSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var jobMetadata []Metadata
		err := json.NewDecoder(req.Body).Decode(&jobMetadata)
		defer req.Body.Close()
		if err != nil {
			logger.Error("Error parsing request body", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(utility.ClientError))
			return
		}

		for _, metadata := range jobMetadata {
			err = metadataHandler.store.CreateOrUpdateJobMetadata(metadata)
			if err != nil {
				logger.Error("Error updating metadata", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(utility.ServerError))
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (metadataHandler *metadataHandler) HandleBulkDisplay() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		jobMetadata, err := metadataHandler.store.GetAllJobsMetadata()
		if err != nil {
			logger.Error("Error fetching metadata", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		jobsMetadataInJSON, err := json.Marshal(jobMetadata)
		if err != nil {
			logger.Error("Error marshalling jobs metadata in json", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		w.Write(jobsMetadataInJSON)
	}
}
