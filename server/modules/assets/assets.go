package assets

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
)

type AssetsModule struct {
	service *assets.AssetService
}

func Setup(router *mux.Router, service *assets.AssetService) {
	m := &AssetsModule{
		service: service,
	}

	assets := router.PathPrefix("/assets").Subrouter()
	assets.HandleFunc("/", m.ListAssets).Methods("GET")

	audiences := assets.PathPrefix("/audiences").Subrouter()
	audiences.HandleFunc("/", m.ListAudience).Methods("GET")
	audiences.HandleFunc("/audience/{id}", m.GetAudience).Methods("GET")
	audiences.HandleFunc("/audience/{id}", m.DeleteAudience).Methods("DELETE")
	audiences.HandleFunc("/audience", m.AddAudience).Methods("POST")
	audiences.HandleFunc("/audience/{id}", m.UpdateAudience).Methods("PATCH")

	//charts := assets.PathPrefix("/charts").Subrouter()
	//charts.HandleFunc("/", m.Li).Methods("GET")
	// charts.HandleFunc("/audience/{id}", m.GetAudience).Methods("GET")
	// charts.HandleFunc("/audience/{id}", m.DeleteAudience).Methods("GET")
	// charts.HandleFunc("/audience", m.UpdateAudience).Methods("PUT")
	// charts.HandleFunc("/audience/{id}", m.UpdateAudience).Methods("PATCH")

}

func (m *AssetsModule) ListAssets(w http.ResponseWriter, r *http.Request) {
	audiences, err := m.service.GetAllAudienceAssets(r.Context())

	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, audiences, w)
}
