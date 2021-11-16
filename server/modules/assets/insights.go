package assets

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
)

func (m *AssetsModule) AddInsight(w http.ResponseWriter, r *http.Request) {
	req := assets.AddInsightCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.AddInsight(r.Context(), req.BuildFromCmd()); err != nil {
		var duplicateErr *app.ErrDuplicateEntry
		if errors.As(err, &duplicateErr) {
			gwihttp.ResponseWithJSON(http.StatusConflict, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		var validationErr *assets.ErrValidation
		if errors.As(err, &validationErr) {
			gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

func (m *AssetsModule) DeleteInsight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	insightID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an insight id"}, w)
		return
	}

	id, err := strconv.Atoi(insightID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "insight id must be a uint"}, w)
		return
	}

	if err := m.service.DeleteInsight(r.Context(), uint32(id)); err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusNoContent, nil, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusNoContent, nil, w)
}

func (m *AssetsModule) UpdateInsight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	insightID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an insight id"}, w)
		return
	}

	id, err := strconv.Atoi(insightID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "insight id must be a uint"}, w)
		return
	}

	req := assets.UpdateInsightCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UpdateInsight(r.Context(), uint32(id), req); err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		var validationErr *assets.ErrValidation
		if errors.As(err, &validationErr) {
			gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

func (m *AssetsModule) GetInsight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	insightID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an insight id"}, w)
		return
	}

	id, err := strconv.Atoi(insightID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "insight id must be a uint"}, w)
		return
	}

	var (
		response assets.Insight
	)

	if response, err = m.service.GetInsight(r.Context(), uint32(id)); err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, response, w)
}

func (m *AssetsModule) ListInsights(w http.ResponseWriter, r *http.Request) {
	var (
		response []assets.Insight
		err      error
	)

	if response, err = m.service.ListInsightAssets(r.Context()); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, response, w)
}

type StarInsightRequest struct {
	UserEmail string `json:"user_email"`
	InsightID uint32 `json:"insight_id"`
}

func (m *AssetsModule) StarInsight(w http.ResponseWriter, r *http.Request) {
	req := StarInsightRequest{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.StartInsight(r.Context(), req.UserEmail, req.InsightID); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

type UnStarInsightRequest struct {
	UserEmail string `json:"user_email"`
	InsightID uint32 `json:"insight_id"`
}

func (m *AssetsModule) UnStarInsight(w http.ResponseWriter, r *http.Request) {
	req := UnStarInsightRequest{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UnstarInsight(r.Context(), req.UserEmail, req.InsightID); err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, notFound) {
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}
