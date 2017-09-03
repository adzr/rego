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
	"runtime"
	"strings"
	"time"
)

// GoTools wraps some Golang commands with some additional flags for quick use.
type GoTools struct {
	// WorkDir is the working directory where the command is being executed.
	WorkDir string
	// Verbose shows more verbose output while execution.
	Verbose bool
}

func (g *GoTools) withGo() Command {
	return NewNamedCommand(runtime.GOROOT()+"/bin/go", g.WorkDir)
}

// Clean invokes: 'go clean -i ./...'.
// See 'go clean --help'
func (g *GoTools) Clean() error {

	if _, err := g.withGo().Execute("clean", "-i", "./..."); err != nil {
		return err
	}

	return nil
}

// Install invokes: 'go install -ldflags -X <pkg>.GitCommit=<commit> -X <pkg>.ReleaseVersion=<releaseVersion> -X <pkg>.BuildTimestamp=<current timestamp formatted in RFC3339>'.
// See 'go install --help'
func (g *GoTools) Install(commit, releaseVersion, pkg string) error {

	var err error
	var goVersion string

	now := time.Now().UTC()

	goVersion, _ = g.withGo().Execute("version")

	vars := []string{
		"-X", fmt.Sprintf("\"%v.GitCommit=%v\"", pkg, commit),
		"-X", fmt.Sprintf("\"%v.BuildTimestamp=%v\"", pkg, now.Format(time.RFC3339)),
		"-X", fmt.Sprintf("\"%v.ReleaseVersion=%v\"", pkg, releaseVersion),
		"-X", fmt.Sprintf("\"%v.GoVersion=%v\"", pkg, goVersion),
	}

	args := []string{"install", "-ldflags", strings.Join(vars, " ")}

	if _, err = NewNamedCommand("go", g.WorkDir).Execute(args...); err != nil {
		return err
	}

	return nil
}
