package icinga2

type HostGroup struct {
	Name        string                 `json:"__name"`
	DisplayName string                 `json:"display_name"`
	Vars        map[string]interface{} `json:"vars"`
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

func (s *Server) GetHostGroup(name string) (HostGroup, error) {
	var hostGroupResults HostGroupResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/hostgroups/"+name, nil, &hostGroupResults, nil)
	if err != nil {
		return HostGroup{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		panic("Did not get 200 OK")
	}
	return hostGroupResults.Results[0].HostGroup, nil
}

func (s *Server) CreateHostGroup(hostGroup HostGroup) error {
	hostGroupCreate := HostGroupCreate{Attrs: hostGroup}
	var result Results
	_, err := s.napping.Put(s.URL+"/v1/objects/hostgroups/"+hostGroup.Name, hostGroupCreate, &result, nil)
	if err != nil {
		panic(err)
	}

	return err
}

func (s *Server) ListHostGroups() (hostGroups []HostGroup, err error) {
	var hostGroupResults HostGroupResults
	hostGroups = []HostGroup{}

	_, err = s.napping.Get(s.URL+"/v1/objects/hostgroups/", nil, &hostGroupResults, nil)
	if err != nil {
		return
	}
	for _, result := range hostGroupResults.Results {
		hostGroups = append(hostGroups, result.HostGroup)
	}

	return
}

func (s *Server) DeleteHostGroup(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/hostgroups/"+name, nil, nil, nil)
	return
}
