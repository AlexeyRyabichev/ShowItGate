package ShowItGate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type NodeCfg struct {
	Gateways Gateways `json:"gateways"`
	Name     string   `json:"name"`
	Base     string   `json:"base"`
	Host     string   `json:"host"`
	Scheme   string   `json:"scheme"`
	ApiKey   string   `json:"-"`
	Token    string   `json:"-"`
}

type internalNodeCfg struct {
	Gateways Gateways `json:"gateways"`
	Name     string   `json:"name"`
	Base     string   `json:"base"`
	Host     string   `json:"host"`
	Scheme   string   `json:"scheme"`
	ApiKey   string   `json:"api_key"`
	Token    string   `json:"token"`
}

func (cfg *NodeCfg) RegisterNode() error {
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

		req, err := http.NewRequest("POST", "http://64.225.109.162:7050/node", bytes.NewBuffer(cfgBytes))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	cfg.Token = resp.Header.Get("X-Token")
	return nil
}

func ReadCfgFromJSON(jsonFile string) (NodeCfg, error) {
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return NodeCfg{}, err
	}

	cfg := internalNodeCfg{}
	if err := json.Unmarshal(file, &cfg); err != nil {
		return NodeCfg{}, err
	}

	return migrate(&cfg), nil
}

func (cfg *NodeCfg) SaveCfgToJSON(jsonFile string) error {
	iCfg := reverseMigrate(cfg)

	file, err := json.MarshalIndent(iCfg, "", "  ")
	if err != nil{
		return err
	}

	if err := ioutil.WriteFile(jsonFile, file, 0666); err != nil{
		return err
	}

	return nil
}

func migrate(cfg *internalNodeCfg) NodeCfg {
	return NodeCfg{
		Gateways: cfg.Gateways,
		Name:     cfg.Name,
		Base:     cfg.Base,
		Host:     cfg.Host,
		Scheme:   cfg.Scheme,
		ApiKey:   cfg.ApiKey,
		Token:    cfg.Token,
	}
}

func reverseMigrate(cfg *NodeCfg) internalNodeCfg {
	return internalNodeCfg{
		Gateways: cfg.Gateways,
		Name:     cfg.Name,
		Base:     cfg.Base,
		Host:     cfg.Host,
		Scheme:   cfg.Scheme,
		ApiKey:   cfg.ApiKey,
		Token:    cfg.Token,
	}
}
