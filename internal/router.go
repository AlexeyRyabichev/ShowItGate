package internal

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/AlexeyRyabichev/ShowItGate/public"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type RouterCfg struct {
	Name string
}

type Router struct {
	cfg    RouterCfg
	routes []Route

	Nodes  map[string]Node
	Router *mux.Router
}

func NewRouter(cfg RouterCfg) *Router {
	router := Router{
		cfg:   cfg,
		Nodes: make(map[string]Node),
	}
	router.routes = []Route{
		{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: router.Index,
		},
		{
			Name:        "Register node",
			Method:      "POST",
			Pattern:     "/node",
			HandlerFunc: router.NodePost,
		},
	}
	router.initRouter()
	return &router
}

func (rt *Router) initRouter() {
	rt.Router = mux.NewRouter().StrictSlash(true)

	for _, route := range rt.routes {
		rt.addRoute(route)
	}

	rt.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"HANDLER NOT FOUND FOR REQUEST: %s %s",
			r.Method,
			r.RequestURI,
		)
		w.WriteHeader(http.StatusNotFound)
	})
}

func (rt *Router) addRoute(route Route) {
	var handler http.Handler
	handler = route.HandlerFunc
	handler = public.Logger(handler, route.Name)

	rt.Router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}
