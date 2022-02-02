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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDissectHTTPS(t *testing.T) {
	rem, err := DissectHTTPS("https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository")

	require.NoError(t, err)
	assert.Equal(t, "eu-west-1", rem.Region)
	assert.Equal(t, "repository", rem.Repository)
	assert.Equal(t, "", rem.Profile)
}

func TestDissectHTTPS_MalformedURL(t *testing.T) {
	_, err := DissectHTTPS("https://git-codecommit..amazonaws.com/v1/repos/repository")

	assert.Errorf(t, err, "malformed codecommit HTTPS URL")
}

func TestDissectGrc(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		region     string
		repository string
		profile    string
	}{
		{
			name:       "NoNamedProfile",
			url:        "codecommit::eu-west-1://repository",
			region:     "eu-west-1",
			repository: "repository",
			profile:    "",
		},
		{
			name:       "NamedProfile",
			url:        "codecommit::eu-west-1://profile@repository",
			region:     "eu-west-1",
			repository: "repository",
			profile:    "profile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rem, err := DissectGrc(tt.url)

			require.NoError(t, err)
			require.Equal(t, tt.region, rem.Region)
			require.Equal(t, tt.repository, rem.Repository)
			require.Equal(t, tt.profile, rem.Profile)
		})
	}
}

func TestDissectGrc_MalformedURL(t *testing.T) {
	_, err := DissectGrc("codecommit::eu-west-1://")

	assert.Errorf(t, err, "malformed codecommit grc URL")
}
