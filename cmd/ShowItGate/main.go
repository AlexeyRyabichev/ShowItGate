package main

import (
	"log"
	"net/http"

	"github.com/AlexeyRyabichev/ShowItGate/internal"
)

func main() {
	routerCfg := internal.RouterCfg{Name: "ShowItGate"}

	router := internal.NewRouter(routerCfg)

	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(":7050", router.Router))
}
