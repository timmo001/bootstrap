package main

import (
	"github.com/charmbracelet/log"

	u "github.com/timmo001/bootstrap/utils"
)

func main() {
	u.PrintSeparator("Install or update neovim")

	if err := u.RunCmd("sudo", "apt", "install", "ninja-build", "gettext", "cmake", "unzip", "curl", "build-essential", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("git@github.com:neovim/neovim", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir("neovim", "make", "CMAKE_BUILD_TYPE=Release"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir("neovim", "sudo", "make", "install"); err != nil {
		log.Fatalf("error: %v", err)
	}
}
