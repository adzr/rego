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

// GitCommit is the git commit hash string,
// gets passed from the command line using a binary release of this tool.
var GitCommit string

// BuildTimestamp is the current timestamp in a string format,
// gets passed from the command line using a binary release of this tool.
var BuildTimestamp string

// ReleaseVersion is the desired release version string that represents the version of this executable.
// gets passed from the command line using a binary release of this tool.
var ReleaseVersion string

// GoVersion indicates which version of Go has been used to build this binary.
// gets passed from the command line using a binary release of this tool.
var GoVersion string

type configurations struct {
	WorkDir         string
	Tag             string
	Release         string
	IgnoreTagPrefix string
	Package         string
	Commit          string
	Branch          string
	Verbose         bool
	Version         bool
}

func configure(conf *configurations) {

	workDirectory := ""

	if wd, err := os.Getwd(); err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	} else {
		workDirectory = wd
	}

	parser := getopt.Options{
		Description: "Builds a release of a Golang source based on the current status of its git repository.",
		Definitions: []getopt.Option{
			{
				OptionDefinition: "work-directory|w|REGO_WORK_DIR",
				Description:      "The working directory that contains the git repository",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     workDirectory,
			}, {
				OptionDefinition: "branch|b|REGO_BRANCH",
				Description:      "The branch where the release be taken from",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "develop",
			}, {
				OptionDefinition: "commit|c|REGO_COMMIT",
				Description:      "The commit hash of the snapshot, overrides the branch option",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "tag|t|REGO_TAG",
				Description:      "The tag of the final release, overrides the branch and commit options",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "release|r|REGO_RELEASE",
				Description:      "The release version, defaults to the most recent tag",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "package|p|REGO_PACKAGE",
				Description: "The package name of which contains the definitions of the public variables" +
					" (GitCommit, BuildTimestamp, ReleaseVersion)",
				Flags:        getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue: "main",
			}, {
				OptionDefinition: "ignore-tag-prefix|i|REGO_IGNORE_TAG_PREFIX",
				Description: "Ignores the specified version/tag prefix when reading from the repository" +
					" to write it without prefix in the binary",
				Flags:        getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue: "",
			}, {
				OptionDefinition: "verbose",
				Description:      "Shows more verbose output",
				Flags:            getopt.Flag,
				DefaultValue:     false,
			}, {
				OptionDefinition: "version|v",
				Description:      "Prints the version and exits",
				Flags:            getopt.Flag,
				DefaultValue:     false,
			},
		},
	}

	if options, _, _, err := parser.ParseCommandLine(); err != nil {
		fmt.Println(err.Error())
		os.Exit(err.ErrorCode)
	} else if help, wantsHelp := options["help"]; wantsHelp && help.String == "usage" {
		fmt.Print(parser.Usage())
		os.Exit(0)
	} else if wantsHelp && help.String == "help" {
		fmt.Print(parser.Help())
		os.Exit(0)
	} else {
		conf.WorkDir = options["work-directory"].String
		conf.Branch = options["branch"].String
		conf.Commit = options["commit"].String
		conf.Tag = options["tag"].String
		conf.Release = options["release"].String
		conf.Package = options["package"].String
		conf.IgnoreTagPrefix = options["ignore-tag-prefix"].String
		conf.Verbose = options["verbose"].Bool
		conf.Version = options["version"].Bool
	}
}

func fail(message string) {
	println(message)
	os.Exit(1)
}

func useTag(conf *configurations) bool {
	return len(conf.Tag) > 0
}

func useCommit(conf *configurations) bool {
	return len(conf.Commit) > 0
}

func useRelease(conf *configurations) bool {
	return len(conf.Release) > 0
}

func generateRelease(conf *configurations) string {
	if useTag(conf) {
		if useRelease(conf) {
			return conf.Release
		}

		return strings.TrimPrefix(conf.Release, conf.IgnoreTagPrefix)
	}

	if useRelease(conf) {
		return conf.Release
	}

	return "SNAPSHOT"
}

func release(commit, releaseVersion, pkg, workDir string, verbose bool) error {
	goTools := &GoTools{WorkDir: workDir, Verbose: verbose}

	if err := goTools.Clean(); err != nil {
		return err
	} else if err := goTools.Install(commit, releaseVersion, pkg); err != nil {
		return err
	} else {
		return nil
	}
}

func reportInput(conf *configurations) string {
	var target string

	if useTag(conf) {
		target = fmt.Sprintf("Tag: %v", conf.Tag)
	} else if useCommit(conf) {
		target = fmt.Sprintf("Commit: %v", conf.Commit)
	} else {
		target = fmt.Sprintf("Branch: %v", conf.Branch)
	}

	return fmt.Sprintf("Working directory: %v%vTarget: %v%vRelease Version: %v%vIgnore tag prefix: %v%vPackage: %v%v",
		conf.WorkDir, NewLine(),
		target, NewLine(),
		conf.Release, NewLine(),
		conf.IgnoreTagPrefix, NewLine(),
		conf.Package, NewLine())
}

func printVersion() (string, error) {
	return fmt.Sprintf("Release: %v%vCommit: %v%vBuild Time: %v%vBuilt with: %v%v",
		ReleaseVersion, NewLine(),
		GitCommit, NewLine(),
		BuildTimestamp, NewLine(),
		GoVersion, NewLine()), nil
}

func assertGitStatus(conf *configurations) error {
	g := &Git{WorkDir: conf.WorkDir, Verbose: conf.Verbose}

	var status string
	var err error

	if status, err = g.Status(); err != nil {
		return err
	}

	if len(status) > 0 {
		return fmt.Errorf("Uncommitted/untracked files:%v%v", NewLine(), status)
	}

	return nil
}

func main() {

	var conf configurations

	configure(&conf)

	if conf.Version {
		fmt.Print(printVersion())
		os.Exit(0)
	} else if conf.Verbose {
		fmt.Println(reportInput(&conf))
	}

	if err := assertGitStatus(&conf); err != nil {
		fail(err.Error())
	}

	var commit, rls string
	var err error

	g := &Git{WorkDir: conf.WorkDir, Verbose: conf.Verbose}

	if useTag(&conf) {
		if conf.Verbose {
			fmt.Printf("Using tag: %v%v", conf.Tag, NewLine())
		}

		if commit, err = g.GetTagCommit(conf.Tag); err != nil {
			fail(err.Error())
		}
	} else if useCommit(&conf) {
		if conf.Verbose {
			fmt.Printf("Using commit: %v%v", conf.Commit, NewLine())
		}

		if exists, err := g.IsCommitExists(conf.Commit); err != nil {
			fail(err.Error())
		} else if !exists {
			fail("Invalid commit specified.")
		} else {
			commit = conf.Commit
		}
	} else if c, err := g.GetBranchCommit(conf.Branch); err != nil {
		if conf.Verbose {
			fmt.Printf("Using branch '%v' commit: %v%v", conf.Branch, conf.Commit, NewLine())
		}

		fail(err.Error())
	} else {
		if conf.Verbose {
			fmt.Printf("Target commit: %v%v", conf.Commit, NewLine())
		}

		commit = c
	}

	if conf.Verbose {
		fmt.Printf("Generating release: %v%v", conf.Release, NewLine())
	}

	rls = generateRelease(&conf)

	if err := g.Checkout(commit); err != nil {
		fail(err.Error())
	} else {
		fmt.Printf("Commit '%v' is checked out, don't forget to switch back to your working reference.%v", commit, NewLine())
	}

	if conf.Verbose {
		fmt.Printf("Building from commit '%v'%v", commit, NewLine())
	}

	if err := release(commit, rls, conf.Package, conf.WorkDir, conf.Verbose); err != nil {
		fail(err.Error())
	}
}
