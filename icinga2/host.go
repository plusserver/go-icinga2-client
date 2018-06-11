package icinga2

import (
	"fmt"
	"net/url"
)

type Host struct {
	Name         string   `json:"display_name"`
	Address      string   `json:"address,omitempty"`
	Address6     string   `json:"address6,omitempty"`
	CheckCommand string   `json:"check_command,omitempty"`
	Notes        string   `json:"notes"`
	NotesURL     string   `json:"notes_url"`
	Vars         Vars     `json:"vars"`
	Groups       []string `json:"groups,omitempty"`
	Zone         string   `json:"zone,omitempty"`
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

func (s *WebClient) GetHost(name string) (Host, error) {
	var hostResults HostResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/hosts/"+name, nil, &hostResults, nil)
	if err != nil {
		return Host{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		return Host{}, fmt.Errorf("Did not get 200 OK")
	}
	return hostResults.Results[0].Host, nil
}

func (s *WebClient) CreateHost(host Host) error {
	hostCreate := HostCreate{Templates: []string{"generic-host"}, Attrs: host}
	err := s.CreateObject("/hosts/"+host.Name, hostCreate)
	return err
}

func (s *WebClient) ListHosts() (hosts []Host, err error) {
	var hostResults HostResults
	hosts = []Host{}

	_, err = s.napping.Get(s.URL+"/v1/objects/hosts/", nil, &hostResults, nil)
	if err != nil {
		return
	}
	for _, result := range hostResults.Results {
		if s.Zone == "" || s.Zone == result.Host.Zone {
			hosts = append(hosts, result.Host)
		}
	}

	return
}

func (s *WebClient) DeleteHost(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/hosts/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return
}

func (s *WebClient) UpdateHost(host Host) error {
	host.Groups = []string{} // must be empty when updating the Host
	hostUpdate := HostCreate{Attrs: host}
	err := s.UpdateObject("/hosts/"+host.Name, hostUpdate)
	return err
}

func (s *MockClient) GetHost(name string) (Host, error) {
	if host, ok := s.Hosts[name]; ok {
		return host, nil
	} else {
		return Host{}, fmt.Errorf("host not found")
	}
}

func (s *MockClient) CreateHost(host Host) error {
	s.mutex.Lock()
	s.Hosts[host.Name] = host
	s.mutex.Unlock()
	return nil
}

func (s *MockClient) ListHosts() ([]Host, error) {
	hosts := []Host{}

	for _, x := range s.Hosts {
		hosts = append(hosts, x)
	}

	return hosts, nil
}

func (s *MockClient) DeleteHost(name string) error {
	s.mutex.Lock()
	delete(s.Hosts, name)
	s.mutex.Unlock()
	return nil
}

func (s *MockClient) UpdateHost(host Host) error {
	s.mutex.Lock()
	s.Hosts[host.Name] = host
	s.mutex.Unlock()
	return nil
}
