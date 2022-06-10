#!/bin/sh

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/manpages.sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/codecommitsign/... man | gzip -c -9 > manpages/codecommit-sign.1.gz