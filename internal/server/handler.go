package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service EnrichService
	log     *logrus.Entry
}

type EnrichService interface {
	Enrich(ctx context.Context, personName models.RequestEnrich) (models.ResponseEnrich, error)
}

func NewHandler(service EnrichService, log *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log.WithField("module", "handler"),
	}
}

func (h *Handler) enrich(w http.ResponseWriter, r *http.Request) {
	var personName models.RequestEnrich

	err := json.NewDecoder(r.Body).Decode(&personName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	enrichedPersonName, err := h.service.Enrich(r.Context(), personName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(enrichedPersonName)
	if err != nil {
		h.log.Warningf("json.NewEncoder.Encode: %s", err)
	}
}
