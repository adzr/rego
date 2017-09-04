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

	var e error
	var workDirectory string

	if workDirectory, e = os.Getwd(); e != nil {
		return "", e
	}

	parser := getopt.Options{
		Description: "Builds and installs a binary release of a Golang source code while embedding its release information - through a group of exported public variables in the source - based on the current status of its Git repository, all the source files must be committed into the local repository before running this command or it will complain, this tool assumes that Golang (with a valid 'GOROOT' and 'GOPATH' environment variables) and Git source control are installed and fully working though shell.",
		Definitions: []getopt.Option{
			{
				OptionDefinition: "work-directory|w|REGO_WORK_DIR",
				Description:      "The working directory that contains the project source files and its Git repository",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     workDirectory,
			}, {
				OptionDefinition: "branch|b|REGO_BRANCH",
				Description:      "The branch name of where the binary release source is going to be taken from, the command automatically picks the most recent commit hash in the specified branch, the commit hash string is passed to the binary release while building through the public variable 'GitCommit'",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "develop",
			}, {
				OptionDefinition: "commit|c|REGO_COMMIT",
				Description:      "The commit hash string of where the binary release source is going to be taken from, specifying this option causes the '--branch' option to be ignored since this option is more specific, the commit hash string is passed to the binary release while building through the public variable 'GitCommit'",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "tag|t|REGO_TAG",
				Description:      "The tag name of where the binary release source is going to be taken from, causes the '--branch' and '--commit' options to be ignored since this option is more specific, the commit hash string is passed to the binary release while building through the public variable 'GitCommit'",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
			}, {
				OptionDefinition: "release|r|REGO_RELEASE",
				Description:      "The string that is meant to represent the final binary release version, if the '--tag' option is specified this option is automatically calculated with consideration of '--ignore-tag-prefix' option if specified to represent the tag name, the value of this option is passed to the binary release while building through the public variable 'ReleaseVersion'",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "SNAPSHOT",
			}, {
				OptionDefinition: "package|p|REGO_PACKAGE",
				Description: "The package name of which contains the declarations of the public variables" +
					" (GitCommit, BuildTimestamp, ReleaseVersion, GoVersion) which represent the commit hash of where the binary release source has been pulled from, the timestamp of when the build has be triggered, the release version string, the Golang version that has been used in the build, respectively",
				Flags:        getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue: "main",
			}, {
				OptionDefinition: "ignore-tag-prefix|i|REGO_IGNORE_TAG_PREFIX",
				Description:      "If the '--tag' option is specified, this option trims the specified prefix off the tag name while calculating the release version string",
				Flags:            getopt.Optional | getopt.ExampleIsDefault,
				DefaultValue:     "",
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

	var err *getopt.GetOptError
	var options map[string]getopt.OptionValue

	if options, _, _, err = parser.ParseCommandLine(); err != nil {
		return "", fmt.Errorf("failed with error code: %v, %v", err.ErrorCode, err.Error())
	} else if help, wantsHelp := options["help"]; wantsHelp && help.String == "usage" {
		return parser.Usage(), nil
	} else if wantsHelp && help.String == "help" {
		return parser.Help(), nil
	} else if options["version"].Bool {
		return fmt.Sprintf("Release: %v%vCommit: %v%vBuild Time: %v%vBuilt with: %v",
			ReleaseVersion, NewLine(),
			GitCommit, NewLine(),
			BuildTimestamp, NewLine(),
			GoVersion), nil
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
