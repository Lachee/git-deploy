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

/*
 merge two maps together from left to right
 [ a: 10, b: 20 ] <<< [ a: 20, c: 40 ]
 [ a: 20, b: 20, c: 40]
*/
func merge(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

//shellCommand creates a new command in the given shell. If no shell is provided, then the OS default will be used.
func shellCommand(shell string, cmd string, args ...string) *exec.Cmd {
	if shell == "" {
		shell = DEFAULT_SHELL
	}
	cmdArgs := append([]string{cmd}, args...)
	return exec.Command(shell, cmdArgs...)
}
