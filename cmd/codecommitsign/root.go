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

package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gembaadvantage/codecommit-sign/pkg/awsv4"
	"github.com/gembaadvantage/codecommit-sign/pkg/translate"
	"github.com/spf13/cobra"
)

const (
	desc = `Generate an AWS authenticated V4 signed CodeCommit URL that can be used to fetch
and push changes to a CodeCommit repository within an AWS account.

Both HTTPS and (git-remote-codecommit) GRC URL formats are supported. If a GRC
URL is provided, it will automatically be translated into a compatible HTTPS
URL`

	exs = `Sign a CodeCommit HTTPS URL:

$ codecommit-sign https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository
https://<USERNAME>:<PASSWORD>@git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository

Sign a CodeCommit GRC URL:

$ codecommit-sign codecommit::eu-west-1://repository
https://<USERNAME>:<PASSWORD>@git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository`
)

type signOptions struct {
	Profile  string
	CloneURL string
}

func newRootCmd(out io.Writer, args []string) *cobra.Command {
	opts := signOptions{}

	cmd := &cobra.Command{
		Use:          "codecommit-sign [URL]",
		Short:        "Generate an AWS authenticated V4 signed CodeCommit URL",
		Long:         desc,
		Example:      exs,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CloneURL = args[0]
			return opts.Run(out)
		},
	}

	f := cmd.Flags()
	f.StringVar(&opts.Profile, "profile", "", "the AWS named profile to use when looking up credentials")

	cmd.AddCommand(newVersionCmd(out), newCompletionCmd(out), newManPagesCmd(out))
	return cmd
}

func (o signOptions) Run(out io.Writer) error {
	// Dynamically load options
	opts := []func(*config.LoadOptions) error{}
	if o.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(o.Profile))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		fmt.Fprintln(out, "\u26a0\ufe0f  failed to retrieve default AWS config")
		return err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Fprintln(out, "\u26a0\ufe0f  failed to retrieve AWS credentials")
		return err
	}

	signer := awsv4.NewSigner(creds)
	// Detect if a GRC URL has been provided and translate
	if strings.HasPrefix(o.CloneURL, "codecommit::") {
		var terr error
		o.CloneURL, terr = translate.FromGRC(o.CloneURL)
		if terr != nil {
			return terr
		}
	}

	surl, err := signer.Sign(o.CloneURL)
	if err != nil {
		return err
	}

	fmt.Fprint(out, surl)
	return nil
}
