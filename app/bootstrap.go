package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

var forceInstall bool

func init() {
	flag.BoolVar(&forceInstall, "force", false, "Force install all packages")
	flag.Parse()
}

func isExecutableInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

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

func existsDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func runCmdNoInput(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmd(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdInDir(dir, name string, arg ...string) error {
	log.Infof("Running command in directory: %s %s %v", dir, name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
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

func updateOrCloneRepo(repoURL, destDir string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		// Directory does not exist, clone the repo
		return runCmd("git", "clone", "--depth", "1", repoURL, destDir)
	} else {
		// Directory exists, pull the latest changes
		return runCmdInDir(destDir, "git", "pull")
	}
}

func printSeparator(msg ...string) {
	log.Print("================================================================================")
	log.Print("")
	log.Print("")
	log.Printf("== %s ==", strings.Join(msg, " "))
	log.Print("")
	log.Print("")
	log.Print("================================================================================")
}

func main() {
	shell := os.Getenv("SHELL")
	home := os.Getenv("HOME")

	log.Info("Bootstrapping...")

	var installedPackages []string

	// Ask if the user is running on a desktop environment
	printSeparator("Checking if running on a desktop environment")
	isDesktop := false
	isWSL := false
	email := "aidan@timmo.dev"
	name := "Aidan Timson"
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you running on a desktop environment?").
				Value(&isDesktop),
		).Title("Desktop"),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you on WSL? 🤮").
				Value(&isWSL),
		).Title("WSL"),
		huh.NewGroup(
			huh.NewInput().
				Title("What is your email?").
				Value(&email),
			huh.NewInput().
				Title("What is your name?").
				Value(&name),
		).Title("Git config"),
	)
	if err := form.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Update apt
	if err := runCmd("sudo", "apt", "update"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Upgrade apt
	printSeparator("Upgrade apt packages")
	if err := runCmd("sudo", "apt", "full-upgrade", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Cleanup apt
	printSeparator("Cleanup apt packages")
	if err := runCmd("sudo", "apt", "autoremove", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Copy .editorconfig
	printSeparator("Copying .editorconfig")
	if err := runCmd("cp", ".editorconfig", home); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Exit if the shell is not zsh
	printSeparator("Checking shell")
	if !strings.Contains(shell, "zsh") {
		log.Fatalf("Please restart your shell and run the script again in zsh to continue.")
	}

	// Install wget
	printSeparator("wget")
	if forceInstall || !isExecutableInstalled("wget") {
		if err := runCmd("sudo", "apt", "install", "wget", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "wget")
	}

	// Install curl
	printSeparator("curl")
	if forceInstall || !isExecutableInstalled("curl") {
		if err := runCmd("sudo", "apt", "install", "curl", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "curl")
	}

	// Install git
	printSeparator("git")
	if forceInstall || !isExecutableInstalled("git") {
		if err := runCmd("sudo", "apt", "install", "git", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "git")
	}
	err := runCmd("git", "config", "--global", "pull.rebase", "true")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runCmd("git", "config", "--global", "rebase.autoStash", "true")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runCmd("git", "config", "--global", "core.editor", "nvim")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runCmd("git", "config", "--global", "push.default", "current")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runCmd("git", "config", "--global", "user.email", email)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runCmd("git", "config", "--global", "user.name", name)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install gh
	printSeparator("GitHub CLI (gh)")
	if forceInstall || !isExecutableInstalled("gh") {
		if err := runCmd("sudo", "mkdir", "-p", "-m", "775", "/etc/apt/keyrings"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := downloadFile("https://cli.github.com/packages/githubcli-archive-keyring.gpg", "githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "mv", "githubcli-archive-keyring.gpg", "/etc/apt/keyrings/githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "chmod", "go+r", "/etc/apt/keyrings/githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("echo", "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "update"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "install", "gh", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "gh")
	}

	// Install stow
	printSeparator("stow")
	if forceInstall || !isExecutableInstalled("stow") {
		if err := runCmd("sudo", "apt", "install", "stow", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "stow")
	}

	// Setup dotfiles
	printSeparator("Setting up dotfiles")
	dotfilesPath := home + "/.config/dotfiles"
	if err := updateOrCloneRepo("git@github.com:timmo001/dotfiles", dotfilesPath); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmdInDir(dotfilesPath, "./install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ruby
	printSeparator("ruby")
	if forceInstall || !isExecutableInstalled("ruby") {
		if err := runCmd("sudo", "apt", "install", "ruby", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "ruby")
	}

	// Install zsh-autosuggestions
	printSeparator("zsh-autosuggestions")
	if forceInstall || !isExecutableInstalled("zsh-autosuggestions") {
		if err := runCmd("sudo", "apt", "install", "zsh-autosuggestions", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "zsh-autosuggestions")
	}

	// Install zsh-syntax-highlighting
	printSeparator("zsh-syntax-highlighting")
	if forceInstall || !isExecutableInstalled("zsh-syntax-highlighting") {
		if err := runCmd("sudo", "apt", "install", "zsh-syntax-highlighting", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "zsh-syntax-highlighting")
	}

	// Install oh-my-zsh
	printSeparator("oh-my-zsh")
	exists, err := existsDir(home + "/.oh-my-zsh")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if forceInstall || !exists {
		if err := deleteDir(home + "/.oh-my-zsh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := downloadFile("https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh", "omz-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmdNoInput("sh", "omz-install.sh"); err != nil {
			log.Errorf("error: %v", err)
			if err := deleteFile("omz-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			log.Fatal("error installing oh-my-zsh")
		}
		if err := deleteFile("omz-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "oh-my-zsh")
	}

	// Download omz plugins
	printSeparator("Downloading oh-my-zsh plugins")
	pluginsDir := home + "/.oh-my-zsh/custom/plugins"
	if err := updateOrCloneRepo("git@github.com:zsh-users/zsh-autosuggestions.git", pluginsDir+"/zsh-autosuggestions"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := updateOrCloneRepo("git@github.com:zsh-users/zsh-syntax-highlighting.git", pluginsDir+"/zsh-syntax-highlighting"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := updateOrCloneRepo("git@github.com:zdharma-continuum/fast-syntax-highlighting.git", pluginsDir+"/fast-syntax-highlighting"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := updateOrCloneRepo("git@github.com:marlonrichert/zsh-autocomplete.git", pluginsDir+"/zsh-autocomplete"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install starship
	printSeparator("starship")
	if forceInstall || !isExecutableInstalled("starship") {
		if err := runCmd("curl", "-fsSL", "https://starship.rs/install.sh", "-o", "starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("chmod", "+x", "starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("./starship-install.sh", "--yes"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "starship")
	}

	// Install nodejs
	printSeparator("Node.js")
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
	if err := runCmd("fnm", "install", "22"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install python + dependencies
	printSeparator("Python and dependencies")
	if err := runCmd("sudo", "apt", "install", "python3", "python3-dev", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("sudo", "apt", "install", "python3-pip", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("sudo", "apt", "install", "python3-venv", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd(
		"sudo", "apt", "install", "autoconf", "libssl-dev", "libxml2-dev", "libxslt1-dev", "libjpeg-dev", "libffi-dev",
		"libudev-dev", "zlib1g-dev", "pkg-config", "libavformat-dev", "libavcodec-dev", "libavdevice-dev", "libavutil-dev",
		"libswscale-dev", "libswresample-dev", "libavfilter-dev", "ffmpeg", "libgammu-dev", "-y",
	); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install rust
	printSeparator("Rust")
	if forceInstall || !isExecutableInstalled("rustc") {
		if err := downloadFile("https://sh.rustup.rs", "rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("chmod", "+x", "rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("./rustup-init.sh", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "rust")
	}

	// Install zig
	printSeparator("Zig")
	if err := runCmd("sudo", "snap", "install", "zig", "--classic", "--beta"); err != nil {
		log.Fatalf("error: %v", err)
	}
	installedPackages = append(installedPackages, "zig")

	// Install docker
	if !isWSL {
		printSeparator("Docker")
		if forceInstall || !isExecutableInstalled("docker") {
			if err := downloadFile("https://get.docker.com", "docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := runCmd("chmod", "+x", "docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := runCmd("./docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := deleteFile("docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			installedPackages = append(installedPackages, "docker")
		}

		// Install docker compose
		printSeparator("Docker Compose")
		if err := runCmd("sudo", "apt", "install", "docker-compose-plugin", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "docker-compose-plugin")
	}

	// Install homebrew
	printSeparator("Homebrew")
	if forceInstall || !isExecutableInstalled("brew") {
		if err := downloadFile("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", "brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("chmod", "+x", "brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("./brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "homebrew")
	}

	// // Setup .zshrc
	// printSeparator("Setting up .zshrc")
	// zshrcPath := home + "/.zshrc"
	// log.Infof("Setting up %s", zshrcPath)

	// // Open ./.zshrc
	// f, err := os.OpenFile("./.zshrc", os.O_RDONLY, 0644)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// defer f.Close()
	// log.Infof("Reading %s", f.Name())
	// // Read the file line by line using scanner
	// scanner := bufio.NewScanner(f)
	// for scanner.Scan() {
	// 	// Read the file line by line
	// 	line := scanner.Text()

	// 	log.Infof("Line: %s", line)

	// 	// Add the line to ~/.zshrc
	// 	if err := addIfMissingToFile(zshrcPath, line); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// }

	// Enable yarn
	printSeparator("Enabling Yarn")
	if err := runCmd("corepack", "enable", "yarn"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Enable pnpm
	printSeparator("Enabling pnpm")
	if err := runCmd("corepack", "enable", "pnpm"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install markdownlint
	printSeparator("markdownlint")
	if err := runCmd("sudo", "gem", "install", "mdl"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install neovim
	printSeparator("Neovim")
	if err := runCmd("sudo", "apt", "install", "ninja-build", "gettext", "cmake", "unzip", "curl", "build-essential", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := updateOrCloneRepo("git@github.com:neovim/neovim", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmdInDir("neovim", "make", "CMAKE_BUILD_TYPE=Release"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmdInDir("neovim", "sudo", "make", "install"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("npm", "install", "-g", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}
	// if err := deleteDir("neovim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }

	// // Install neovim config
	// printSeparator("Neovim config")
	// if err := deleteDir(home + "/.config/nvim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// if err := updateOrCloneRepo("git@github.com:timmo001/nvim-config", home+"/.config/nvim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	if err := runCmd("nvim", "+PlugInstall", "+qall"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ascii-image-converter
	printSeparator("ascii-image-converter")
	if err := runCmd("go", "install", "github.com/TheZoraiz/ascii-image-converter@latest"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ripgrep
	printSeparator("ripgrep")
	if err := runCmd("sudo", "apt", "install", "ripgrep", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install fzf
	printSeparator("fzf")
	if err := runCmd("sudo", "apt", "install", "fzf", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install bat
	printSeparator("bat")
	if err := runCmd("sudo", "apt", "install", "bat", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("sudo", "ln", "-s", "/usr/bin/batcat", "/usr/bin/bat"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install lynx
	printSeparator("lynx")
	if err := runCmd("sudo", "apt", "install", "lynx", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install lazygit if not installed
	printSeparator("lazygit")
	if forceInstall || !isExecutableInstalled("lazygit") {
		if err := runCmd("brew", "install", "lazygit"); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Install lazydocker if not installed
	printSeparator("lazydocker")
	if forceInstall || !isExecutableInstalled("lazydocker") {
		if err := runCmd("brew", "install", "lazydocker"); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Install nerd fonts
	printSeparator("Nerd Fonts")
	if err := runCmd("sudo", "apt", "install", "fonts-firacode", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("sudo", "apt", "install", "fonts-hack", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := updateOrCloneRepo("https://github.com/ryanoasis/nerd-fonts", "nerd-fonts"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmdInDir("nerd-fonts", "bash", "install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmdInDir("nerd-fonts", "sudo", "bash", "install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := runCmd("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", "'FiraMono Nerd Font Medium 13'"); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("isDesktop: %v", isDesktop)

	// Install desktop environment packages
	if isDesktop {
		// Install gnome-tweaks and gnome-shell-extensions
		printSeparator("gnome-tweaks and gnome-shell-extensions")
		if err := runCmd("sudo", "apt", "install", "gnome-tweaks", "gnome-shell-extensions", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Setup flatpak and flathub
		printSeparator("Setting up flatpak and flathub")
		if err := runCmd("sudo", "apt", "install", "flatpak", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "install", "gnome-software-plugin-flatpak", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install zen browser
		printSeparator("Zen Browser")
		if err := runCmd("flatpak", "install", "flathub", "io.github.zen_browser.zen", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install vs*ode
		if forceInstall || !isExecutableInstalled("code") {
			printSeparator("VS C*de")
			if err := downloadFile("https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64", "vscode.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := runCmd("sudo", "apt", "install", "./vscode.deb", "-y"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := deleteFile("vscode.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
		}

		// Install postman
		printSeparator("Postman")
		if err := downloadFile("https://dl.pstmn.io/download/latest/linux_64", "postman.tar.gz"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "rm", "-rf", "/usr/bin/postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "rm", "-rf", "/opt/Postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "tar", "-xzf", "postman.tar.gz", "-C", "/opt"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "ln", "-s", "/opt/Postman/Postman", "/usr/bin/postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("postman.tar.gz"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install ghostty
		printSeparator("Ghostty")
		if err := runCmd("sudo", "apt", "install", "libgtk-4-dev", "libadwaita-1-dev", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "libgtk-4-dev", "libadwaita-1-dev")
		if err := updateOrCloneRepo("https://github.com/ghostty-org/ghostty", "ghostty"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmdInDir("ghostty", "sudo", "zig", "build", "-p", "/usr", "-Doptimize=ReleaseFast"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "ghostty")
		// Set CTRL+ALT+T to open ghostty
		if err := runCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "name", "'Open Ghostty'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "binding", "'<Primary><Alt>t'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "command", "'/usr/bin/ghostty'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		// // Scan the ~/.config/ghostty/config and for differences with ghostty-config
		// ghosttyConfigPath := home + "/.config/ghostty/config"
		// // Open ./ghostty-config
		// f, err := os.OpenFile("./ghostty-config", os.O_RDONLY, 0644)
		// if err != nil {
		// 	log.Fatalf("error: %v", err)
		// }
		// defer f.Close()
		// log.Infof("Reading %s", f.Name())
		// // Read the file line by line using scanner
		// scanner := bufio.NewScanner(f)
		// for scanner.Scan() {
		// 	// Read the file line by line
		// 	line := scanner.Text()

		// 	log.Infof("Line: %s", line)

		// 	// Add the line to ~/.zshrc
		// 	if err := addIfMissingToFile(ghosttyConfigPath, line); err != nil {
		// 		log.Fatalf("error: %v", err)
		// 	}
		// }

		// Install chrome
		printSeparator("Google Chrome")
		if err := downloadFile("https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb", "chrome.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "install", "./chrome.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("chrome.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install slack
		printSeparator("Slack")
		if err := runCmd("sudo", "snap", "install", "slack", "--classic"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install discord
		printSeparator("Discord")
		if err := downloadFile("https://discord.com/api/download?platform=linux&format=deb", "discord.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "install", "./discord.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install steam
		printSeparator("Steam")
		if forceInstall || !isExecutableInstalled("steam") {
			if err := downloadFile("https://cdn.fastly.steamstatic.com/client/installer/steam.deb", "steam.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := runCmd("sudo", "apt", "install", "./steam.deb", "-y"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := deleteFile("steam.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
		}

		// Install sunshine
		printSeparator("Sunshine")
		if err := downloadFile("https://github.com/LizardByte/Sunshine/releases/download/v0.23.1/sunshine-ubuntu-24.04-amd64.deb", "sunshine.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := runCmd("sudo", "apt", "install", "./sunshine.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := deleteFile("sunshine.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install moonlight
		printSeparator("Moonlight")
		if err := runCmd("flatpak", "install", "flathub", "com.moonlight_stream.Moonlight", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// // Check for GoXLR and install GoXLR-Utility
		// if _, err := os.Stat("/dev/snd/GoXLR"); err == nil {
		// 	if forceInstall || !isExecutableInstalled("goxlr-utility") {
		// 		printSeparator("GoXLR-Utility")
		// 		if err := downloadFile("https://github.com/GoXLR-on-Linux/goxlr-utility/releases/download/v1.1.4/goxlr-utility_1.1.4-1_amd64.deb", "goxlr-utility.deb"); err != nil {
		// 			log.Fatalf("error: %v", err)
		// 		}
		// 		if err := runCmd("sudo", "apt", "install", "./goxlr-utility.deb", "-y"); err != nil {
		// 			log.Fatalf("error: %v", err)
		// 		}
		// 		if err := deleteFile("goxlr-utility.deb"); err != nil {
		// 			log.Fatalf("error: %v", err)
		// 		}
		// 	}
		// }

		// Install hyprland
		printSeparator("Hyprland")
		if forceInstall || !isExecutableInstalled("hyprland") {
			if err := runCmd("sudo", "apt", "install", "hyprland", "hyprland-backgrounds", "wofi", "wofi-pass", "wl-clipboard", "pseudo", "libgtk-4-dev", "waybar", "fonts-font-awesome", "clang-tidy", "gobject-introspection", "libdbusmenu-gtk3-dev", "libevdev-dev", "libfmt-dev", "libgirepository1.0-dev", "libgtk-3-dev", "libgtkmm-3.0-dev", "libinput-dev", "libjsoncpp-dev", "libmpdclient-dev", "libnl-3-dev", "libnl-genl-3-dev", "libpulse-dev", "libsigc++-2.0-dev", "libspdlog-dev", "libwayland-dev", "scdoc", "upower", "libxkbregistry-dev", "-y"); err != nil {
				log.Fatalf("error: %v", err)
			}
		}
	}

	log.Info("Bootstrapping complete.")
	log.Infof("Installed packages: %v", installedPackages)
}
