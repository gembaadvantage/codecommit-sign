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
	"fmt"
	"regexp"
	"strings"
)

var (
	urlRgx = regexp.MustCompile(`^https://(.+@)?git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	grcRgx = regexp.MustCompile(`^codecommit::(.+)://(.+)$`)
)

// ToGrc translates a CodeCommit HTTPS URL to a compatible CodeCommit (git-remote-codecommit)
// GRC based URL that can be used to fetch and push changes to a CodeCommit repository
func ToGrc(url string) (string, error) {
	m := urlRgx.FindStringSubmatch(url)
	if len(m) < 4 {
		return "", errors.New("malformed codecommit HTTPS URL")
	}

	region := m[2]
	repo := m[len(m)-1]

	return fmt.Sprintf("codecommit::%s://%s", region, repo), nil
}

// FromGrc translates a CodeCommit (git-remote-codecommit) GRC URL to a compatible HTTPS URL
// that can be used to fetch and push changes to a CodeCommit repository
func FromGrc(url string) (string, error) {
	m := grcRgx.FindStringSubmatch(url)
	if len(m) < 3 {
		return "", errors.New("malformed codecommit grc URL")
	}

	region := m[1]
	repo := m[len(m)-1]
	if strings.Contains(repo, "@") {
		repo = strings.Split(repo, "@")[1]
	}

	return fmt.Sprintf("https://git-codecommit.%s.amazonaws.com/v1/repos/%s", region, repo), nil
}
