package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func deleteDir(dir string) error {
	log.Infof("Deleting directory: %s", dir)

	// Delete the directory
	return os.RemoveAll(dir)
}

func deleteFile(file string) error {
	log.Infof("Deleting file: %s", file)

	// Delete the file
	return os.Remove(file)
}

func downloadFile(url, dest string) error {
	log.Infof("Downloading file: %s", url)

	// Download the file
	cmd := exec.Command("curl", "-L", "-o", dest, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmd(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isLineInFile(file, line string) (bool, error) {
	// Check if the line is in the file
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if scanner.Text() == line {
			return true, nil
		}
	}

	return false, nil
}

func addIfMissingToFile(file, line string) error {
	if exists, err := isLineInFile(file, line); err != nil {
		log.Fatalf("error: %v", err)
		return err
	} else if exists {
		// Line is already in the file
		log.Infof("Line is already in %s: %s", file, line)
		return nil
	} else {
		// Add the line to the file

		// Open the file
		f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write the line to the file
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}

		log.Infof("Added line to %s: %s", file, line)

		return nil
	}
}

func main() {
	shell := os.Getenv("SHELL")
	home := os.Getenv("HOME")

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
	if !strings.Contains(shell, "zsh") {
		log.Fatalf("Please restart your shell and run the script again in zsh to continue.")
	}

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

	// Install zsh-autosuggestions
	if err := runCmd("sudo", "apt", "install", "zsh-autosuggestions", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install zsh-syntax-highlighting
	if err := runCmd("sudo", "apt", "install", "zsh-syntax-highlighting", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install oh-my-zsh
	omzDir := home + "/.oh-my-zsh"
	if err := deleteDir(omzDir); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := downloadFile("https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh", "install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("zsh", "install.sh"); err != nil {
		log.Errorf("error: %v", err)
		if err := deleteFile("install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Fatal("error installing oh-my-zsh")
	}
	if err := deleteFile("install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Download omz plugins
	pluginsDir := home + "/.oh-my-zsh/custom/plugins"
	if err := runCmd("git", "clone", "https://github.com/zsh-users/zsh-autosuggestions.git", pluginsDir+"/zsh-autosuggestions"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("git", "clone", "https://github.com/zsh-users/zsh-syntax-highlighting.git", pluginsDir+"/zsh-syntax-highlighting"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("git", "clone", "https://github.com/zdharma-continuum/fast-syntax-highlighting.git", pluginsDir+"/fast-syntax-highlighting"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("git", "clone", "--depth=1", "--", "https://github.com/marlonrichert/zsh-autocomplete.git", pluginsDir+"/zsh-autocomplete"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install nodejs
	if err := downloadFile("https://fnm.vercel.app/install", "fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("chmod", "+x", "fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("./fnm-install.sh", "--skip-shell"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := deleteFile("fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Setup .zshrc
	zshrcPath := home + "/.zshrc"
	log.Infof("Setting up %s", zshrcPath)

	// Open ./.zshrc
	f, err := os.OpenFile("./.zshrc", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer f.Close()

	log.Infof("Reading %s", f.Name())

	// Read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// Read the file line by line
		line := scanner.Text()

		log.Infof("Line: %s", line)

		// Add the line to ~/.zshrc
		if err := addIfMissingToFile(zshrcPath, line); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Enable yarn
	if err := runCmd("corepack", "enable", "yarn"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Enable pnpm
	if err := runCmd("corepack", "enable", "pnpm"); err != nil {
		log.Fatalf("error: %v", err)
	}

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
