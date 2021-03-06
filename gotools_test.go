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
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	content = `
package main

import (
	"fmt"
)

var GitCommit string
var BuildTimestamp string
var ReleaseVersion string
var GoVersion string

func main() {
	fmt.Printf("Release: %v\nCommit: %v\nBuilt with: %v\n",
		ReleaseVersion,
		GitCommit,
		GoVersion)
}
`
)

type GoToolsTestSuite struct {
	suite.Suite
	goTools *GoTools
	goPath  string
}

func (suite *GoToolsTestSuite) SetupTest() {
	var err error

	if suite.goPath, err = ioutil.TempDir("/tmp/", "test-rego-go-"); err != nil {
		suite.Fail("failed to create temporary directory before test setup", err.Error())
	}

	os.Setenv("GOPATH", suite.goPath)

	suite.goTools = &GoTools{WorkDir: suite.goPath + "/src/project", Verbose: true}

	if err = os.MkdirAll(suite.goTools.WorkDir, os.ModePerm); err != nil {
		suite.Fail("failed to create project source directory before test setup", err.Error())
	}

	if err = ioutil.WriteFile(suite.goTools.WorkDir+"/main.go", []byte(content), 0600); err != nil {
		suite.Fail("failed to create 'main.go'", err.Error())
	}

	git := NewNamedCommand("git", suite.goTools.WorkDir)

	git.Execute("init")
	git.Execute("config", "commit.gpgsign", "false")

	if _, err = git.Execute("add", "main.go"); err != nil {
		suite.Fail("failed to add 'main.go' before test setup", err.Error())
	}

	if _, err = git.Execute("commit", "-n", "-m", "'Initial commit'"); err != nil {
		suite.Fail("failed to commit 'Initial commit' before test setup", err.Error())
	}
}

func (suite *GoToolsTestSuite) TearDownTest() {
	if len(suite.goPath) > 0 {
		if err := os.RemoveAll(suite.goPath); err != nil {
			suite.Fail("failed to remove temporary directory after test teardown", err.Error())
		}
	}
}

func TestGoToolsTestSuite(t *testing.T) {
	suite.Run(t, new(GoToolsTestSuite))
}

func (suite *GoToolsTestSuite) TestGoTools_Clean_Success() {
	var err error

	goCommand := NewNamedCommand(runtime.GOROOT()+"/bin/go", suite.goTools.WorkDir)

	if _, err = goCommand.Execute("install"); err != nil {
		suite.Fail("failed to execute 'go install'", err.Error())
	}

	if err = suite.goTools.Clean(); err != nil {
		suite.Fail(err.Error())
	}

	assert.Nil(suite.T(), err)
}

func (suite *GoToolsTestSuite) TestGoTools_Clean_Failure() {
	dir := suite.goTools.WorkDir
	suite.goTools.WorkDir = "/some-nonexistent-path/"
	assert.NotNil(suite.T(), suite.goTools.Clean())
	suite.goTools.WorkDir = dir
}

func (suite *GoToolsTestSuite) TestGoTools_Install_Success() {
	var err error
	var commit string
	var out string
	var goVersion string

	if commit, err = NewNamedCommand("git", suite.goTools.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get last commit", err.Error())
	}

	if err = suite.goTools.Install(commit, "1.0", "main"); err != nil {
		suite.Fail("failed to install binary", err.Error())
	}

	commandName := strings.Replace(suite.goTools.WorkDir, "/src/project", "/bin/project", 1)

	if out, err = NewNamedCommand(commandName, suite.goTools.WorkDir).Execute(); err != nil {
		suite.Fail("failed to execute output binary", err.Error())
	}

	if goVersion, err = suite.goTools.withGo().Execute("version"); err != nil {
		suite.Fail("failed to get go version", err.Error())
	}

	expected := fmt.Sprintf("Release: %v\nCommit: %v\nBuilt with: %v",
		"1.0",
		commit,
		goVersion)

	assert.Equal(suite.T(), expected, out)
}

func (suite *GoToolsTestSuite) TestGoTools_Install_Failure() {
	var err error
	var commit string

	if commit, err = NewNamedCommand("git", suite.goTools.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get last commit", err.Error())
	}

	dir := suite.goTools.WorkDir
	suite.goTools.WorkDir = "/some-nonexistent-path/"
	assert.NotNil(suite.T(), suite.goTools.Install(commit, "1.0", "main"))
	suite.goTools.WorkDir = dir
}
