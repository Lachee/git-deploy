package main

//TODO:
// * use this file for SSH stuff https://github.com/melbahja/goph

import (
	"errors"
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

//readBranch gets the branch of the current project
func (p *project) getBranch(sys system) string {
	branch, _ := gitCurrentBranch(sys)
	return branch
}

//updateRepository performs a series of actions to update the codebase
func (p *project) updateRepository() (bool, error) {
	if p.config.Update != "" {
		// TODO: Execute p.config.Update
		cmd := shellCommand("", p.config.Update)
		return true, cmd.Run()
	}

	/*
		changed := false
		status := git.status()
		if status.hasChanges {
			git.stash()
			changed = git.pull()
			git.stash("pop")
		} else {
			changed = git.pull()
		}
		return changed, nil
	*/

	return false, errors.New("git clone: not yet implemented")
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

	// TODO: Establish an SSH connection
	sys = newLocalSystem(p.config.ProjectDirectory)

	// Read the configuration
	localConfig, err = p.readLocalConfig(sys)
	if err != nil {
		return err
	}

	// Check if the configurations apply to this branch
	foundConfig := false
	branch := p.getBranch(sys)
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
	gitChanged, gitError := p.updateRepository()
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
