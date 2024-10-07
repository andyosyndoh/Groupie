package routes

import (
	"net/http"
	"strings"

	"groupie/internals/handlers"
	"groupie/internals/renders"
)

// Allowed routes
var allowedRoutes = map[string]bool{
	"/":        true,
	"/details": true,
	"/api/search": true,
}

// RouteChecker is a middleware that checkes allowed routes
func RouteChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := allowedRoutes[r.URL.Path]; !ok {
			handlers.NotFoundHandler(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes manages the routes
func RegisterRoutes(mux *http.ServeMux) {
	staticDir := renders.GetProjectRoot("views", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		handlers.Details(w, r)
	})

	mux.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.SearchArtists(w, r)
	})
}
