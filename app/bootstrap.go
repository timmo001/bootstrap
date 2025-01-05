package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func runCmd(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	log.Info("Bootstrapping...")

	// Update apt
	if err := runCmd("sudo", "apt", "update"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Upgrade apt
	if err := runCmd("sudo", "apt", "full-upgrade", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

  // Exit if the shell is not zsh
  shell := os.Getenv("SHELL")
  if strings.Contains(shell, "zsh") == false {
    log.Fatalf("Please restart your shell and run the script again in zsh to continue.")
  }

  //
  // Applications
  //

  // Install curl
  if err := runCmd("sudo", "apt", "install", "curl", "-y"); err != nil {
    log.Fatalf("error: %v", err)
  }

  // Install git
  if err := runCmd("sudo", "apt", "install", "git", "-y"); err != nil {
    log.Fatalf("error: %v", err)
  }

  // Install ruby
  if err := runCmd("sudo", "apt", "install", "ruby", "-y"); err != nil {
    log.Fatalf("error: %v", err)
  }


  //
  // Tools
  //

  // Install markdownlint
  if err := runCmd("sudo", "gem", "install", "mdl"); err != nil {
    log.Fatalf("error: %v", err)
  }

	// Ask if the user is running on a desktop environment
	isDesktop := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Are you running on a desktop environment?").Value(&isDesktop),
		),
	)
	if err := form.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("isDesktop: %v", isDesktop)

	log.Info("Bootstrapping complete.")
}
