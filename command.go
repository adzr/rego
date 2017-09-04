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
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Command represents a shell command.
type Command interface {
	// Execute executes this command supplied with its arguments as method parameters.
	// It returns the output as string on success, otherwise error.
	Execute(args ...string) (string, error)
}

type namedCommand struct {
	name    string
	workDir string
}

func (comm *namedCommand) Execute(args ...string) (string, error) {
	c := exec.Command(comm.name, args...)
	c.Dir = comm.workDir

	var out string
	var err error
	var stderr bytes.Buffer
	c.Stderr = &stderr

	if output, e := c.Output(); e != nil {
		out, err = strings.TrimSpace(string(output)),
			errors.Wrap(errors.New(strings.TrimSpace(string(stderr.Bytes()))), e.Error())
	} else {
		out, err = strings.TrimSpace(string(output)), nil
	}

	return out, err
}

// NewNamedCommand return an instance of Command.
func NewNamedCommand(name, workDir string) Command {
	return &namedCommand{name: name, workDir: workDir}
}

// NewLine returns the new line character within a string.
func NewLine() string {
	return fmt.Sprintln()
}
