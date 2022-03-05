/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package translate

import (
	"errors"
	"regexp"
	"strings"
)

var (
	urlRgx = regexp.MustCompile(`^https://(.+@)?git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	grcRgx = regexp.MustCompile(`^codecommit:(:.+:)?//(.+)$`)
)

// Remote provides both details about how the repository is hosted within AWS
// CodeCommit and it will be accessed
type Remote struct {
	// Repository contains the name of the target repository
	Repository string

	// Region identifies if an AWS Region was used when targeting the remote
	Region string

	// Profile identifies if a named AWS Profile was used when targeting the remote
	Profile string
}

// RemoteHTTPS identifies details about an AWS CodeCommit remote based on the provided
// HTTPS clone URL
func RemoteHTTPS(url string) (Remote, error) {
	m := urlRgx.FindStringSubmatch(url)
	if len(m) < 4 {
		return Remote{}, errors.New("malformed codecommit HTTPS URL")
	}

	return Remote{
		Repository: m[len(m)-1],
		Region:     m[2],
	}, nil
}

// RemoteGRC identifies details about an AWS CodeCommit remote based on the provided
// GRC (git-remote-codecommit) clone URL
func RemoteGRC(url string) (Remote, error) {
	m := grcRgx.FindStringSubmatch(url)
	if len(m) < 3 {
		return Remote{}, errors.New("malformed codecommit GRC URL")
	}

	rem := Remote{
		Region: strings.ReplaceAll(m[1], ":", ""),
	}

	// GRC supports prefixing the repository name with an optional AWS profile
	rem.Repository = m[len(m)-1]
	if strings.Contains(rem.Repository, "@") {
		p := strings.Split(rem.Repository, "@")

		rem.Repository = p[1]
		rem.Profile = p[0]
	}

	return rem, nil
}
