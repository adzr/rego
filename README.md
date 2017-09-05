# rego

Rego (release-go) is a command line tool to help building [Golang](https://golang.org) code committed under a [Git](https://git-scm.com/) repository into binaries with embedded release information.

[![Build Status](https://travis-ci.org/adzr/rego.svg?branch=master)](https://travis-ci.org/adzr/rego) [![Coverage Status](https://coveralls.io/repos/github/adzr/rego/badge.svg?branch=master)](https://coveralls.io/github/adzr/rego?branch=master)


## Brief

Simply, instead of running `go install`, `rego` can be executed against the desired a project directory, and it will feed the output binary with the release information specified in the command arguments.

## Installation

```bash
go get -u github.com/adzr/rego
```

## Usage

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
For detailed help type ```rego --help```

## Example
Create a new Golang project named ```example-go```, initialize a new git repository and add a ```main.go``` file
```bash
cd $GOPATH/src/ \
&& mkdir example-go \
&& cd example-go \
&& git init \
&& touch main.go
```
Edit the recently created file ```main.go``` with your favorite editor and paste the following code then save the file
```golang
package main

import (
	"fmt"
)

var GitCommit string
var BuildTimestamp string
var ReleaseVersion string
var GoVersion string

func main() {
	fmt.Printf("Release: %v\nCommit: %v\nBuild Time: %v\nBuilt with: %v\n",
		ReleaseVersion, GitCommit, BuildTimestamp, GoVersion)
}

```
Now in the command line type
```bash
rego -w $GOPATH/src/example-go
```
Output
```
Uncommitted/untracked files:
 ?? main.go
```
Oops, looks like we missed something, it appears that we haven't committed all our files, let's commit them and try again
```bash
cd $GOPATH/src/example-go \
&& git add . \
&& git commit -m 'Initial commit' \
&& rego -r 1.0
```
Have we succeeded?
```
branch 'develop' is not found
```
Aaa..nope, it seems that ```rego``` is trying by default to pull from a branch named ```develop```, and that's because it tries to follow [GitFlow](https://datasift.github.io/gitflow/IntroducingGitFlow.html) branching model, releasing from ```develop``` branch is what you will mostly be doing before merging to the master, but our tiny example project has only a ```master``` branch so let's be more specific here
```bash
rego -r 1.0 -b master -w $GOPATH/src/example-go
```
And here is what we get back
```
commit '03c7ac7ddd8563cf513a5925c85193c405d66c12' is checked out, don't forget to switch back to your working reference
```
It seems that it picked the most recent commit in our selected branch and checked it out for releasing, it also notifies us not to forget to check out our previous branch back again so we don't continue working/committing to an unreferenced branch, so let's checkout the ```master``` branch again.
```bash
git checkout master
```
Ok, now we're back on master, but what's happened with our ```example-go``` binary that we were trying to build? let's check out by running it
```bash
$GOPATH/bin/example-go
```
Output
```
Release: 1.0
Commit: 134492f9327867b715fdc552993305179b7bc23f
Build Time: 2017-09-04T19:07:57Z
Built with: go version go1.9 linux/amd64
```
Finally, we can see our code is built and embedding the correct release information, so now you can try to play more with the command options to see different results, e.g like a different release version (which defaults to SNAPSHOT if not specified) or try to tag your commit and pass the tag name as an option to the command, so refer back to the help page for more information by typing ```rego --help```.

## Contributing
Pull requests and issue reports are welcomed.

## License
This project is licensed under [Apache License Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.txt)
