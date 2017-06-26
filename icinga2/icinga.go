package icinga2

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/jmcvetta/napping.v3"
)

type Client interface {
	GetHost(string) (Host, error)
	CreateHost(Host) error
	ListHosts() ([]Host, error)
	DeleteHost(string) error
	UpdateHost(Host) error

	GetHostGroup(string) (HostGroup, error)
	CreateHostGroup(HostGroup) error
	ListHostGroups() ([]HostGroup, error)
	DeleteHostGroup(string) error
	UpdateHostGroup(HostGroup) error

	GetService(string) (Service, error)
	CreateService(Service) error
	ListServices() ([]Service, error)
	DeleteService(string) error
	UpdateService(Service) error
}

type WebClient struct {
	napping     napping.Session
	URL         string
	Username    string
	Password    string
	Debug       bool
	InsecureTLS bool
}

type MockClient struct {
	Hostgroups map[string]HostGroup
	Hosts      map[string]Host
	Services   map[string]Service
}

type Vars map[string]interface{}

func New(s WebClient) (*WebClient, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: s.InsecureTLS},
	}
	client := &http.Client{Transport: transport}

	s.napping = napping.Session{
		Log:      s.Debug,
		Client:   client,
		Userinfo: url.UserPassword(s.Username, s.Password),
	}
	return &s, nil
}

func NewMockClient() (c *MockClient) {
	c = new(MockClient)
	c.Hostgroups = make(map[string]HostGroup)
	c.Hosts = make(map[string]Host)
	c.Services = make(map[string]Service)
	return
}

type Results struct {
	Results []struct {
		Code   float64  `json:"code"`
		Errors []string `json:"errors,omitempty"`
		Status string   `json:"status,omitempty"`
		Name   string   `json:"name,omitempty"`
		Type   string   `json:"type,omitempty"`
	} `json:"results"`
}

func (s *WebClient) CreateObject(path string, create interface{}) error {
	var results Results

	resp, err := s.napping.Put(s.URL+"/v1/objects"+path, create, &results, nil)
	if err != nil {
		return err
	}

	return s.handleResults("create", path, resp, &results)
}

func (s *WebClient) UpdateObject(path string, create interface{}) error {
	var results Results

	resp, err := s.napping.Post(s.URL+"/v1/objects"+path, create, &results, nil)
	if err != nil {
		return err
	}

	return s.handleResults("update", path, resp, &results)
}

func (s *WebClient) handleResults(typ, path string, resp *napping.Response, results *Results) error {
	if resp.HttpResponse().StatusCode >= 400 {
		return fmt.Errorf("%s %s : %s %+v\n", typ, path, resp.HttpResponse().Status)
	}

	if len(results.Results) <= 0 {
		return nil
	}

	for _, r := range results.Results {
		if r.Code >= 400.0 {
			return fmt.Errorf("%s %s : %s\n", typ, path, r.Status)
		}
	}

	return nil

}
