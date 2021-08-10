package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

type deployment interface {
	//deploy runs the deployment within the given system
	deploy(sys system)
}

type sshDeploy struct {
	ssh    sshConfig
	deploy deployment
}

func createDeployment(name string, with map[string]interface{}) deployment {
	var d deployment = nil
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

func (d *scriptDeploy) deploy(sys system) {
	log.Printf("Executing %s\n", d.Script)

	_, err := sys.execShell(d.Shell, d.Script, d.Args...)
	if err != nil {
		log.Println("Failed to deploy script", d.Script, err)
		return
	}

	log.Printf("Finished %s\n", d.Script)
}

type composerDeploy struct {
}

type npmDeploy struct {
	Script string
}

func (d *npmDeploy) deploy(sys system) {

	// Setup the default
	script := d.Script
	if script == "" {
		script = "build"
	}

	var err error
	var result []byte

	log.Printf("Executing NPM Install\n")
	_, err = sys.exec("npm", "i")
	if err != nil {
		log.Println("failed to install dependencies", err)
		return
	}

	log.Printf("Executing NPM Deploy %s\n", script)
	result, err = sys.exec("npm", "run", script)
	if err != nil {
		log.Println("failed to run script", err)
		return
	}

	log.Println("finished running script", string(result))
}

type dotnetDeploy struct {
}
type dockerDeploy struct {
}
