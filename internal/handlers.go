package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

func (rt *Router) NodePost(w http.ResponseWriter, r *http.Request) {
	var node Node

	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := rt.Nodes[node.Base]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	for _, gateway := range node.Gateways {
		rt.addRoute(Route{
			Name:        gateway.Name,
			Method:      gateway.Method,
			Pattern:     path.Join(node.Base, gateway.Path),
			HandlerFunc: rt.proxyFunc,
		})
	}
	rt.Nodes[node.Base] = node

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *Router) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (rt *Router) proxyFunc(w http.ResponseWriter, r *http.Request) {
	pathElements := strings.Split(r.URL.Path, "/")
	base := fmt.Sprintf("/%s/%s", pathElements[1], pathElements[2])

	if _, ok := rt.Nodes[base]; !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	node := rt.Nodes[base]

	newURL := r.URL
	newURL.Host = node.Host
	newURL.Scheme = node.Scheme

	req, err := http.NewRequest(r.Method, newURL.String(), r.Body)
	req.Header = r.Header

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
