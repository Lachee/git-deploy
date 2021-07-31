package main

import (
	"log"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type deployer interface {
	deploy()
}

type sshDeploy struct {
	ssh    sshConfig
	deploy deployer
}

func createDeployer(name string, with map[interface{}]interface{}) deployer {
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
	var cmd *exec.Cmd
	args := []string{d.Script}
	args = append(args, d.Args...)
	if d.Shell == "" {
		cmd = exec.Command("bash", args...)
	} else {
		cmd = exec.Command(d.Shell, args...)
	}

	log.Printf("Executing %s\n", d.Script)
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
