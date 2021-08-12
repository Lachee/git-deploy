package main

import (
	"log"
	"net/http"
)

type provider interface {
	verify(project *project, w http.ResponseWriter, r *http.Request) bool
}

func createProvider(name string) provider {
	var p provider = nil
	switch name {
	case "":
	case "web":
		p = &webProvider{}
	}

	if p == nil {
		log.Printf("error: there is no deployer for %s\n", name)
		return nil
	}

	return p
}

type webProvider struct {
}

func (p *webProvider) verify(project *project, w http.ResponseWriter, r *http.Request) bool {
	return true
}
