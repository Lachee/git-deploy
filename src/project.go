package main

//TODO:
// * use this file for SSH stuff https://github.com/melbahja/goph

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type project struct {
	config projectConfig
}

func newProject(config projectConfig) *project {
	project := &project{
		config: config,
	}
	return project
}

//hasSSH determines if the project should connect via ssh first
func (p *project) hasSSH() bool {
	return p.config.SSH.Host != "" && p.config.SSH.User != ""
}

//readLocalConfig reads the configuration from the project
func (p *project) readLocalConfig(system system) (localProjectConfig, error) {
	// Prepare config
	config := localProjectConfig{}

	// Read the data
	data, err := system.read(p.config.ConfigPath)
	if err != nil {
		return config, err
	}

	// Parse the yaml
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// Return the defaults
	return config.check()
}

//updateRepository performs a series of actions to update the codebase
func (p *project) updateRepository(git git) (bool, error) {

	var err error
	var hasChange bool

	// Check if there is a change
	hasChange, err = git.hasChanges()
	if err != nil {
		log.Println("failed to check status change", err)
		return false, err
	}

	// If there was a change, we need to handle stashes.
	//  otherwise, we will just do a pull
	if hasChange {
		err = git.pushStash()
		if err != nil {
			log.Println("failed to stash", err)
			return false, err
		}

		hasChange, err = git.pull()
		if err != nil {
			log.Println("failed to pull", err)
			return false, err
		}

		err = git.popStash()
		if err != nil {
			log.Println("failed to pop", err)
			return false, err
		}
	} else {
		hasChange, err = git.pull()
		if err != nil {
			log.Println("failed to pull", err)
			return false, err
		}
	}

	// Return if we have changed.
	return hasChange, nil
}

//deploy the project's latest changes
func (p *project) deploy() error {

	// Revert the CWD back afterwards (incase there was any changes)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	var localConfig localProjectConfig
	var deployConfig deployConfig
	var err error
	var sys system
	var git git

	// TODO: Establish an SSH connection
	sys = newLocalSystem(p.config.ProjectDirectory)
	git = newGit(sys)

	// Read the configuration
	localConfig, err = p.readLocalConfig(sys)
	if err != nil {
		return err
	}

	// Get active branch
	branch, err := git.currentBranch()
	if err != nil {
		log.Println("failed to load git branch", err)
		return err
	}

	// See if we can find a deploy for thsi branch
	foundConfig := false
	for _, dconfig := range localConfig.Deploys {
		if containsStr(dconfig.Branches, branch) {
			deployConfig = dconfig
			foundConfig = true
			break
		}
	}
	if !foundConfig {
		return fmt.Errorf("failed to find an appropriate deploy script")
	}

	// Update the repository, if there was changes reload the configuration
	gitChanged, gitError := p.updateRepository(git)
	if gitError != nil {
		return gitError
	}
	if gitChanged {
		// TODO: Reload the configuration
		log.Println("Git Changed")
	}

	//Establish a deployer. If none was found assume it was a script
	deployer := createDeployment(deployConfig.Use, deployConfig.With)
	if deployer == nil {
		deployConfig.With["script"] = deployConfig.Use
		deployer = createDeployment("script", deployConfig.With)
	}

	//Run the deployer with the appropriate enviroment variables
	env := p.enviromentVariables(deployConfig.EnviromentVariables)
	sys.setEviromentVariables(env)
	deployer.deploy(sys)

	// TODO: Return the result in some kind of history
	return nil
}

//enviromentVariables creates a list of enviroment variable definitions, with the project's configured enviros taking presidence
func (p *project) enviromentVariables(additionalVariables map[string]string) []string {
	enviros := []string{}
	for k, v := range additionalVariables {
		if k != "" {
			enviros = append(enviros, k+"="+v)
		}
	}
	for k, v := range p.config.EnviromentVariables {
		if k != "" {
			enviros = append(enviros, k+"="+v)
		}
	}
	return enviros
}
