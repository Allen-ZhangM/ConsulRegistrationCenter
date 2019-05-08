package models

import (
	. "RegistrationCenter/conf"
	consulapi "github.com/hashicorp/consul/api"
)

type CatalogRegistration struct {
	Node       string                    `json:"node"`
	Address    string                    `json:"address"`
	NodeMeta   map[string]string         `json:"node_meta"`
	Datacenter string                    `json:"datacenter"`
	Service    *AgentServiceRegistration `json:"service"`
}

type AgentServiceRegistration struct {
	ID      string              `json:"id,omitempty"`
	Name    string              `json:"name,omitempty"`
	Tags    []string            `json:"tags,omitempty"`
	Port    int                 `json:"port,omitempty"`
	Address string              `json:"address,omitempty"`
	Meta    map[string]string   `json:"meta,omitempty"`
	Checks  []AgentServiceCheck `json:"checks,omitempty"`
}

type AgentServiceCheck struct {
	CheckID       string              `json:"check_id,omitempty"`
	Name          string              `json:"name,omitempty"`
	Interval      string              `json:"interval,omitempty"`
	Timeout       string              `json:"timeout,omitempty"`
	TTL           string              `json:"ttl,omitempty"`
	HTTP          string              `json:"http,omitempty"`
	Header        map[string][]string `json:"header,omitempty"`
	Method        string              `json:"method,omitempty"`
	TCP           string              `json:"tcp,omitempty"`
	Status        string              `json:"status,omitempty"`
	Notes         string              `json:"notes,omitempty"`
	TLSSkipVerify bool                `json:"tls_skip_verify,omitempty"`
}

// 返回值service信息
type SimpleAgentService struct {
	ID      string
	Name    string
	Tags    []string
	Meta    map[string]string
	Port    int
	Address string
}

func ToSimpleAgentService(as consulapi.CatalogService) SimpleAgentService {
	return SimpleAgentService{
		ID:      as.ServiceID,
		Name:    as.ServiceName,
		Tags:    as.ServiceTags,
		Meta:    as.ServiceMeta,
		Port:    as.ServicePort,
		Address: as.ServiceAddress,
	}
}

func ToAgentServiceRegistration(asr *AgentServiceRegistration) *consulapi.AgentServiceRegistration {
	casc := consulapi.AgentServiceChecks{}
	for _, v := range asr.Checks {
		casc = append(casc, &consulapi.AgentServiceCheck{
			CheckID:       v.CheckID,
			Name:          v.Name,
			Interval:      v.Interval,
			Timeout:       v.Timeout,
			TTL:           v.TTL,
			HTTP:          v.HTTP,
			Header:        v.Header,
			Method:        v.Method,
			TCP:           v.TCP,
			Status:        v.Status,
			Notes:         v.Notes,
			TLSSkipVerify: v.TLSSkipVerify,
		})
	}
	return &consulapi.AgentServiceRegistration{
		ID:      asr.ID,
		Name:    asr.Name,
		Tags:    asr.Tags,
		Port:    asr.Port,
		Address: asr.Address,
		Meta:    asr.Meta,
		Checks:  casc,
	}
}

func DefaultCatalogRegistration() *CatalogRegistration {
	server := RandomServer()
	return &CatalogRegistration{
		Node:       server.Name,
		Address:    server.Address,
		NodeMeta:   nil,
		Datacenter: Datacenters[0],
	}
}
