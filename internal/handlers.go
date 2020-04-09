package internal

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/AlexeyRyabichev/ShowItGate"
)

func (rt *Router) NodePost(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-Api-Key")

	if _, ok := rt.cfg.ApiKeys[key]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var node ShowItGate.NodeCfg
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := rt.cfg.Nodes[node.Base]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	node.Token = generateRandToken()

	for _, gateway := range node.Gateways {
		rt.addRoute(Route{
			Name:        gateway.Name,
			Method:      gateway.Method,
			Pattern:     path.Join(node.Base, gateway.Path),
			HandlerFunc: rt.proxyFunc,
		})
	}
	rt.cfg.Nodes[node.Base] = node

	if err := rt.cfg.SaveCfgToJSON("cfg.json"); err != nil { //TODO: fix hardcode
		log.Printf("cannot save configuration: %v", err)
	}

	w.Header().Set("X-Token", node.Token)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func generateRandToken() string {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	return fmt.Sprintf("%x", tokenBytes)
}

func (rt *Router) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (rt *Router) proxyFunc(w http.ResponseWriter, r *http.Request) {
	pathElements := strings.Split(r.URL.Path, "/")
	base := fmt.Sprintf("/%s/%s", pathElements[1], pathElements[2])

	if _, ok := rt.cfg.Nodes[base]; !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	node := rt.cfg.Nodes[base]

	newURL := r.URL
	newURL.Host = node.Host
	newURL.Scheme = node.Scheme

	req, err := http.NewRequest(r.Method, newURL.String(), r.Body)
	req.Header = r.Header
	req.Header.Set("X-Token", node.Token)

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
	}

	for name, value := range resp.Header {
		w.Header().Set(name, strings.Join(value, ""))
	}

	w.Write(bodyBytes)
}
