package main

import (
	"os/exec"
	"strings"
)

type git struct {
	system system
}

func newGit(system system) git {
	return git{system: system}
}

//currentBranch of the system root
func (git git) currentBranch() (string, error) {
	data, err := git.system.exec("git", "branch", "--show-current")
	return strings.Trim(string(data), "\n *"), err
}

//hasChanges checks if the git repo has any changes
func (git git) hasChanges() (bool, error) {
	_, err := git.system.exec("git", "diff", "--quiet")
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() == 1 {
			return true, nil
		}
	}
	return false, err
}

//pushStash creates a new stash
func (git git) pushStash() error {
	_, err := git.system.exec("git", "stash", "push", "-u", "-q")
	return err
}

//popStash pops the latest stash
func (git git) popStash() error {
	_, err := git.system.exec("git", "stash", "pop", "-q")
	return err
}

//pull from and integrate with another repository or a local branch. Returns true if there have been a change
func (git git) pull() (bool, error) {
	msg, err := git.system.exec("git", "pull")
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(string(msg), "Already up to date"), nil
}

/*
func gitPull(system system) error {

}

func gitStash(system system, stash string) error {

}

func gitStashPop(system system) error {

}

*/
