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
	"time"
)

// GoTools wraps some Golang commands with some additional flags for quick use.
type GoTools struct {
	workDir string
	verbose bool
}

// Clean invokes: 'go clean -i ./...'
// See 'go clean --help'
func (g *GoTools) Clean() error {

	if _, err := (&command{Name: "go", WorkDir: g.workDir}).Execute("clean", "-i", "./..."); err != nil {
		return err
	}

	return nil
}

// Install invokes: 'go install -ldflags -X <pkg>.GitCommit=<commit> -X <pkg>.ReleaseVersion=<releaseVersion> -X <pkg>.BuildTimestamp=<current timestamp formatted in RFC3339>'.
// See 'go install --help'
func (g *GoTools) Install(commit, releaseVersion, pkg string) error {

	now := time.Now().UTC()

	vars := []string{
		"-X", fmt.Sprintf("\"%v.GitCommit=%v\"", pkg, commit),
		"-X", fmt.Sprintf("\"%v.BuildTimestamp=%v\"", pkg, now.Format(time.RFC3339)),
		"-X", fmt.Sprintf("\"%v.ReleaseVersion=%v\"", pkg, releaseVersion),
	}

	args := []string{"install", "-ldflags", strings.Join(vars, " ")}

	if _, err := (&command{Name: "go", WorkDir: g.workDir}).Execute(args...); err != nil {
		return err
	}

	return nil
}
