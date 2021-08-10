package main

type remoteSystem struct {
	root   string
	config sshConfig
}

func newRemoteSystem(workingDirectory string, ssh sshConfig) *remoteSystem {
	return &remoteSystem{
		root:   workingDirectory,
		config: ssh,
	}
}

func (system *remoteSystem) rootDirectory() string {
	return system.root
}
