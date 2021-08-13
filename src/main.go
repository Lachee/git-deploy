package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
	"github.com/lachee/git-deploy/src/web"
)

var (
	configPath string
	config     globalConfig
	processing map[string]time.Time = make(map[string]time.Time)
)

//loadProject loads the configuration and finds the appropriate project.
func loadProject(name string) (*project, error) {
	// Load the configuration
	// configData, configError := loadConfiguration(configPath)
	// config = configData
	// if configError != nil {
	// 	log.Fatalln("Failed to parse configuration", configError)
	// }

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
	flag.Parse()

	// Set config path
	configPath = *configPathPtr

	// Load the configuration
	configData, configError := loadConfiguration(configPath)
	config = configData
	if configError != nil {
		log.Fatalln("Failed to parse configuration", configError)
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
func createRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	// Setup the routes
	router.HandleFunc("/", routeBase)
	router.HandleFunc("/{project}/deploy/{provider}", func(w http.ResponseWriter, r *http.Request) {
		response := handleAPI(w, r)
		response.SetContentTypeFromRequest(r)
		response.Write(w)
	}).Methods("POST")

	lmt := tollbooth.NewLimiter(1.0/30.0, nil)
	lmt.
		SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetMethods([]string{"POST", "PUT"})

	return tollbooth.LimitHandler(lmt, router)
}

//routeBase handles any request that isn't matching
func routeBase(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/Lachee/git-deploy", http.StatusTemporaryRedirect)
}

func handleAPI(w http.ResponseWriter, r *http.Request) web.Response {
	vars := mux.Vars(r)

	// Ensure we have the project
	project, err := loadProject(vars["project"])
	if err != nil {
		return web.NewErrorResponse(http.StatusNotFound, err)
	}

	// Ensure we have the provider
	provider := createProvider(vars["provider"])
	if provider == nil {
		log.Println("cannot find the provider:", vars["provider"])
		return web.NewErrorResponse(http.StatusNotFound, errors.New("provider cannot be found"))
	}

	// Ensure the provider is correct
	verified := provider.verify([]byte(project.config.Secret), w, r)
	if !verified {
		return web.NewErrorResponse(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Ensure we are not already deploying
	_, alreadyDeploying := processing[project.config.Name]
	if alreadyDeploying {
		return web.NewErrorResponse(http.StatusConflict, errors.New("already deploying"))
	}

	// Deploy if we can on a new go-routine
	processing[project.config.Name] = time.Now()
	go func() {
		project.deploy()
		delete(processing, project.config.Name)
	}()

	// Return the status
	return web.NewResponse(http.StatusOK, "deployed")
}
