package main

import "os/exec"

// contains checks if a string is present in a slice
func containsStr(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//shellCommand creates a new command in the given shell. If no shell is provided, then the OS default will be used.
func shellCommand(shell string, cmd string, args ...string) *exec.Cmd {
	if shell == "" {
		shell = DEFAULT_SHELL
	}
	cmdArgs := append([]string{cmd}, args...)
	return exec.Command(shell, cmdArgs...)
}
