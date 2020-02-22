package internal

import "github.com/AlexeyRyabichev/ShowItGate/public"

type Node struct {
	Gateways public.Gateways `json:"gateways"`
	Name     string          `json:"name"`
	Base     string          `json:"base"`
	Host     string          `json:"host"`
	Scheme   string          `json:"scheme"`
}