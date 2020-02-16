package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"

	"github.com/AlexeyRyabichev/ShowItGate/internal"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Gateway struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

type Gateways []Gateway

type Node struct {
	Gateways Gateways `json:"gateways"`
	Name     string   `json:"name"`
	Base     string   `json:"base"`
	Host     string   `json:"host"`
	Scheme   string   `json:"scheme"`
}

var router *mux.Router
var nodes map[string]Node

var routes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Index,
	},

	Route{
		Name:        "Register node",
		Method:      "POST",
		Pattern:     "/node",
		HandlerFunc: NodePost,
	},
}

func main() {
	nodes = make(map[string]Node)
	initRouter()

	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(":7050", router))
}

func initRouter() {
	router = mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		addRoute(route)
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"HANDLER NOT FOUND FOR REQUEST: %s %s",
			r.Method,
			r.RequestURI,
		)
	})
}

func addRoute(route Route) {
	var handler http.Handler
	handler = route.HandlerFunc
	handler = internal.Logger(handler, route.Name)

	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

func NodePost(w http.ResponseWriter, r *http.Request) {
	var node Node

	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := nodes[node.Base]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	for _, gateway := range node.Gateways {
		addRoute(Route{
			Name:        gateway.Name,
			Method:      gateway.Method,
			Pattern:     path.Join(node.Base, gateway.Path),
			HandlerFunc: nil,
		})
	}
	nodes[node.Base] = node

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
