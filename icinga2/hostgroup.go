package icinga2

import (
	"fmt"
)

type HostGroup struct {
	Name string `json:"display_name,omitempty"`
	Vars Vars   `json:"vars"`
	Zone string `json:"zone,omitempty"`
}

type HostGroupResults struct {
	Results []struct {
		HostGroup HostGroup `json:"attrs"`
	} `json:"results"`
}

type HostGroupCreate struct {
	Templates []string  `json:"templates"`
	Attrs     HostGroup `json:"attrs"`
}

func (s *WebClient) GetHostGroup(name string) (HostGroup, error) {
	var hostGroupResults HostGroupResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/hostgroups/"+name, nil, &hostGroupResults, nil)
	if err != nil {
		return HostGroup{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		return HostGroup{}, fmt.Errorf("Did not get 200 OK")
	}
	return hostGroupResults.Results[0].HostGroup, nil
}

func (s *WebClient) CreateHostGroup(hostGroup HostGroup) error {
	hostGroupCreate := HostGroupCreate{Attrs: hostGroup}
	err := s.CreateObject("/hostgroups/"+hostGroup.Name, hostGroupCreate)
	return err
}

func (s *WebClient) ListHostGroups() (hostGroups []HostGroup, err error) {
	var hostGroupResults HostGroupResults
	hostGroups = []HostGroup{}

	_, err = s.napping.Get(s.URL+"/v1/objects/hostgroups/", nil, &hostGroupResults, nil)
	if err != nil {
		return
	}
	for _, result := range hostGroupResults.Results {
		if s.Zone == "" || s.Zone == result.HostGroup.Zone {
			hostGroups = append(hostGroups, result.HostGroup)
		}
	}

	return
}

func (s *WebClient) DeleteHostGroup(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/hostgroups/"+name, nil, nil, nil)
	return
}

func (s *WebClient) UpdateHostGroup(hostGroup HostGroup) error {
	hostGroupUpdate := HostGroupCreate{Attrs: hostGroup}

	err := s.UpdateObject("/hostgroups/"+hostGroup.Name, hostGroupUpdate)
	return err
}

func (s *MockClient) GetHostGroup(name string) (HostGroup, error) {
	if hg, ok := s.Hostgroups[name]; ok {
		return hg, nil
	} else {
		return HostGroup{}, fmt.Errorf("hostgroup not found")
	}
}

func (s *MockClient) CreateHostGroup(hostGroup HostGroup) error {
	s.mutex.Lock()
	s.Hostgroups[hostGroup.Name] = hostGroup
	s.mutex.Unlock()
	return nil
}

func (s *MockClient) ListHostGroups() ([]HostGroup, error) {
	hostGroups := []HostGroup{}

	for _, x := range s.Hostgroups {
		hostGroups = append(hostGroups, x)
	}

	return hostGroups, nil
}

func (s *MockClient) DeleteHostGroup(name string) error {
	s.mutex.Lock()
	delete(s.Hostgroups, name)
	s.mutex.Unlock()
	return nil
}

func (s *MockClient) UpdateHostGroup(hostGroup HostGroup) error {
	s.mutex.Lock()
	s.Hostgroups[hostGroup.Name] = hostGroup
	s.mutex.Unlock()
	return nil
}
