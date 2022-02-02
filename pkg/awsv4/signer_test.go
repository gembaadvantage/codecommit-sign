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

package awsv4

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Use a static time to ensure consistent test results
	requestTime = time.Date(2021, 9, 1, 10, 25, 23, 0, time.UTC)
)

const (
	repoURL = "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/dummy-repo"
)

func TestSign(t *testing.T) {
	creds := aws.Credentials{
		AccessKeyID:     "ACCESS_KEY_ID",
		SecretAccessKey: "SECRET_ACCESS_KEY",
		SessionToken:    "SESSION_TOKEN",
	}

	s := NewSigner(creds)

	req, err := s.Sign(repoURL)
	require.NoError(t, err)

	// Check all parts of the signed URL
	u, err := url.Parse(req)
	require.NoError(t, err)

	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, creds.AccessKeyID+"%"+creds.SessionToken, u.User.Username())

	// Test password using a regex
	passw, _ := u.User.Password()
	rgx := regexp.MustCompile(fmt.Sprintf("^%sZ[a-zA-Z0-9]{64}$", s.requestTime.Format("20060102T150405")))
	assert.True(t, rgx.MatchString(passw))

	assert.Equal(t, "git-codecommit.eu-west-1.amazonaws.com", u.Host)
	assert.Equal(t, "/v1/repos/dummy-repo", u.Path)
}

func TestSign_MalformedUrl(t *testing.T) {
	s := NewSigner(aws.Credentials{})

	sig, err := s.Sign("https://codecommit.amazonaws.com")

	require.Error(t, err)
	assert.Empty(t, sig)
}

func TestIdentifyRegion(t *testing.T) {
	rgn, err := identifyRegion(repoURL)

	require.NoError(t, err)
	assert.Equal(t, "eu-west-1", rgn)
}

func TestIdentifyRegion_MalformedUrl(t *testing.T) {
	rgn, err := identifyRegion("https://codecommit.amazonaws.com")

	require.Error(t, err)
	assert.EqualError(t, err, "no region found in malformed codecommit URL")
	assert.Empty(t, rgn)
}

func TestCanonicalRequest(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	// Construct a GIT request
	req, err := http.NewRequest("GIT", repoURL, http.NoBody)
	require.NoError(t, err)

	cr := v4.canonicalRequest(req)
	assert.Equal(t, "GIT\n/v1/repos/dummy-repo\n\nhost:git-codecommit.eu-west-1.amazonaws.com\n\nhost\n", string(cr))
}

func TestCanonicalRequest_IgnoresPayload(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	payload := []byte("payload")

	// Construct a GIT request
	req, err := http.NewRequest("GIT", repoURL, bytes.NewReader(payload))
	require.NoError(t, err)

	cr := v4.canonicalRequest(req)

	// Payload would be a hex-encoded string after the final newline
	assert.True(t, strings.HasSuffix(string(cr), "\n"))
}

func TestCanonicalRequest_IgnoresQueryParameters(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	// Construct a GIT request
	req, err := http.NewRequest("GIT", repoURL+"?param=12345", http.NoBody)
	require.NoError(t, err)

	cr := v4.canonicalRequest(req)

	queryParams := ""
	assert.True(t, strings.HasPrefix(string(cr), fmt.Sprintf("GIT\n/v1/repos/dummy-repo\n%s\n", queryParams)))
}

func TestStringToSign(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	canonicalReq := "GIT\n/v1/repos/dummy-repo\n\nhost:git-codecommit.eu-west-1.amazonaws.com\n\nhost\n"

	sts := v4.stringToSign([]byte(canonicalReq))
	assert.Equal(t, "AWS4-HMAC-SHA256\n20210901T102523\n20210901/eu-west-1/codecommit/aws4_request\nb7cad41c14b37f02e4d2deaf4f0773423b7dbe5db34af6b45362223291f968ef", string(sts))
}

func TestSignature(t *testing.T) {
	v4 := Signer{
		service: "codecommit",
		region:  "eu-west-1",
		credentials: aws.Credentials{
			SecretAccessKey: "SECRET_ACCESS_KEY",
		},
		requestTime: requestTime,
	}

	stringToSign := "AWS4-HMAC-SHA256\n20210901T102523\n20210901/eu-west-1/codecommit/aws4_request\nb7cad41c14b37f02e4d2deaf4f0773423b7dbe5db34af6b45362223291f968ef"

	sig := v4.signature([]byte(stringToSign))
	assert.Equal(t, "670ef0c32f13fc847ed0c001a2b90bc772451c1eb89eadf84a73f3e82da9f56a", fmt.Sprintf("%x", sig))
}
