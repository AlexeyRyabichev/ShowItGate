package main

import (
	"log"
	"net/http"

	"github.com/AlexeyRyabichev/ShowItGate/internal"
)

var cfgFile = "cfg.json"

func main() {
	routerCfg, err := internal.ReadCfgFromJSON(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	router := internal.NewRouter(routerCfg)

	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(":7050", router.Router))
}
