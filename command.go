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
)

// Command represents a shell command.
type Command interface {
	// Execute executes this command supplied with its arguments as method parameters.
	// It returns the output as string on success, otherwise error.
	Execute(args ...string) (string, error)
}

// command is a default Command implementation that gets executed
// against a certain working directory.
type command struct {
	// Name is the command name.
	Name string
	// WorkDir is the directory path meant to executed the command against.
	WorkDir string
}

// NestedError is a recursive structure meant to store the error stack.
type NestedError struct {
	Message string
	Cause   *NestedError
}

// Error displays the error stack stored in the NestedError structure.
func (e *NestedError) Error() string {
	if e.Cause != nil {
		if message := e.Cause.Error(); len(message) > 0 {
			return fmt.Sprintf("%v%v\tCaused by: %v", e.Message, NewLine(), message)
		}
	}

	return e.Message
}

func (comm *command) Execute(args ...string) (string, error) {
	c := exec.Command(comm.Name, args...)
	c.Dir = comm.WorkDir

	var out string
	var err error
	var stderr bytes.Buffer
	c.Stderr = &stderr

	if output, e := c.Output(); e != nil {
		out, err = strings.TrimSpace(string(output)),
			&NestedError{Message: e.Error(), Cause: &NestedError{Message: strings.TrimSpace(string(stderr.Bytes()))}}
	} else {
		out, err = strings.TrimSpace(string(output)), nil
	}

	return out, err
}

// NewLine returns the new line character within a string.
func NewLine() string {
	return fmt.Sprintln()
}
