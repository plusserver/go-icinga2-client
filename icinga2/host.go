package icinga2

import (
	"net/url"
)

type Host struct {
	Name         string                 `json:"__name"`
	Address      string                 `json:"address,omitempty"`
	Address6     string                 `json:"address6,omitempty"`
	DisplayName  string                 `json:"display_name"`
	CheckCommand string                 `json:"check_command,omitempty"`
	Vars         map[string]interface{} `json:"vars"`
	Groups       []string               `json:"groups"`
}

type HostResults struct {
	Results []struct {
		Host Host `json:"attrs"`
	} `json:"results"`
}

type HostCreate struct {
	Templates []string `json:"templates"`
	Attrs     Host     `json:"attrs"`
}

func (s *Server) GetHost(name string) (Host, error) {
	var hostResults HostResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/hosts/"+name, nil, &hostResults, nil)
	if err != nil {
		return Host{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		panic("Did not get 200 OK")
	}
	return hostResults.Results[0].Host, nil
}

func (s *Server) CreateHost(host Host) error {
	hostCreate := HostCreate{Templates: []string{"generic-host"}, Attrs: host}
	var result Results
	_, err := s.napping.Put(s.URL+"/v1/objects/hosts/"+host.Name, hostCreate, &result, nil)
	if err != nil {
		panic(err)
	}

	return err
}

func (s *Server) ListHosts() (hosts []Host, err error) {
	var hostResults HostResults
	hosts = []Host{}

	_, err = s.napping.Get(s.URL+"/v1/objects/hosts/", nil, &hostResults, nil)
	if err != nil {
		return
	}
	for _, result := range hostResults.Results {
		hosts = append(hosts, result.Host)
	}

	return
}

func (s *Server) DeleteHost(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/hosts/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return
}
