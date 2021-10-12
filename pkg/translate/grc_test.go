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

package translate

import (
	"testing"
)

func TestToGrc(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		err      string
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
		{
			name:     "NoRegion",
			url:      "https://git-codecommit..amazonaws.com/v1/repos/repository",
			expected: "",
			err:      "malformed codecommit HTTPS URL",
		},
		{
			name:     "NoRepositoryName",
			url:      "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/",
			expected: "",
			err:      "malformed codecommit HTTPS URL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ToGrc(tt.url)

			if err != nil {
				if tt.err == "" {
					t.Fatalf("unexpected error '%s'", err.Error())
				}

				if err.Error() != tt.err {
					t.Fatalf("expected error '%s' but received error '%s'\n", tt.err, err)
				}
			}

			if actual != tt.expected {
				t.Fatalf("expected %s but received %s\n", tt.expected, actual)
			}
		})
	}
}

func TestFromGrc(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		err      string
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
		{
			name:     "NoRegion",
			url:      "codecommit::://repository",
			expected: "",
			err:      "malformed codecommit grc URL",
		},
		{
			name:     "NoRepositoryName",
			url:      "codecommit::eu-west-1://",
			expected: "",
			err:      "malformed codecommit grc URL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FromGrc(tt.url)

			if err != nil {
				if err.Error() != tt.err {
					t.Fatalf("expected %s but received %s\n", tt.err, err)
				}
			}

			if actual != tt.expected {
				t.Fatalf("expected %s but received %s\n", tt.expected, actual)
			}
		})
	}
}