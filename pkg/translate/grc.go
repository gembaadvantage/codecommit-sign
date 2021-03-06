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
	"os"
)

// ToGrc translates a CodeCommit HTTPS URL to a compatible CodeCommit (git-remote-codecommit)
// GRC based URL that can be used to fetch and push changes to a CodeCommit repository
func ToGRC(url string) (string, error) {
	rem, err := RemoteHTTPS(url)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("codecommit::%s://%s", rem.Region, rem.Repository), nil
}

// FromGrc translates a CodeCommit (git-remote-codecommit) GRC URL to a compatible HTTPS URL
// that can be used to fetch and push changes to a CodeCommit repository
func FromGRC(url string) (string, error) {
	rem, err := RemoteGRC(url)
	if err != nil {
		return "", err
	}

	// If a region is not set, check if one is provided through the AWS_REGION environment variable
	if rem.Region == "" {
		if rem.Region = os.Getenv("AWS_REGION"); rem.Region == "" {
			return "", errors.New("no aws region identified")
		}
	}

	domain := "amazonaws.com"
	if rem.Region == "cn-north-1" || rem.Region == "cn-northwest-1" {
		domain = "amazonaws.com.cn"
	}

	return fmt.Sprintf("https://git-codecommit.%s.%s/v1/repos/%s", rem.Region, domain, rem.Repository), nil
}
