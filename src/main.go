package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	config globalConfig
)

func main() {
	// Prepare the flags
	addrPtr := flag.String("address", ":7096", "IP address to bind the HTTP server to")
	configPathPtr := flag.String("config", "./config.yaml", "path to the configuration")
	flag.Parse()

	// Load the configuration
	configData, configError := loadConfiguration(*configPathPtr)
	config = configData
	if configError != nil {
		log.Fatalln("Failed to parse configuration", configError)
	}

	// Setup the router
	router := createRouter()
	err := http.ListenAndServe(*addrPtr, router)
	if err != nil {
		log.Fatalln("Fatal Error has occured", err)
	} else {
		log.Println("Closing gracefully")
	}
}

//createRouter initializes the routes
func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Setup the routes
	router.HandleFunc("/", routeBase)
	router.HandleFunc("/{project}/deploy/", routeDeploy).
		Methods("POST")
	router.HandleFunc("/{project}/deploy/{provider}/", routeProvider).
		Methods("POST")

	return router
}

func routeBase(w http.ResponseWriter, r *http.Request) {
	// Return Content
}

func routeProvider(w http.ResponseWriter, r *http.Request) {
	// 1. Get provider
	// 2. Validate provider auth
}

func routeDeploy(w http.ResponseWriter, r *http.Request) {
	// 1. Validate secret
}
