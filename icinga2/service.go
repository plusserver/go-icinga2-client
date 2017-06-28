package icinga2

import (
	"fmt"
	"net/url"
)

type Service struct {
	Name         string `json:"display_name"`
	HostName     string `json:"host_name"`
	CheckCommand string `json:"check_command"`
	Notes        string `json:"notes"`
	NotesURL     string `json:"notes_url"`
	Vars         Vars   `json:"vars"`
}

type ServiceResults struct {
	Results []struct {
		Service Service `json:"attrs"`
	} `json:"results"`
}

type ServiceCreate struct {
	Attrs Service `json:"attrs"`
}

func (s *Service) fullName() string {
	return s.HostName + "!" + s.Name
}

func (s *WebClient) GetService(name string) (Service, error) {
	var serviceResults ServiceResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/services/"+name, nil, &serviceResults, nil)
	if err != nil {
		return Service{}, err
	}
	if resp.HttpResponse().StatusCode != 200 {
		return Service{}, fmt.Errorf("Did not get 200 OK")
	}
	return serviceResults.Results[0].Service, nil
}

func (s *WebClient) CreateService(service Service) error {
	serviceCreate := ServiceCreate{Attrs: service}
	err := s.CreateObject("/services/"+service.fullName(), serviceCreate)
	return err
}

func (s *WebClient) ListServices() (services []Service, err error) {
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

func (s *WebClient) DeleteService(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/services/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return
}

func (s *WebClient) UpdateService(service Service) error {
	serviceUpdate := ServiceCreate{Attrs: service}

	err := s.UpdateObject("/services/"+service.fullName(), serviceUpdate)
	return err
}

func (s *MockClient) GetService(name string) (Service, error) {
	return s.Services[name], nil
}

func (s *MockClient) CreateService(service Service) error {
	s.Services[service.fullName()] = service
	return nil
}

func (s *MockClient) ListServices() ([]Service, error) {
	services := []Service{}

	for _, x := range s.Services {
		services = append(services, x)
	}

	return services, nil
}

func (s *MockClient) DeleteService(name string) error {
	delete(s.Services, name)
	return nil
}

func (s *MockClient) UpdateService(service Service) error {
	s.Services[service.fullName()] = service
	return nil
}
