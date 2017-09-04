REGO
===
Rego (release-go) is a command line tool to help building [Golang](https://golang.org) code committed under a [Git](https://git-scm.com/) repository into binaries with release information.

[![Build Status](https://travis-ci.org/adzr/rego.svg?branch=master)](https://travis-ci.org/adzr/rego)
[![Coverage Status](https://coveralls.io/repos/github/adzr/rego/badge.svg?branch=master)](https://coveralls.io/github/adzr/rego?branch=master)
---

Brief
-----
Simply, instead of running `go install`, `rego` can be executed against the desired a project directory, and it will feed the output binary with the release information specified in the command arguments.

Only as a prerequisite the developer has to define the following public variables in his project main package:

```golang
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

```

Command usage can be as follows:

```
Usage: rego [-w <work-directory>] [-b <branch>] [-c <commit>] [-t <tag>] [-r <release>] [-p <package>] [-i <ignore-tag-prefix>] --verbose -v

Builds a release of a Golang source based on the current status of its git repository.

Options:
    -w, --work-directory=<work-directory>         The working directory that contains the git repository (default: /home/adzr/Documents/code/golang/src/github.com/adzr/rego); setable via $REGO_WORK_DIR
    -b, --branch=<branch>                         The branch where the release is taken from (default: develop); setable via $REGO_BRANCH
    -c, --commit=<commit>                         The commit hash of the snapshot, causes the branch option to be ignored; setable via $REGO_COMMIT
    -t, --tag=<tag>                               The tag of the final release, causes the branch and commit options to be ignored; setable via $REGO_TAG
    -r, --release=<release>                       The release version, defaults to the most recent tag or to the tag option if specified (default: SNAPSHOT); setable via $REGO_RELEASE
    -p, --package=<package>                       The package name of which contains the definitions of the public variables (GitCommit, BuildTimestamp, ReleaseVersion) (default: main); setable via $REGO_PACKAGE
    -i, --ignore-tag-prefix=<ignore-tag-prefix>   Ignores the specified version/tag prefix when reading from the repository to write it without prefix in the binary; setable via $REGO_IGNORE_TAG_PREFIX
        --verbose                                 Shows more verbose output
    -v, --version                                 Prints the version and exits
    -h, --help                                    usage (-h) / detailed help text (--help)

```
