/*
Copyright 2017 Ahmed Zaher

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strings"
)

// Git is a context structure for a git command.
type Git struct {
	// WorkDir is the working directory where the command is being executed.
	WorkDir string
	// Verbose shows more verbose output while execution.
	Verbose bool
}

func (g *Git) withGit() Command {
	return NewNamedCommand("git", g.WorkDir)
}

// IsCommitExists takes a git commit hash and returns true if it's found, false otherwise.
// It returns error if something goes wrong while checking.
func (g *Git) IsCommitExists(hash string) (bool, error) {
	var out string
	var err error

	if out, err = g.withGit().Execute("show", "-s", "--format=%H", hash); err != nil {
		if !strings.HasPrefix(err.Error(), "exit status 128 :: fatal: bad object") {
			return false, err
		}
	}

	return strings.Compare(hash, out) == 0 && len(hash) > 0, nil
}

// Checkout checks out the specified git commit hash, it returns error if it fails.
func (g *Git) Checkout(hash string) error {
	_, err := g.withGit().Execute("checkout", fmt.Sprintf("%v", hash))
	return err
}

// Status returns the status of the current git repository targeted by this Git object, the status returned
// is in string format, the string is empty if all is committed, it returns error on failure.
func (g *Git) Status() (string, error) {
	var out string
	var err error
	if out, err = g.withGit().Execute("status", "-s", "-uall"); err != nil {
		return "", err
	}
	return out, nil
}

// GetTagCommit returns the git commit hash of the specified git tag, it returns an error on failure.
func (g *Git) GetTagCommit(tag string) (string, error) {
	var out string
	var err error
	if out, err = g.withGit().Execute("for-each-ref",
		fmt.Sprintf("refs/tags/%v", tag), "--format='%(objectname)'"); err != nil {
		return "", err
	} else if len(strings.Trim(out, "'")) == 0 {
		return "", fmt.Errorf("tag '%v' is not found", tag)
	}
	return strings.Trim(out, "'"), nil
}

// GetBranchCommit returns the git commit hash of the specified git branch, it returns an error on failure.
func (g *Git) GetBranchCommit(branch string) (string, error) {
	var out string
	var err error
	if out, err = g.withGit().Execute("for-each-ref",
		fmt.Sprintf("refs/heads/%v", branch), "--format='%(objectname)'"); err != nil {
		return "", err
	} else if len(strings.Trim(out, "'")) == 0 {
		return "", fmt.Errorf("branch '%v' is not found", branch)
	}
	return strings.Trim(out, "'"), nil
}
