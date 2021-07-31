package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	config   globalConfig
	projects map[string]*project
)

func main() {
	// Prepare the flags
	addrPtr := flag.String("address", "localhost:7096", "IP address to bind the HTTP server to")
	configPathPtr := flag.String("config", "./config.yaml", "path to the configuration")
	deployPtr := flag.String("deploy", "", "project to immediately deploy and then abort")
	flag.Parse()

	// Load the configuration
	configData, configError := loadConfiguration(*configPathPtr)
	config = configData
	if configError != nil {
		log.Fatalln("Failed to parse configuration", configError)
	}

	// Setup the projects
	projects = make(map[string]*project, len(config.Projects))
	for _, pconfig := range config.Projects {
		projects[pconfig.Name] = newProject(pconfig)
	}

	// If we are early deploying, then do so
	if *deployPtr != "" {
		log.Printf("Deploying %s/\n", *deployPtr)
		projects[*deployPtr].deploy()
		return
	}

	// Setup the router
	router := createRouter()
	log.Printf("Listening to http://%s/\n", *addrPtr)
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
	//router.HandleFunc("/", routeBase)
	router.HandleFunc("/{project}/deploy/", routeDeploy).
		Methods("POST")
	router.HandleFunc("/{project}/deploy/{provider}", routeProvider).
		Methods("POST")

	return router
}

//routeBase handles any request that isn't matching
func routeBase(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/Lachee/git-deploy", http.StatusTemporaryRedirect)
}

func routeProvider(w http.ResponseWriter, r *http.Request) {
	// 1. Get provider
	// 2. Validate provider auth
}

func routeDeploy(w http.ResponseWriter, r *http.Request) {
	// 1. Validate secret
}
