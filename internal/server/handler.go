package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/AlexZav1327/name-enricher/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

const defaultLimit = 20

type Handler struct {
	service EnricherService
	log     *logrus.Entry
	metrics *metrics
}

type EnricherService interface {
	EnrichUser(ctx context.Context, userName models.RequestEnrich) (models.ResponseEnrich, error)
	GetUsersList(ctx context.Context, params models.ListingQueryParams) ([]models.ResponseEnrich, error)
	UpdateUser(ctx context.Context, user models.ResponseEnrich) (models.ResponseEnrich, error)
	DeleteUser(ctx context.Context, userName string) error
}

func NewHandler(service EnricherService, log *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log.WithField("module", "handler"),
		metrics: newMetrics(),
	}
}

func (h *Handler) enrich(w http.ResponseWriter, r *http.Request) {
	var userName models.RequestEnrich

	err := json.NewDecoder(r.Body).Decode(&userName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	userNameEnriched, err := h.service.EnrichUser(r.Context(), userName)
	if errors.Is(err, models.ErrNameNotValid) {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userNameEnriched)
	if err != nil {
		h.log.Warningf("json.NewEncoder(w).Encode(userNameEnriched): %s", err)
	}
}

func (h *Handler) getList(w http.ResponseWriter, r *http.Request) {
	var params models.ListingQueryParams

	params.TextFilter = r.URL.Query().Get("textFilter")

	params.ItemsPerPage, _ = strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	if params.ItemsPerPage == 0 {
		params.ItemsPerPage = defaultLimit
	}

	params.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	params.Sorting = r.URL.Query().Get("sorting")
	params.Descending, _ = strconv.ParseBool(r.URL.Query().Get("descending"))

	usersList, err := h.service.GetUsersList(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(usersList)
	if err != nil {
		h.log.Warningf("json.NewEncoder(w).Encode(usersList): %s", err)
	}
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var user models.ResponseEnrich

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	user.Name = chi.URLParam(r, "name")

	updatedUser, err := h.service.UpdateUser(r.Context(), user)
	if errors.Is(err, storage.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		h.log.Warningf("json.NewEncoder(w).Encode(updatedUser): %s", err)
	}
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "name")

	err := h.service.DeleteUser(r.Context(), userName)
	if errors.Is(err, storage.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
