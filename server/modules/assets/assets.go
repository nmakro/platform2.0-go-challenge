package assets

import (
	"errors"
	"net/http"

	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
)

type AssetsModule struct {
	service assets.AssetService
}

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
	req := DeleteAudienceRequest{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.DeleteAudience(r.Context(), req.AudienceID); err != nil {
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
	req := assets.UpdateAudienceCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.service.UpdateAudience(r.Context(), req.BuildFromCmd()); err != nil {
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
