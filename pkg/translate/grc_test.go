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

func TestToGRC(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "ValidNoAuthentication",
			url:      "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository",
			expected: "codecommit::eu-west-1://repository",
		},
		{
			name:     "ValidAuthentication",
			url:      "https://username:password@git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository",
			expected: "codecommit::eu-west-1://repository",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ToGRC(tt.url)

			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestToGRC_MalformedURL(t *testing.T) {
	url, err := ToGRC("https://git-codecommit..amazonaws.com/v1/repos/repository")

	require.Error(t, err)
	assert.Equal(t, "", url)
}

func TestFromGRC(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "NoNamedProfile",
			url:      "codecommit::eu-west-1://repository",
			expected: "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository",
		},
		{
			name:     "IgnoresNamedProfile",
			url:      "codecommit::eu-west-1://profile@repository",
			expected: "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FromGRC(tt.url)

			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestFromGrc_MalformedURL(t *testing.T) {
	url, err := FromGRC("codecommit::eu-west-1://")

	require.Error(t, err)
	assert.Equal(t, "", url)
}
