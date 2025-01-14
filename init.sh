#!/usr/bin/env zsh

# Stop on error
set -e

HOST_ARCH=$(uname -m)
echo "Host architecture: " $HOST_ARCH

# Get current shell
CURRENT_SHELL=$(basename $SHELL)
echo "Current shell:     " $CURRENT_SHELL

GO_ARCH="amd64"
echo "Go architecture:   " $GO_ARCH

GO_VERSION=1.23.4
echo "Go version:        " $GO_VERSION

# If -s flag is passed, skip init
if [[ "$1" == "-s" ]]; then
  echo "Skipping init"
else
  echo "Setup rc files"
  touch ~/.bashrc
  touch ~/.zshrc

  install_go() {
    echo "Removing any existing go installation"
    sudo rm -rf /usr/local/go

    # Install go
    echo "Downloading go" $GO_VERSION " for " $GO_ARCH

    url="https://golang.org/dl/go$GO_VERSION.linux-$GO_ARCH.tar.gz"

    filename="go$GO_VERSION.linux-$GO_ARCH.tar.gz"
    wget $url -O $filename

    echo "Installing go"
    sudo tar -C /usr/local -xzf $filename

    # Check PATH variable
    if [[ ":$PATH:" == *":/usr/local/go/bin:"* ]]; then
      echo "go is already in PATH"
    else
      # Add go to PATH
      echo "Adding go to PATH"
      echo "export PATH=\$PATH:/usr/local/go/bin" >>~/.bashrc
      echo "export PATH=\$PATH:/usr/local/go/bin" >>~/.zshrc
    fi

    # Remove downloaded file
    rm $filename

    # Source current shell rc file
    echo "source ~/.$CURRENT_SHELL"rc""
    source ~/.$CURRENT_SHELL"rc"

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
  set +e
  echo "source ~/.$CURRENT_SHELL"rc""
  source ~/.$CURRENT_SHELL"rc"
  set -e

  if [ $CURRENT_SHELL != "zsh" ]; then
    echo "Install zsh"
    sudo apt update
    sudo apt install -y zsh
  fi

  echo "Init complete"
fi

if [ $CURRENT_SHELL != "zsh" ]; then
  chsh -s $(which zsh)
  sudo chsh -s $(which zsh)

  echo "Please restart your terminal to switch to zsh"
fi

echo "running go mod tidy"
go mod tidy

echo "running go run app/bootstrap.go"
go run app/bootstrap.go

set +e
echo "source ~/.$CURRENT_SHELL"rc""
source ~/.$CURRENT_SHELL"rc"
set -e
