package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func main() {
	log.Info("Bootstrapping...")

	isDesktop := false

	// Ask if the user is running on a desktop environment
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Are you running on a desktop environment?").Value(&isDesktop),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("isDesktop: %v", isDesktop)

	// Update apt
	if err := RunCmd("sudo", "apt", "update"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Upgrade apt
	if err := RunCmd("sudo", "apt", "full-upgrade"); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Info("Bootstrapping complete.")
}
