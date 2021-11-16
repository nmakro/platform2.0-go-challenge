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

	// Register audience routers.
	audiences := assets.PathPrefix("/audiences").Subrouter()
	audiences.HandleFunc("/", m.ListAudience).Methods("GET")
	audiences.HandleFunc("/audience/{id}", m.GetAudience).Methods("GET")
	audiences.HandleFunc("/audience/{id}", m.DeleteAudience).Methods("DELETE")
	audiences.HandleFunc("/audience/{id}", m.UpdateAudience).Methods("PATCH")
	audiences.HandleFunc("/audience", m.AddAudience).Methods("POST")

	// Register charts routers.
	charts := assets.PathPrefix("/charts").Subrouter()
	charts.HandleFunc("/", m.ListCharts).Methods("GET")
	charts.HandleFunc("/chart/{id}", m.GetChart).Methods("GET")
	charts.HandleFunc("/chart/{id}", m.DeleteChart).Methods("DELETE")
	charts.HandleFunc("/chart/{id}", m.UpdateChart).Methods("PATCH")
	charts.HandleFunc("/chart", m.AddChart).Methods("PUT")

	insights := assets.PathPrefix("/insights").Subrouter()
	insights.HandleFunc("/", m.ListInsights).Methods("GET")
	insights.HandleFunc("/insight/{id}", m.GetInsight).Methods("GET")
	insights.HandleFunc("/insight/{id}", m.DeleteInsight).Methods("DELETE")
	insights.HandleFunc("/insight/{id}", m.UpdateInsight).Methods("PATCH")
	insights.HandleFunc("/insight", m.AddInsight).Methods("PUT")

}

func (m *AssetsModule) ListAssets(w http.ResponseWriter, r *http.Request) {
	audiences, err := m.service.ListAudienceAssets(r.Context())

	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, audiences, w)
}
