package assets

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
)

func (m *AssetsModule) AddAudience(w http.ResponseWriter, r *http.Request) {
	req := assets.AddAudienceCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.AddAudience(r.Context(), req.BuildFromCmd()); err != nil {
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

type DeleteAudienceRequest struct {
	AudienceID uint32 `json:"audience_id"`
}

func (m *AssetsModule) DeleteAudience(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	audienceID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an audience id"}, w)
		return
	}

	id, err := strconv.Atoi(audienceID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "audience id must be a uint"}, w)
		return
	}

	if err := m.service.DeleteAudience(r.Context(), uint32(id)); err != nil {
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

func (m *AssetsModule) UpdateAudience(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	audienceID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an audience id"}, w)
		return
	}

	id, err := strconv.Atoi(audienceID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "audience id must be a uint"}, w)
		return
	}

	req := assets.UpdateAudienceCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UpdateAudience(r.Context(), uint32(id), req); err != nil {
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

type GetAudienceRequest struct {
	AudienceID uint32 `json:"id"`
}

func (m *AssetsModule) GetAudience(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	audienceID, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify an audience id"}, w)
		return
	}

	id, err := strconv.Atoi(audienceID)
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "audience id must be a uint"}, w)
		return
	}

	var (
		response assets.Audience
	)

	if response, err = m.service.GetAudience(r.Context(), uint32(id)); err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			fmt.Println("not found!!!!")
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, response, w)
}

func (m *AssetsModule) ListAudience(w http.ResponseWriter, r *http.Request) {
	var (
		response []assets.Audience
		err      error
	)

	if response, err = m.service.GetAllAudienceAssets(r.Context()); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, response, w)
}

type StarAudienceRequest struct {
	UserEmail  string `json:"user_email"`
	AudienceID uint32 `json:"audience_id"`
}

func (m *AssetsModule) StarAudience(w http.ResponseWriter, r *http.Request) {
	req := StarAudienceRequest{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.StarAudience(r.Context(), req.UserEmail, req.AudienceID); err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

type UnStarAudienceRequest struct {
	UserEmail  string `json:"user_email"`
	AudienceID uint32 `json:"audience_id"`
}

func (m *AssetsModule) UnStarAudience(w http.ResponseWriter, r *http.Request) {
	req := UnStarAudienceRequest{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UnstarAudience(r.Context(), req.UserEmail, req.AudienceID); err != nil {
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
