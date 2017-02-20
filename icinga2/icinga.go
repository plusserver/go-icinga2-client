package icinga2

import (
	"crypto/tls"
	"gopkg.in/jmcvetta/napping.v3"
	"net/http"
	"net/url"
)

type Server struct {
	napping     napping.Session
	URL         string
	Username    string
	Password    string
	Debug       bool
	InsecureTLS bool
}

func New(s Server) (*Server, error) {
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

type Results struct {
	Results []struct {
		Code   float64  `json:"code"`
		Errors []string `json:"errors"`
	} `json:"results"`
}
