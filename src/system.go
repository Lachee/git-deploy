package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type system interface {
	//rootDirectory the system is working from
	rootDirectory() string

	//exec runs a command on the system
	exec(cmd string, args ...string) ([]byte, error)

	//execShell using a specific shell
	execShell(shell string, cmd string, args ...string) ([]byte, error)

	//read a file from the system
	read(filePath string) ([]byte, error)

	//write a file to the system
	//write(filePath string, data []byte) error

	//setEnviromentVariables sets the _additional_ enviroment variables to be sent to commands
	setEviromentVariables(env []string)
}

type localSystem struct {
	root string
	envs []string
}

func newLocalSystem(workingDirectory string) *localSystem {
	return &localSystem{
		root: workingDirectory,
		envs: make([]string, 0),
	}
}

//rootDirectory the system is working from
func (system *localSystem) rootDirectory() string { return system.root }

//exec runs a command on the system
func (system *localSystem) exec(cmd string, args ...string) ([]byte, error) {
	os.Chdir(system.rootDirectory())
	execCMD := exec.Command(cmd, args...)
	return execCMD.Output()
}

//execShell using a specific shell
func (system *localSystem) execShell(shell string, cmd string, args ...string) ([]byte, error) {
	os.Chdir(system.rootDirectory())
	if shell == "" {
		shell = DEFAULT_SHELL
	}
	execCMD := exec.Command(shell, append([]string{cmd}, args...)...)
	return execCMD.Output()
}

//read a file from the system
func (system *localSystem) read(filePath string) ([]byte, error) {
	os.Chdir(system.rootDirectory())
	return ioutil.ReadFile(filePath)
}

//write a file to the system
//func (system *localSystem) write(filePath string, data []byte) error { return nil }

//setEnviromentVariables sets the _additional_ enviroment variables to be sent to commands
func (system *localSystem) setEviromentVariables(envs []string) {
	system.envs = envs
}
