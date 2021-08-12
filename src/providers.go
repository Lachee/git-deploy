package main

import (
	"crypto/subtle"
	"log"
	"net/http"
)

type provider interface {
	verify(secret string, w http.ResponseWriter, r *http.Request) bool
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

func (p *webProvider) verify(secret string, w http.ResponseWriter, r *http.Request) bool {
	key := r.Header.Get("X-API-Key")
	if key == "" {
		w.Header().Add("X-Reason", "No X-API-KEY supplied")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No API Key supplied"))
		return false
	}

	// TODO: Pad the secret and the key
	// ConstantTimeCompare exits early on mismatch length
	return subtle.ConstantTimeCompare([]byte(key), []byte(secret)) == 1
}
