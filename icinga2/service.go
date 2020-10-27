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
	Zone         string `json:"zone,omitempty"`
}

type ServiceResults struct {
	Results []struct {
		Service Service `json:"attrs"`
	} `json:"results"`
}

type ServiceCreate struct {
	Attrs Service `json:"attrs"`
}

func (s Service) GetCheckCommand() string {
	return s.CheckCommand
}

func (s Service) GetVars() Vars {
	return s.Vars
}

func (s Service) GetNotes() string {
	return s.Notes
}

func (s Service) GetNotesURL() string {
	return s.NotesURL
}

func (s *Service) FullName() string {
	return s.HostName + "!" + s.Name
}

func (s *WebClient) GetService(name string) (Service, error) {
	var serviceResults ServiceResults
	resp, err := s.napping.Get(s.URL+"/v1/objects/services/"+name, nil, &serviceResults, nil)
	if err != nil {
		return Service{}, err
	}

	switch resp.HttpResponse().StatusCode {
	case 404:
		return Service{}, ErrNotFound
	case 403:
		return Service{}, ErrForbidden
	case 401:
		return Service{}, ErrUnauthorized
	case 200:
		return serviceResults.Results[0].Service, nil
	}
	return Service{}, fmt.Errorf("Got http error: %d: %w", resp.HttpResponse().StatusCodem, ErrUnknown)
}

func (s *WebClient) CreateService(service Service) error {
	serviceCreate := ServiceCreate{Attrs: service}
	err := s.CreateObject("/services/"+service.FullName(), serviceCreate)
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
		if s.Zone == "" || s.Zone == result.Service.Zone {
			services = append(services, result.Service)
		}
	}

	return
}

func (s *WebClient) DeleteService(name string) (err error) {
	_, err = s.napping.Delete(s.URL+"/v1/objects/services/"+name, &url.Values{"cascade": []string{"1"}}, nil, nil)
	return
}

func (s *WebClient) UpdateService(service Service) error {
	serviceUpdate := ServiceCreate{Attrs: service}

	err := s.UpdateObject("/services/"+service.FullName(), serviceUpdate)
	return err
}

func (s *MockClient) GetService(name string) (Service, error) {
	if sv, ok := s.Services[name]; ok {
		return sv, nil
	} else {
		return Service{}, fmt.Errorf("service not found")
	}
}

func (s *MockClient) CreateService(service Service) error {
	s.mutex.Lock()
	s.Services[service.FullName()] = service
	s.mutex.Unlock()
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
	s.mutex.Lock()
	delete(s.Services, name)
	s.mutex.Unlock()
	return nil
}

func (s *MockClient) UpdateService(service Service) error {
	s.mutex.Lock()
	s.Services[service.FullName()] = service
	s.mutex.Unlock()
	return nil
}
