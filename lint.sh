#!/usr/bin/env bash
set -e

# TODO: Use docker for reproducible environment

cd "$(dirname "${BASH_SOURCE[0]}")"

echoerr() { echo "$@" 1>&2; }

set +e
command -v golangci-lint > /dev/null
if [[ $? != 0 ]]; then
    echoerr "Cannot find golangci-lint. Please make sure it is installed."
    exit 1
fi
set -e

golangci-lint run
