package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service EnricherService
	log     *logrus.Entry
}

type EnricherService interface {
	HandleUser(ctx context.Context, userName models.RequestEnrich) (models.ResponseEnrich, error)
}

func NewHandler(service EnricherService, log *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log.WithField("module", "handler"),
	}
}

func (h *Handler) enrich(w http.ResponseWriter, r *http.Request) {
	var userName models.RequestEnrich

	err := json.NewDecoder(r.Body).Decode(&userName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	userNameEnriched, err := h.service.HandleUser(r.Context(), userName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		h.log.Infof("err: %s", err)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userNameEnriched)
	if err != nil {
		h.log.Warningf("json.NewEncoder.Encode: %s", err)
	}
}
