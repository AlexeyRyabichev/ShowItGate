package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/AlexeyRyabichev/ShowItGate"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type RouterCfg struct {
	Name    string
	Nodes   map[string]ShowItGate.NodeCfg
	ApiKeys map[string]bool
}

type Router struct {
	cfg    RouterCfg
	routes []Route

	//Nodes  map[string]ShowItGate.NodeCfg
	Router *mux.Router
}

func NewRouter(cfg RouterCfg) *Router {
	router := Router{
		cfg: cfg,
		//Nodes: make(map[string]ShowItGate.NodeCfg),
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

	for _, node := range rt.cfg.Nodes{
		for _, gateway := range node.Gateways {
			rt.addRoute(Route{
				Name:        gateway.Name,
				Method:      gateway.Method,
				Pattern:     path.Join(node.Base, gateway.Path),
				HandlerFunc: rt.proxyFunc,
			})
		}
	}

	rt.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"HANDLER NOT FOUND FOR REQUEST: %s %s",
			r.Method,
			r.RequestURI,
		)
		w.WriteHeader(http.StatusNotFound)
	})

	rt.Router.Use(mux.CORSMethodMiddleware(rt.Router))
}

func (rt *Router) addRoute(route Route) {
	var handler http.Handler
	handler = route.HandlerFunc
	handler = ShowItGate.Logger(handler, route.Name)

	rt.Router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)

	rt.Router.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}).Methods(http.MethodOptions)

	log.Printf("registered route %s", route.Pattern)
}

func ReadCfgFromJSON(jsonFile string) (RouterCfg, error) {
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return newRouterCfg(), err
	}

	cfg := RouterCfg{}
	if err := json.Unmarshal(file, &cfg); err != nil {
		return newRouterCfg(), err
	}

	return cfg, nil
}

func (cfg *RouterCfg) SaveCfgToJSON(jsonFile string) error {
	file, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(jsonFile, file, 0666); err != nil {
		return err
	}

	return nil
}

func newRouterCfg() RouterCfg{
	return RouterCfg{
		Name:    "ShowItGate",
		Nodes:   make(map[string]ShowItGate.NodeCfg),
		ApiKeys: make(map[string]bool),
	}
}