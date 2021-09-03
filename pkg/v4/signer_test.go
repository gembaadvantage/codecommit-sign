/*
Copyright (c) 2021 Gemba Advantage

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

package v4

import (
	"bytes"
	"fmt"
	"net/http"
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

func TestCanonicalRequest(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	// Construct a GIT request
	req, err := http.NewRequest("GIT", "https://eu-west-1.codecommit/repo", http.NoBody)
	require.NoError(t, err)

	cr := v4.canonicalRequest(req)
	assert.Equal(t, "GIT\n/repo\n\nhost:eu-west-1.codecommit\n\nhost\n", string(cr))
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
	req, err := http.NewRequest("GIT", "https://eu-west-1.codecommit/repo", bytes.NewReader(payload))
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
	req, err := http.NewRequest("GIT", "https://eu-west-1.codecommit/repo?param=1234", http.NoBody)
	require.NoError(t, err)

	cr := v4.canonicalRequest(req)

	queryParams := ""
	assert.True(t, strings.HasPrefix(string(cr), fmt.Sprintf("GIT\n/repo\n%s\n", queryParams)))
}

func TestStringToSign(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	canonicalReq := "GIT\n/repo\n\nhost:eu-west-1.codecommit\n\nhost\n"

	sts := v4.stringToSign([]byte(canonicalReq))
	assert.Equal(t, "AWS4-HMAC-SHA256\n20210901T102523\n20210901/eu-west-1/codecommit/aws4_request\n6f8de131b391656e8b0905c87c7ac2b4efd0c78e8d6aef1abf6e4a642bba0e43", string(sts))
}

func TestSignature(t *testing.T) {
	v4 := Signer{
		service:     "codecommit",
		region:      "eu-west-1",
		credentials: aws.Credentials{},
		requestTime: requestTime,
	}

	stringToSign := "AWS4-HMAC-SHA256\n20210901T102523\n20210901/eu-west-1/codecommit/aws4_request\n6f8de131b391656e8b0905c87c7ac2b4efd0c78e8d6aef1abf6e4a642bba0e43"

	sig := v4.signature([]byte(stringToSign))
	assert.Equal(t, "c315e423bdd21769d77186409a55df2b92cf79869fb5382a207e848b44073aa9", fmt.Sprintf("%x", sig))
}
