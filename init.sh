#!/usr/bin/env zsh

# Stop on error
set -e

HOST_ARCH=$(uname -m)
echo "Host architecture: " $HOST_ARCH

# Get current shell
CURRENT_SHELL=$(basename $SHELL)
echo "Current shell:     " $CURRENT_SHELL

echo "running go run app/bootstrap.go"
go run app/bootstrap.go

