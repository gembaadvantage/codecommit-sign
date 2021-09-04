# codecommit-sign

[![Build status](https://img.shields.io/github/workflow/status/gembaadvantage/codecommit-sign/ci?style=flat-square&logo=go)](https://github.com/gembaadvantage/codecommit-sign/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gembaadvantage/codecommit-sign?style=flat-square)](https://goreportcard.com/report/github.com/gembaadvantage/codecommit-sign)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gembaadvantage/codecommit-sign.svg?style=flat-square)](go.mod)
[![codecov](https://codecov.io/gh/gembaadvantage/codecommit-sign/branch/main/graph/badge.svg)](https://codecov.io/gh/gembaadvantage/codecommit-sign)

DESCRIPTION

## Install

Binary downloads can be found on the [Releases](https://github.com/gembaadvantage/codecommit-sign/releases) page. Unpack the `codecommit-sign` binary and add it to your PATH.

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew tap gembaadvantage/tap
brew install codecommit-sign
```

### Fish

To use [Fish](https://gofi.sh/):

```sh
gofish install codecommit-sign
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop install codecommit-sign
```

## Quick Start

TODO - link to CLI command or retrieve from the AWS console

```sh
codecommit-sign https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/<REPOSITORY_NAME>
```

If not using a tool like aws-vault (environment variables aren't correctly set) then a named profile can be provided:

```sh
codecommit-sign https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/<REPOSITORY_NAME> --profile <AWS_PROFILE>
```

## Examples

TODO

### AWS Lambda

TODO

### AWS CodeBuild

TODO
