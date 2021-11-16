package assets

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
)

type AssetsModule struct {
	service      *assets.AssetService
	sessionStore *sessions.CookieStore
}

func Setup(router *mux.Router, service *assets.AssetService, sessionStore *sessions.CookieStore) {
	m := &AssetsModule{
		service:      service,
		sessionStore: sessionStore,
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
	charts.HandleFunc("/chart", m.AddChart).Methods("POST")

	insights := assets.PathPrefix("/insights").Subrouter()
	insights.HandleFunc("/", m.ListInsights).Methods("GET")
	insights.HandleFunc("/insight/{id}", m.GetInsight).Methods("GET")
	insights.HandleFunc("/insight/{id}", m.DeleteInsight).Methods("DELETE")
	insights.HandleFunc("/insight/{id}", m.UpdateInsight).Methods("PATCH")
	insights.HandleFunc("/insight", m.AddInsight).Methods("POST")

	stars := router.PathPrefix("starred-assets").Subrouter()
	stars.HandleFunc("/", m.ListFavoritesAssetsForUser).Methods("GET")
	stars.HandleFunc("/audience/{id}", m.LoggedIn(m.StarAudience)).Methods("PUT")
	stars.HandleFunc("/audience/{id}", m.LoggedIn(m.StarAudience)).Methods("DELETE")
	stars.HandleFunc("/insight/{id}", m.LoggedIn(m.StarAudience)).Methods("PUT")
	stars.HandleFunc("/insight/{id}", m.LoggedIn(m.StarAudience)).Methods("DELETE")
	stars.HandleFunc("/chart/{id}", m.LoggedIn(m.StarAudience)).Methods("PUT")
	stars.HandleFunc("/chart/{id}", m.LoggedIn(m.StarAudience)).Methods("DELETE")
}

func (m *AssetsModule) ListAssets(w http.ResponseWriter, r *http.Request) {
	var (
		wg sync.WaitGroup
	)
	wg.Add(3)

	errChan := make(chan error, 3)

	audiences := make([]assets.Audience, 0, 50)
	go func() {
		defer wg.Done()
		aud, err := m.service.ListAudienceAssets(r.Context())

		if err != nil {
			errChan <- err
			return
		}

		for i := range audiences {
			audiences = append(audiences, aud[i])
		}
	}()

	chartsSlice := make([]assets.Chart, 0, 50)
	go func() {
		defer wg.Done()
		charts, err := m.service.ListChartAssets(r.Context())

		if err != nil {
			errChan <- err
			return
		}

		for i := range charts {
			chartsSlice = append(charts, charts[i])
		}
	}()

	insights := make([]assets.Insight, 0, 50)
	go func() {
		defer wg.Done()
		ins, err := m.service.ListInsightAssets(r.Context())

		if err != nil {
			errChan <- err
			return
		}

		for i := range insights {
			insights = append(insights, ins[i])
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
			return
		}
	}

	result := struct {
		Audiences []assets.Audience `json:"audience"`
		Charts    []assets.Chart    `json:"charts"`
		Insights  []assets.Insight  `json:"insights"`
	}{
		Audiences: audiences,
		Charts:    chartsSlice,
		Insights:  insights,
	}

	gwihttp.ResponseWithJSON(http.StatusOK, result, w)
}

func (m *AssetsModule) ListFavoritesAssetsForUser(w http.ResponseWriter, r *http.Request) {
	session, _ := m.sessionStore.Get(r, "gwi-cookie")
	userEmal, ok := session.Values["user_email"].(string)
	if !ok || userEmal == "" {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, nil, w)
		return
	}

	var (
		wg sync.WaitGroup
	)
	wg.Add(3)

	errChan := make(chan error, 3)

	audiences := make([]assets.Audience, 0, 50)
	go func() {
		defer wg.Done()
		aud, err := m.service.GetAudiencesForUser(r.Context(), userEmal)

		if err != nil {
			errChan <- err
			return
		}

		for i := range audiences {
			audiences = append(audiences, aud[i])
		}
	}()

	chartsSlice := make([]assets.Chart, 0, 50)
	go func() {
		defer wg.Done()
		charts, err := m.service.GetChartsForUser(r.Context(), userEmal)

		if err != nil {
			errChan <- err
			return
		}

		for i := range charts {
			chartsSlice = append(charts, charts[i])
		}
	}()

	insights := make([]assets.Insight, 0, 50)
	go func() {
		defer wg.Done()
		ins, err := m.service.GetInsightsForUser(r.Context(), userEmal)

		if err != nil {
			errChan <- err
			return
		}

		for i := range insights {
			insights = append(insights, ins[i])
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			var notFound *app.ErrEntityNotFound
			if errors.As(err, &notFound) {
				gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": "cannot find logged in user"}, w)
				return
			}
			gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
			return
		}
	}

	result := struct {
		Audiences []assets.Audience `json:"audience"`
		Charts    []assets.Chart    `json:"charts"`
		Insights  []assets.Insight  `json:"insights"`
	}{
		Audiences: audiences,
		Charts:    chartsSlice,
		Insights:  insights,
	}

	gwihttp.ResponseWithJSON(http.StatusOK, result, w)
}

func (m *AssetsModule) LoggedIn(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := m.sessionStore
		session, _ := store.Get(r, "gwi-cookie")
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			nextHandler(w, r)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}
