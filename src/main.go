package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	configPath string
	config     globalConfig
)

func test() bool {
	system := newLocalSystem("C:/Users/lachl/go/src/github.com/lachee/git-deploy")
	branch, err := gitCurrentBranch(system)
	if err != nil {
		log.Fatalln("Failed to get the current branch", err)
		return false
	}
	log.Println("Branch", branch)
	return false
}

//loadProject loads the configuration and finds the appropriate project.
func loadProject(name string) (*project, error) {
	// Load the configuration
	configData, configError := loadConfiguration(configPath)
	config = configData
	if configError != nil {
		log.Fatalln("Failed to parse configuration", configError)
	}

	// Find the correct project
	for _, pconfig := range config.Projects {
		if pconfig.Name == name {
			return newProject(pconfig), nil
		}
	}

	return nil, errors.New("failed to find the project")
}

func main() {
	// Prepare the flags
	addrPtr := flag.String("address", "localhost:7096", "IP address to bind the HTTP server to")
	configPathPtr := flag.String("config", "./config.yaml", "path to the configuration")
	deployPtr := flag.String("deploy", "", "project to immediately deploy and then abort")
	testPtr := flag.Bool("test", false, "Runs a test function")
	flag.Parse()

	// Set config path
	configPath = *configPathPtr

	if *testPtr {
		log.Println("Testing Function")
		if test() {
			log.Println("Aborted test")
			return
		}
	}

	// If we are early deploying, then do so
	if *deployPtr != "" {
		log.Printf("Deploying '%s'\n", *deployPtr)
		project, err := loadProject(*deployPtr)
		if err != nil {
			log.Fatalln("cannot find the project:", err)
			return
		}

		deployError := project.deploy()
		if deployError != nil {
			log.Fatalln("Failed to deploy: ", deployError)
		}

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
