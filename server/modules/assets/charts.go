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

func (m *AssetsModule) AddChart(w http.ResponseWriter, r *http.Request) {
	req := assets.AddChartCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.AddChart(r.Context(), req.BuildFromCmd()); err != nil {
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

func (m *AssetsModule) DeleteChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chartID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify a chart id"}, w)
		return
	}

	id, err := strconv.Atoi(chartID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "chart id must be a uint"}, w)
		return
	}

	if err := m.service.DeleteChart(r.Context(), uint32(id)); err != nil {
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

func (m *AssetsModule) UpdateChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chartID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify a chart id"}, w)
		return
	}

	id, err := strconv.Atoi(chartID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "chart id must be a uint"}, w)
		return
	}

	req := assets.UpdateChartCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UpdateChart(r.Context(), uint32(id), req); err != nil {
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

func (m *AssetsModule) GetChart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	audienceID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify a chart id"}, w)
		return
	}

	id, err := strconv.Atoi(audienceID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "chart id must be a uint"}, w)
		return
	}

	var (
		response assets.Chart
	)

	if response, err = m.service.GetChart(r.Context(), uint32(id)); err != nil {
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

func (m *AssetsModule) ListCharts(w http.ResponseWriter, r *http.Request) {
	var (
		response []assets.Chart
		err      error
	)

	if response, err = m.service.ListChartAssets(r.Context()); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, response, w)
}

func (m *AssetsModule) StarChart(w http.ResponseWriter, r *http.Request) {
	session, _ := m.sessionStore.Get(r, "gwi-cookie")
	userEmail, ok := session.Values["user_email"].(string)
	if !ok || userEmail == "" {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, nil, w)
		return
	}

	vars := mux.Vars(r)
	chartID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an audience id"}, w)
		return
	}

	id, err := strconv.Atoi(chartID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "audience id must be a uint"}, w)
		return
	}

	if err := m.service.StarChart(r.Context(), userEmail, uint32(id)); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

func (m *AssetsModule) UnStarChart(w http.ResponseWriter, r *http.Request) {
	session, _ := m.sessionStore.Get(r, "gwi-cookie")
	userEmail, ok := session.Values["user_email"].(string)
	if !ok || userEmail == "" {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, nil, w)
		return
	}

	vars := mux.Vars(r)
	chartID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an audience id"}, w)
		return
	}

	id, err := strconv.Atoi(chartID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "audience id must be a uint"}, w)
		return
	}

	if err := m.service.UnstarChart(r.Context(), userEmail, uint32(id)); err != nil {
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
