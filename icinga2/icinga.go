package icinga2

import (
	"crypto/tls"
	"gopkg.in/jmcvetta/napping.v3"
	"net/http"
	"net/url"
)

type Server struct {
	napping  napping.Session
	URL      string
	Username string
	Password string
}

func New(incigaUrl, username, password string) (*Server, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	sess := napping.Session{
		Log:      true,
		Client:   client,
		Userinfo: url.UserPassword(username, password),
	}

	server := Server{URL: incigaUrl, Username: username, Password: password, napping: sess}
	return &server, nil
}

type Results struct {
	Results []struct {
		Code   float64  `json:"code"`
		Errors []string `json:"errors"`
	} `json:"results"`
}
