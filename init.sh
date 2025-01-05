#!/bin/bash

# Stop on error
set -e

GO_VERSION=1.23.4
HOST_ARCH=$(uname -m)

# Get current shell
#   use last part of SHELL variable so that it works for both /bin/bash and /usr/bin/zsh etc.
CURRENT_SHELL=$(echo $SHELL | awk -F/ '{print $NF}')
echo "Current shell:     " $CURRENT_SHELL

# Check if host is x86_64
if [ $HOST_ARCH == "x86_64" ]; then
  HOST_ARCH="amd64"
fi
echo "Host architecture: " $HOST_ARCH

# If -s flag is passed, skip init
if [ "$1" == "-s" ]; then
  echo "Skipping init"
else
  echo "Setup rc files"
  touch ~/.bashrc
  touch ~/.zshrc

  install_go() {
    echo "Removing any existing go installation"
    sudo rm -rf /usr/local/go

    # Install go
    echo "Downloading go" $GO_VERSION " for " $HOST_ARCH

    url="https://golang.org/dl/go$GO_VERSION.linux-$HOST_ARCH.tar.gz"

    filename="go$GO_VERSION.linux-$HOST_ARCH.tar.gz"
    wget $url -O $filename

    echo "Installing go"
    sudo tar -C /usr/local -xzf $filename

    # Check PATH variable
    if [[ ":$PATH:" == *":/usr/local/go/bin:"* ]]; then
      echo "go is already in PATH"
    else
      # Add go to PATH
      echo "Adding go to PATH"
      echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
      echo "export PATH=$PATH:/usr/local/go/bin" >>~/.zshrc
    fi

    # Remove downloaded file
    rm $filename

    echo "Verify go installation"
    go version
  }

  # Check go is installed
  if ! command -v go &>/dev/null; then
    echo "go could not be found"

    install_go
  else
    echo "go is already installed"
    go version

    # Check if go version is correct
    #  if not, reinstall go
    if [ $(go version | awk '{print $3}') != "go$GO_VERSION" ]; then
      echo "go version is not correct"
      install_go
    fi
  fi

  # Source current shell rc file
  echo "source ~/.$CURRENT_SHELL"rc""
  source ~/.$CURRENT_SHELL"rc"

  echo "Init complete"
fi

echo "running go mod tidy"
go mod tidy

echo "running go run app/bootstrap.go"
go run app/bootstrap.go
