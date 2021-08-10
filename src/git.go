package main

import "strings"

//gitCurrentBranch gets the current git branch
func gitCurrentBranch(system system) (string, error) {
	data, err := system.exec("git", "branch", "--show-current")
	return strings.Trim(string(data), "\n *"), err
}
