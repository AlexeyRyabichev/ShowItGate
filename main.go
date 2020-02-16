package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/AlexeyRyabichev/ShowItGate/internal"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var router *mux.Router
var routes = Routes{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Index,
	},
}

func main() {
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

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
