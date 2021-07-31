package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type globalConfig struct {
	Projects []projectConfig `yaml:"projects"`
}

type projectConfig struct {
	Name             string            `yaml:"name"`
	ProjectDirectory string            `yaml:"project"`
	ConfigPath       string            `yaml:"config"`
	Key              string            `yaml:"key"`
	Providers        []string          `yaml:"providers"`
	Webhook          string            `yaml:"webhook"`
	Env              map[string]string `yaml:"env"`
}

type deployConfig struct {
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

	if config.Projects == nil || len(config.Projects) == 0 {
		return config, errors.New("failed to load any projects")
	}

	return config, err
}
