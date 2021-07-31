package main

import "log"

type deployer interface {
	deploy(opts map[string]string)
}

type sshDeploy struct {
	ssh    sshConfig
	deploy deployer
}

func createDeployer(name string) deployer {
	switch name {
	case "script":
		return &scriptDeploy{}
		// case "composer":
		// 	return &composerDeploy{}
		// case "npm":
		// 	return &npmDeploy{}
		// case "dotnet":
		// 	return &dotnetDeploy{}
		// case "docker":
		// 	return &dockerDeploy{}
	}

	log.Printf("error: there is no deployer for %s\n", name)
	return nil
}

type scriptDeploy struct {
	script string
}

func (d *scriptDeploy) deploy(opts map[string]string) {

}

type composerDeploy struct {
}
type npmDeploy struct {
}
type dotnetDeploy struct {
}
type dockerDeploy struct {
}
