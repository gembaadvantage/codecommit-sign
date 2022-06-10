#!/bin/sh

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/completions.sh
set -e
rm -rf completions
mkdir completions

# Directly invoke codecommit-sign and generate the shell completion scripts
for SH in bash zsh fish; do
	go run ./cmd/codecommit-sign/... completion "${SH}" > "completions/codecommit-sign.${SH}"
done