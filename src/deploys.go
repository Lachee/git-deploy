package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

type deployer interface {
	deploy()
}

type sshDeploy struct {
	ssh    sshConfig
	deploy deployer
}

func createDeployer(name string, with map[string]interface{}) deployer {
	var d deployer = nil
	switch name {
	case "script":
		d = &scriptDeploy{}
		break
		// case "composer":
		// 	return &composerDeploy{}
		// case "npm":
		// 	return &npmDeploy{}
		// case "dotnet":
		// 	return &dotnetDeploy{}
		// case "docker":
		// 	return &dockerDeploy{}
	}

	if d == nil {
		log.Printf("error: there is no deployer for %s\n", name)
		return nil
	}
	data, _ := yaml.Marshal(with)
	yaml.Unmarshal([]byte(data), d)
	return d
}

type scriptDeploy struct {
	Shell  string
	Script string
	Args   []string
}

func (d *scriptDeploy) deploy() {
	log.Printf("Executing %s\n", d.Script)

	cmd := shellCommand(d.Shell, d.Script, d.Args...)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Finished %s\n", d.Script)
}

type composerDeploy struct {
}
type npmDeploy struct {
}
type dotnetDeploy struct {
}
type dockerDeploy struct {
}
