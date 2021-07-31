package main

import (
	"fmt"
	"io/ioutil"

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
func (p *project) readLocalConfig() (localProjectConfig, error) {
	config := localProjectConfig{}

	// Read the data
	data, err := ioutil.ReadFile(p.config.ConfigPath)
	if err != nil {
		return config, err
	}

	// Parse the yaml
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config.check()
}

//readBranch gets the branch of the current project
func (p *project) getBranch() string {
	//TODO: Setup branch detection
	return "master"
}

func (p *project) deploy() error {

	// TODO: Establish an SSH connection
	var localConfig localProjectConfig
	var deployConfig deployConfig
	var err error

	// Read the configuration
	localConfig, err = p.readLocalConfig()
	if err != nil {
		return err
	}

	// Check if the configurations apply to this branch
	foundConfig := false
	branch := p.getBranch()
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

	//Establish a deployer. If none was found assume it was a script
	deployer := createDeployer(deployConfig.Use)
	if deployer == nil {
		deployer = &scriptDeploy{
			script: deployConfig.Use,
		}
	}

	//TODO: Update git repo
	//TODO: Setup enviroment variables

	//Run the deployer
	deployer.deploy(deployConfig.With)

	// Prepare the correct loader type
	//deployer := createDeployer(config.)

	return nil
}
