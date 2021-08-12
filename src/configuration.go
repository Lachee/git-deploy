package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//globalConfig holds the configuration for the git-deploy instance
type globalConfig struct {
	Projects []projectConfig `yaml:"projects"`
}

//check setup the defaults for the configuration and throws an error if invalid
func (config globalConfig) check() (globalConfig, error) {
	// Ensure we have proejcts
	if config.Projects == nil || len(config.Projects) == 0 {
		return config, errors.New("failed to load any projects")
	}

	// Check all the children
	for i, element := range config.Projects {
		c, e := element.check()
		if e != nil {
			return config, e
		}
		config.Projects[i] = c
	}

	//Return
	return config, nil
}

//projectConfig holds the configuration for a individual project
type projectConfig struct {
	Name                string            `yaml:"name"`
	ProjectDirectory    string            `yaml:"project"`
	ConfigPath          string            `yaml:"config"`
	Secret              []byte            `yaml:"secret"`
	Providers           []string          `yaml:"providers"`
	Webhook             string            `yaml:"webhook"`
	EnviromentVariables map[string]string `yaml:"env"`
	SSH                 sshConfig         `yaml:"ssh"`
	Update              string            `yaml:"update"`
}

//check setup the defaults for the configuration and throws an error if invalid
func (config projectConfig) check() (projectConfig, error) {
	if config.ConfigPath == "" {
		config.ConfigPath = "./git-deploy.yaml"
	}
	if len(config.Secret) == 0 {
		return config, fmt.Errorf("project %s has a missing secret", config.Name)
	}
	if config.ProjectDirectory == "" {
		return config, fmt.Errorf("project %s does not have a valid project directory", config.Name)
	}
	if config.ConfigPath == "" {
		return config, fmt.Errorf("project %s does not have a valid config path", config.Name)
	}

	// TODO: Pad the secret
	return config, nil
}

//localProjectConfig holds the configuration of the local git instance
type localProjectConfig struct {
	Deploys []deployConfig `yaml:"deploys"`
}

//check setup the defaults for the configuration and throws an error if invalid
func (config localProjectConfig) check() (localProjectConfig, error) {
	// Check all the children
	for i, element := range config.Deploys {
		c, e := element.check()
		if e != nil {
			return config, e
		}
		config.Deploys[i] = c
	}

	return config, nil
}

//deployConfig holds the configuration for a deployment
type deployConfig struct {
	Name                string            `yaml:"name"`
	Branches            []string          `yaml:"branches"`
	EnviromentVariables map[string]string `yaml:"env"`

	PreScript  string `yaml:"pre"`
	PostScript string `yaml:"post"`

	Use  string                 `yaml:"use"`
	With map[string]interface{} `yaml:"with"`
}

//check setup the defaults for the configuration and throws an error if invalid
func (config deployConfig) check() (deployConfig, error) {
	if config.EnviromentVariables == nil {
		config.EnviromentVariables = make(map[string]string)
	}
	if config.With == nil {
		config.With = make(map[string]interface{})
	}
	return config, nil
}

//sshConfig handles configuration for SSH
type sshConfig struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	PrivateKey string `yaml:"key"`
}

//loadConfiguration loads the givne file path
func loadConfiguration(filepath string) (globalConfig, error) {
	config := globalConfig{}

	// Load the file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	// Parse the yaml
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// Check the configuration
	return config.check()
}
