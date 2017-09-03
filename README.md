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
