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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitTestSuite struct {
	suite.Suite
	git   *Git
	touch Command
}

func (suite *GitTestSuite) SetupTest() {
	var dir string
	var err error

	if dir, err = ioutil.TempDir("", "test_rego_git_"); err != nil {
		suite.Fail("failed to create temporary directory before test setup", err.Error())
		return
	}

	suite.touch = &namedCommand{workDir: dir, name: "touch"}

	suite.git = &Git{Verbose: true, WorkDir: dir}

	git := NewNamedCommand("git", dir)
	git.Execute("init")
	git.Execute("config", "commit.gpgsign", "false")

	if _, err = suite.touch.Execute("README"); err != nil {
		suite.Fail("failed to touch 'README' before test setup", err.Error())
	}

	if _, err = git.Execute("add", "README"); err != nil {
		suite.Fail("failed to add 'README' before test setup", err.Error())
	}

	if _, err = git.Execute("commit", "-n", "-m", "'Initial commit'"); err != nil {
		suite.Fail("failed to commit 'Initial commit' before test setup", err.Error())
	}

	if _, err = git.Execute("checkout", "-b", "develop"); err != nil {
		suite.Fail("failed to checkout new branch 'develop' before test setup", err.Error())
	}

	if _, err = suite.touch.Execute("empty.go"); err != nil {
		suite.Fail("failed to touch 'empty.go' before test setup", err.Error())
	}

	if _, err = git.Execute("add", "empty.go"); err != nil {
		suite.Fail("failed to add 'empty.go' before test setup", err.Error())
	}

	if _, err = git.Execute("commit", "-n", "-m", "'Adding empty.go'"); err != nil {
		suite.Fail("failed to commit 'Adding empty.go' before test setup", err.Error())
	}

	if _, err = git.Execute("checkout", "master"); err != nil {
		suite.Fail("failed to checkout branch 'master' before test setup", err.Error())
	}

	if _, err = git.Execute("tag", "v1.0"); err != nil {
		suite.Fail("failed to commit 'Adding empty.go' before test setup", err.Error())
	}
}

func (suite *GitTestSuite) TearDownTest() {
	if suite.git != nil {
		if len(suite.git.WorkDir) > 0 {
			if err := os.RemoveAll(suite.git.WorkDir); err != nil {
				suite.Fail("failed to remove temporary directory after test teardown", err.Error())
			}
		}
	}
}

func TestGitTestSuite(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}

func (suite *GitTestSuite) TestGit_Status_Committed() {
	var status string
	var err error

	if status, err = suite.git.Status(); err != nil {
		suite.Fail("failed to get status", err.Error())
	}

	assert.Empty(suite.T(), status)
}

func (suite *GitTestSuite) TestGit_Status_Uncommitted() {
	var status string
	var err error

	if _, err = suite.touch.Execute("dirty.go"); err != nil {
		suite.Fail("failed to touch 'dirty.go'", err.Error())
	}

	if status, err = suite.git.Status(); err != nil {
		suite.Fail("failed to get status", err.Error())
	}

	assert.NotEmpty(suite.T(), status)
}

func (suite *GitTestSuite) TestGit_Status_Failure() {
	var err error
	var status string

	if err = os.RemoveAll(suite.git.WorkDir); err != nil {
		suite.Fail("failed to remove work directory", err.Error())
	}

	status, err = suite.git.Status()
	assert.Empty(suite.T(), status)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_IsCommitExists_FailureEmpty() {
	var err error
	var exists bool

	exists, err = suite.git.IsCommitExists("")
	assert.False(suite.T(), exists)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_IsCommitExists_FailureNotEmpty() {
	var err error
	var exists bool

	exists, err = suite.git.IsCommitExists("invalid_hash")
	assert.False(suite.T(), exists)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_IsCommitExists_FailureNotExist() {
	var err error
	var exists bool

	exists, err = suite.git.IsCommitExists("4f0c1d3161c94c10847e96c79a1806836b1bad12")
	assert.False(suite.T(), exists)
	assert.Nil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_IsCommitExists_SuccessExists() {
	var err error
	var exists bool
	var hash string

	if hash, err = NewNamedCommand("git", suite.git.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get recent commit")
		return
	}

	exists, err = suite.git.IsCommitExists(hash)
	assert.True(suite.T(), exists)
	assert.Nil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_Checkout_Failure() {
	err := suite.git.Checkout("4f0c1d3161c94c10847e96c79a1806836b1bad12")
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_Checkout_Success() {
	var err error
	var hash string

	if hash, err = NewNamedCommand("git", suite.git.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get recent commit")
		return
	}

	err = suite.git.Checkout(hash)
	assert.Nil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetTagCommit_Success() {
	var err error
	var hash string
	var commit string

	if hash, err = NewNamedCommand("git", suite.git.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get recent commit")
		return
	}

	commit, err = suite.git.GetTagCommit("v1.0")
	assert.Equal(suite.T(), hash, commit)
	assert.Nil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetTagCommit_FailureNotFound() {
	var err error
	var commit string

	commit, err = suite.git.GetTagCommit("v1.1")
	assert.Empty(suite.T(), commit)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetTagCommit_FailureNoRepo() {
	var err error
	var commit string

	if err = os.RemoveAll(suite.git.WorkDir); err != nil {
		suite.Fail("failed to remove work directory", err.Error())
	}

	commit, err = suite.git.GetTagCommit("v1.1")
	assert.Empty(suite.T(), commit)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetBranchCommit_FailureNoRepo() {
	var err error
	var commit string

	if err = os.RemoveAll(suite.git.WorkDir); err != nil {
		suite.Fail("failed to remove work directory", err.Error())
	}

	commit, err = suite.git.GetBranchCommit("master")
	assert.Empty(suite.T(), commit)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetBranchCommit_FailureNotFound() {
	var err error
	var commit string

	commit, err = suite.git.GetBranchCommit("release")
	assert.Empty(suite.T(), commit)
	assert.NotNil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetBranchCommit_SuccessMaster() {
	var err error
	var hash string
	var commit string

	if hash, err = NewNamedCommand("git", suite.git.WorkDir).Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get recent commit")
		return
	}

	commit, err = suite.git.GetBranchCommit("master")
	assert.Equal(suite.T(), hash, commit)
	assert.Nil(suite.T(), err)
}

func (suite *GitTestSuite) TestGit_GetBranchCommit_SuccessDevelop() {
	var err error
	var hash string
	var commit string

	g := NewNamedCommand("git", suite.git.WorkDir)

	if _, err = g.Execute("checkout", "develop"); err != nil {
		suite.Fail("failed to checkout branch 'develop'", err.Error())
	}

	if hash, err = g.Execute("show", "-s", "--format=%H"); err != nil {
		suite.Fail("failed to get recent commit")
		return
	}

	commit, err = suite.git.GetBranchCommit("develop")
	assert.Equal(suite.T(), hash, commit)
	assert.Nil(suite.T(), err)
}
