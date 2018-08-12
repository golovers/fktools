package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Route represent rest api routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes list of route
type Routes []Route

var routes = Routes{
	Route{"Index", "GET", "/", index},
	Route{"Group Status", "GET", "/group", groupStatus},
	Route{"Sprint Status", "GET", "/sprint", sprintStatus},
	Route{"Links Configurations", "GET", "/links", linkConfig},
	Route{"Add link", "POST", "/add_links", addLink},
	Route{"Delete link", "DELETE", "/links", deleteLinks},
}

// NewRouter return a new router with middlewares registered
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = PrometheusMiddleware(handler)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	router.Handle("/api/v1/metrics", promhttp.Handler())
	return router
}
