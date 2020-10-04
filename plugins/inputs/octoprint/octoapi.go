package octoprint

import (
	gooctoprint "github.com/mcuadros/go-octoprint"
)

// OctoAPI an interface for external connections made by go-octoprint
type OctoAPI interface {
	StateRequest() (*gooctoprint.FullStateResponse, error)
	ConnectionRequest() (*gooctoprint.ConnectionResponse, error)
}

// GoOcto implements the OctoAPI interface to expose API calls for the go-octoprint library
type GoOcto struct {
	client *gooctoprint.Client
}

// NewGoOcto returns a GoOcto instances with a valid go-octoprint client
func NewGoOcto(URL string, APIKey string) *GoOcto {
	var goOcto GoOcto
	goOcto.client = gooctoprint.NewClient(URL, APIKey)
	return &goOcto
}

// StateRequest uses the go-octoprint client to query octoprint's API For the current state
func (o *GoOcto) StateRequest() (*gooctoprint.FullStateResponse, error) {
	r := gooctoprint.StateRequest{}
	s, err := r.Do(o.client)
	if err != nil {
		return &gooctoprint.FullStateResponse{}, err
	}

	return s, nil
}

// ConnectionRequest uses the go-octoprint clien to query octoprint's API for the current connection information
func (o *GoOcto) ConnectionRequest() (*gooctoprint.ConnectionResponse, error) {
	r := gooctoprint.ConnectionRequest{}
	c, err := r.Do(o.client)
	if err != nil {
		return &gooctoprint.ConnectionResponse{}, err
	}

	return c, err
}
