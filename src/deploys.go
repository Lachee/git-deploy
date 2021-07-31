package main

import (
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type deployer interface {
	//deploy runs the deployer with the given working directory and enviroment settings
	deploy(cwd string, env []string)
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
	case "npm":
		d = &npmDeploy{}
		// case "composer":
		// 	return &composerDeploy{}
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

func (d *scriptDeploy) deploy(cwd string, env []string) {
	log.Printf("Executing %s\n", d.Script)

	cmd := shellCommand(d.Shell, d.Script, d.Args...)
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Finished %s\n", d.Script)
}

type composerDeploy struct {
}

type npmDeploy struct {
	Script string
}

func (d *npmDeploy) deploy(cwd string, env []string) {

	// Setup the default
	script := d.Script
	if script == "" {
		script = "build"
	}

	// Execute
	log.Printf("Executing NPM Deploy %s\n", script)
	var err error
	var result []byte

	env = append(os.Environ(), env...)
	npmiCMD := exec.Command("npm", "i")
	npmiCMD.Env = env
	npmRun := exec.Command("npm", "run", script)
	npmRun.Env = env

	err = npmiCMD.Run()
	if err != nil {
		log.Fatalln("Failed to install dependencies: ", err)
		return
	}

	result, err = npmRun.CombinedOutput()
	if err != nil {
		msg := string(result)
		log.Fatalln("Failed to build: ", err, msg)
		return
	}
}

type dotnetDeploy struct {
}
type dockerDeploy struct {
}
