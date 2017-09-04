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
}

func configure(conf *configurations) (string, error) {

	var err error
	var workDirectory string

	if workDirectory, err = os.Getwd(); err != nil {
		return "", err
	}

	parser := getopt.Options{
		Description: "Builds a release of a Golang source based on the current status of its git repository.",
		Definitions: []getopt.Option{
			{
				OptionDefinition: "work-directory|w|REGO_WORK_DIR",
				Description:      "The working directory that contains the git repository.",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     workDirectory,
			}, {
				OptionDefinition: "branch|b|REGO_BRANCH",
				Description:      "The branch where the release is taken from.",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "develop",
			}, {
				OptionDefinition: "commit|c|REGO_COMMIT",
				Description:      "The commit hash of the snapshot, causes the branch option to be ignored.",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "tag|t|REGO_TAG",
				Description:      "The tag of the final release, causes the branch and commit options to be ignored.",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "release|r|REGO_RELEASE",
				Description:      "The release version, defaults to the most recent tag or to the tag option if specified.",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "SNAPSHOT",
			}, {
				OptionDefinition: "package|p|REGO_PACKAGE",
				Description: "The package name of which contains the definitions of the public variables" +
					" (GitCommit, BuildTimestamp, ReleaseVersion).",
				Flags:        getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue: "main",
			}, {
				OptionDefinition: "ignore-tag-prefix|i|REGO_IGNORE_TAG_PREFIX",
				Description: "Ignores the specified version/tag prefix when reading from the repository" +
					" to write it without prefix in the binary.",
				Flags:        getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue: "",
			}, {
				OptionDefinition: "verbose",
				Description:      "Shows more verbose output.",
				Flags:            getopt.Flag,
				DefaultValue:     false,
			}, {
				OptionDefinition: "version|v",
				Description:      "Prints the version and exits.",
				Flags:            getopt.Flag,
				DefaultValue:     false,
			},
		},
	}

	var options map[string]getopt.OptionValue

	if options, _, _, err = parser.ParseCommandLine(); err != nil {
		return "", err
	} else if help, wantsHelp := options["help"]; wantsHelp && help.String == "usage" {
		return parser.Usage(), nil
	} else if wantsHelp && help.String == "help" {
		return parser.Help(), nil
	} else if options["version"].Bool {
		return fmt.Sprintf("Release: %v%vCommit: %v%vBuild Time: %v%vBuilt with: %v%v",
			ReleaseVersion, NewLine(),
			GitCommit, NewLine(),
			BuildTimestamp, NewLine(),
			GoVersion, NewLine()), nil
	}

	conf.Verbose = options["verbose"].Bool
	conf.WorkDir = strings.TrimSpace(options["work-directory"].String)
	conf.Package = strings.TrimSpace(options["package"].String)
	conf.Branch = strings.TrimSpace(options["branch"].String)
	conf.Commit = strings.TrimSpace(options["commit"].String)
	conf.Tag = strings.TrimSpace(options["tag"].String)
	conf.IgnoreTagPrefix = strings.TrimSpace(options["ignore-tag-prefix"].String)
	conf.Release = strings.TrimSpace(options["release"].String)

	return "", nil
}
