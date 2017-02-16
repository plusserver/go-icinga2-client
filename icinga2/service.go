package icinga2

import (
	"net/url"
)

type Service struct {
	Name         string                 `json:"__name"`
	HostName     string                 `json:"host_name"`
	CheckCommand string                 `json:"check_command"`
	Vars         map[string]interface{} `json:"vars"`
}

type ServiceResults struct {
	Results []struct {
		Service Service `json:"attrs"`
	} `json:"results"`
}

type ServiceCreate struct {
	Templates []string `json:"templates"`
	Attrs     Service  `json:"attrs"`
}

func (s *Server) GetService(name string) (Service, error) {
	var serviceResults ServiceResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/services/"+name, nil, &serviceResults, nil)
	if err != nil {
		return Service{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		panic("Did not get 200 OK")
	}
	return serviceResults.Results[0].Service, nil
}

func (s *Server) CreateService(service Service) error {
	serviceCreate := ServiceCreate{Templates: []string{}, Attrs: service}
	var result Results
	_, err := s.napping.Put(s.URL+"/v1/objects/services/"+service.HostName+"!"+service.Name, serviceCreate, &result, nil)
	if err != nil {
		panic(err)
	}

	return err
}

func (s *Server) ListServices() (services []Service, err error) {
	var serviceResults ServiceResults
	services = []Service{}

	_, err = s.napping.Get(s.URL+"/v1/objects/services/", nil, &serviceResults, nil)
	if err != nil {
		return
	}
	for _, result := range serviceResults.Results {
		services = append(services, result.Service)
	}

	return
}

func (s *Server) DeleteService(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/services/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return
}
