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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLine(t *testing.T) {
	assert.Equal(t, fmt.Sprintln(), NewLine())
}

func TestNestedError_Error(t *testing.T) {
	assert.Equal(t, "FAILED!", (&NestedError{Message: "FAILED!", Cause: nil}).Error())
	assert.Equal(t, "FAILED!", (&NestedError{Message: "FAILED!", Cause: &NestedError{Message: ""}}).Error())
	assert.Equal(t, "FAILED!\n\tCaused by: NESTED!",
		(&NestedError{Message: "FAILED!", Cause: &NestedError{Message: "NESTED!"}}).Error())
}

func TestCommandImpl_Execute_Success(t *testing.T) {
	if wd, err := os.Getwd(); err != nil {
		assert.Fail(t, err.Error())
	} else if out, err := (&command{Name: "go", WorkDir: wd}).Execute("version"); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.True(t, len(out) > 0)
	}
}

func TestCommandImpl_Execute_Failure(t *testing.T) {
	if wd, err := os.Getwd(); err != nil {
		assert.Fail(t, err.Error())
	} else if out, err := (&command{Name: "go", WorkDir: wd}).Execute("invalid_argument"); err == nil {
		assert.Fail(t, "Error expected!")
	} else {
		assert.True(t, len(out) == 0)
	}
}
