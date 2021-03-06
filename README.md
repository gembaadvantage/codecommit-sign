# codecommit-sign

[![Build status](https://img.shields.io/github/workflow/status/gembaadvantage/codecommit-sign/ci?style=flat-square&logo=go)](https://github.com/gembaadvantage/codecommit-sign/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gembaadvantage/codecommit-sign?style=flat-square)](https://goreportcard.com/report/github.com/gembaadvantage/codecommit-sign)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gembaadvantage/codecommit-sign.svg?style=flat-square)](go.mod)
[![codecov](https://codecov.io/gh/gembaadvantage/codecommit-sign/branch/main/graph/badge.svg)](https://codecov.io/gh/gembaadvantage/codecommit-sign)

Generate a signed AWS V4 CodeCommit URL without the need for dedicated IAM user credentials.

## Install

Binary downloads can be found on the [Releases](https://github.com/gembaadvantage/codecommit-sign/releases) page. Unpack the `codecommit-sign` binary and add it to your PATH.

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew tap gembaadvantage/tap
brew install codecommit-sign
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop install codecommit-sign
```

### Yum

To install using the yum package manager:

```sh
echo '[codecommit-sign]
name=uplift
baseurl=https://yum.fury.io/ga-paul-t/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/codecommit-sign.repo
sudo yum install -y codecommit-sign

```

### Apt

```sh
echo 'deb [trusted=yes] https://apt.fury.io/ga-paul-t/ /' | sudo tee /etc/apt/sources.list.d/codecommit-sign.list
sudo apt update
sudo apt install -y codecommit-sign
```

### Aur

To install from the [aur](https://archlinux.org/) using [yay](https://github.com/Jguer/yay):

```sh
yay -S codecommit-sign-bin
```

### Linux Packages

Download and manually install one of the .deb, .rpm or .apk packages from the [Releases](https://github.com/gembaadvantage/codecommit-sign/releases) page.

```sh
sudo apt install codecommit-sign_*.deb
```

```sh
sudo yum localinstall codecommit-sign-*.rpm
```

```sh
sudo apk add --no-cache --allow-untrusted codecommit-sign_*.apk
```

### Script

To install using a shell script:

```sh
curl https://raw.githubusercontent.com/gembaadvantage/codecommit-sign/main/scripts/install > install
chmod 700 install
./install
```

## Quick Start

Retrieve (_or construct_) the clone URL to your chosen CodeCommit repository and then sign it. Depending on your chosen authentication mechanism, you may need to provide an AWS named profile through the optional `--profile` flag.

### HTTPS

```sh
codecommit-sign https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository
```

### GRC

```sh
codecommit-sign codecommit::eu-west-1://repository
```

All GRC variants are supported:

- `codecommit://repository`
- `codecommit://profile@repository`
- `codecommit::region://profile@repository`
