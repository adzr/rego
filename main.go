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
	"os"
	"strings"

	getopt "github.com/kesselborn/go-getopt"
)

const (
	executionErrorCode = 126
)

func fail(code int, message string, args ...string) {
	fmt.Fprintf(os.Stderr, "%v"+NewLine(), message, args)
	os.Exit(code)
}

func print(format string, args ...string) {
	fmt.Fprintf(os.Stdout, format+NewLine(), args)
}

func exit(format string, args ...string) {
	print(format, args...)
	os.Exit(0)
}

func release(conf *configurations) {

	if conf.Verbose {
		print("generating release: %v", conf.Release)
	}

	var err error

	g := &Git{WorkDir: conf.WorkDir, Verbose: conf.Verbose}

	if err = g.Checkout(conf.Commit); err != nil {
		fail(executionErrorCode, err.Error())
	} else {
		print("commit '%v' is checked out, don't forget to switch back to your working reference", conf.Commit)
	}

	if conf.Verbose {
		print("building from commit '%v'", conf.Commit)
	}

	gt := &GoTools{WorkDir: conf.WorkDir, Verbose: conf.Verbose}

	if err = gt.Clean(); err != nil {
		fail(executionErrorCode, err.Error())
	} else if err = gt.Install(conf.Commit, conf.Release, conf.Package); err != nil {
		fail(executionErrorCode, err.Error())
	}
}

func validate(conf *configurations) {
	g := &Git{WorkDir: conf.WorkDir, Verbose: conf.Verbose}

	var status string
	var err error

	if status, err = g.Status(); err != nil {
		fail(executionErrorCode, err.Error())
	}

	if len(status) > 0 {
		fail(executionErrorCode, "Uncommitted/untracked files:%v", status)
	}

	if len(conf.Tag) > 0 {
		if conf.Verbose {
			print("requested tag: %v", conf.Tag)
		}

		if conf.Commit, err = g.GetTagCommit(conf.Tag); err != nil {
			fail(executionErrorCode, err.Error())
		}

		conf.Release = strings.TrimPrefix(conf.Tag, conf.IgnoreTagPrefix)
	} else if len(conf.Commit) > 0 {
		if conf.Verbose {
			print("requested commit: %v", conf.Commit)
		}

		var exists bool

		if exists, err = g.IsCommitExists(conf.Commit); err != nil {
			fail(executionErrorCode, err.Error())
		} else if !exists {
			fail(executionErrorCode, "invalid commit specified")
		}
	} else {
		if conf.Commit, err = g.GetBranchCommit(conf.Branch); err != nil {
			fail(executionErrorCode, err.Error())
		}

		if conf.Verbose {
			print("requested branch '%v' commit: %v", conf.Branch, conf.Commit)
		}
	}

	if conf.Verbose {
		print("target commit: %v", conf.Commit)
	}
}

func read(conf *configurations) {
	if out, err := configure(conf); err != nil {
		if e, ok := err.(*getopt.GetOptError); ok {
			fail(e.ErrorCode, e.Error())
		} else {
			fail(executionErrorCode, err.Error())
		}
	} else if len(out) > 0 {
		exit(out)
	}

	if conf.Verbose {
		print(`
Branch: %v
Commit hash: %v
Tag: %v
Working directory: %v
Release version: %v
Ignore tag prefix: %v
Package: %v"
`, conf.Branch, conf.Commit, conf.Tag, conf.WorkDir, conf.Release, conf.IgnoreTagPrefix, conf.Package)
	}
}

func main() {

	var conf configurations

	read(&conf)

	validate(&conf)

	release(&conf)
}
