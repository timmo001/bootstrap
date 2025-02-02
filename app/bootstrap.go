package main

import (
	"flag"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"

	u "github.com/timmo001/bootstrap/utils"
)

var forceInstall bool

func init() {
	flag.BoolVar(&forceInstall, "force", false, "Force install all packages")
	flag.Parse()
}

func main() {
	shell := os.Getenv("SHELL")
	home := os.Getenv("HOME")

	log.Info("Bootstrapping...")

	var installedPackages []string

	// Ask if the user is running on a desktop environment
	u.PrintSeparator("Checking if running on a desktop environment")
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
	if err := u.RunCmd("sudo", "apt", "update"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Upgrade apt
	u.PrintSeparator("Upgrade apt packages")
	if err := u.RunCmd("sudo", "apt", "full-upgrade", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Cleanup apt
	u.PrintSeparator("Cleanup apt packages")
	if err := u.RunCmd("sudo", "apt", "autoremove", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Copy .editorconfig
	u.PrintSeparator("Copying .editorconfig")
	if err := u.RunCmd("cp", ".editorconfig", home); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Exit if the shell is not zsh
	u.PrintSeparator("Checking shell")
	if !strings.Contains(shell, "zsh") {
		log.Fatalf("Please restart your shell and run the script again in zsh to continue.")
	}

	// Install wget
	u.PrintSeparator("wget")
	if forceInstall || !u.IsExecutableInstalled("wget") {
		if err := u.RunCmd("sudo", "apt", "install", "wget", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "wget")
	}

	// Install curl
	u.PrintSeparator("curl")
	if forceInstall || !u.IsExecutableInstalled("curl") {
		if err := u.RunCmd("sudo", "apt", "install", "curl", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "curl")
	}

	// Install git
	u.PrintSeparator("git")
	if forceInstall || !u.IsExecutableInstalled("git") {
		if err := u.RunCmd("sudo", "apt", "install", "git", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "git")
	}
	err := u.RunCmd("git", "config", "--global", "pull.rebase", "true")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = u.RunCmd("git", "config", "--global", "rebase.autoStash", "true")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = u.RunCmd("git", "config", "--global", "core.editor", "nvim")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = u.RunCmd("git", "config", "--global", "push.default", "current")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = u.RunCmd("git", "config", "--global", "user.email", email)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = u.RunCmd("git", "config", "--global", "user.name", name)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install gh
	u.PrintSeparator("GitHub CLI (gh)")
	if forceInstall || !u.IsExecutableInstalled("gh") {
		if err := u.RunCmd("sudo", "mkdir", "-p", "-m", "775", "/etc/apt/keyrings"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DownloadFile("https://cli.github.com/packages/githubcli-archive-keyring.gpg", "githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "mv", "githubcli-archive-keyring.gpg", "/etc/apt/keyrings/githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "chmod", "go+r", "/etc/apt/keyrings/githubcli-archive-keyring.gpg"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("echo", "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "update"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "gh", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "gh")
	}

	// Install stow
	u.PrintSeparator("stow")
	if forceInstall || !u.IsExecutableInstalled("stow") {
		if err := u.RunCmd("sudo", "apt", "install", "stow", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "stow")
	}

	// Setup dotfiles
	u.PrintSeparator("Setting up dotfiles")
	dotfilesPath := home + "/.config/dotfiles"
	if err := u.UpdateOrCloneRepo("git@github.com:timmo001/dotfiles", dotfilesPath); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir(dotfilesPath, "./install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ruby
	u.PrintSeparator("ruby")
	if forceInstall || !u.IsExecutableInstalled("ruby") {
		if err := u.RunCmd("sudo", "apt", "install", "ruby", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "ruby")
	}

	// Install zsh-autosuggestions
	u.PrintSeparator("zsh-autosuggestions")
	if forceInstall || !u.IsExecutableInstalled("zsh-autosuggestions") {
		if err := u.RunCmd("sudo", "apt", "install", "zsh-autosuggestions", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "zsh-autosuggestions")
	}

	// Install zsh-syntax-highlighting
	u.PrintSeparator("zsh-syntax-highlighting")
	if forceInstall || !u.IsExecutableInstalled("zsh-syntax-highlighting") {
		if err := u.RunCmd("sudo", "apt", "install", "zsh-syntax-highlighting", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "zsh-syntax-highlighting")
	}

	// Install oh-my-zsh
	u.PrintSeparator("oh-my-zsh")
	exists, err := u.ExistsDir(home + "/.oh-my-zsh")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if forceInstall || !exists {
		if err := u.DeleteDir(home + "/.oh-my-zsh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DownloadFile("https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh", "omz-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmdNoInput("sh", "omz-install.sh"); err != nil {
			log.Errorf("error: %v", err)
			if err := u.DeleteFile("omz-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			log.Fatal("error installing oh-my-zsh")
		}
		if err := u.DeleteFile("omz-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "oh-my-zsh")
	}

	// Download omz plugins
	u.PrintSeparator("Downloading oh-my-zsh plugins")
	pluginsDir := home + "/.oh-my-zsh/custom/plugins"
	if err := u.UpdateOrCloneRepo("git@github.com:zsh-users/zsh-autosuggestions.git", pluginsDir+"/zsh-autosuggestions"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("git@github.com:zsh-users/zsh-syntax-highlighting.git", pluginsDir+"/zsh-syntax-highlighting"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("git@github.com:zdharma-continuum/fast-syntax-highlighting.git", pluginsDir+"/fast-syntax-highlighting"); err != nil {
		log.Errorf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("git@github.com:marlonrichert/zsh-autocomplete.git", pluginsDir+"/zsh-autocomplete"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install starship
	u.PrintSeparator("starship")
	if forceInstall || !u.IsExecutableInstalled("starship") {
		if err := u.RunCmd("curl", "-fsSL", "https://starship.rs/install.sh", "-o", "starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("chmod", "+x", "starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("./starship-install.sh", "--yes"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("starship-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "starship")
	}

	// Install nodejs
	u.PrintSeparator("Node.js")
	if err := u.DownloadFile("https://fnm.vercel.app/install", "fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("chmod", "+x", "fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("./fnm-install.sh", "--skip-shell"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.DeleteFile("fnm-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("fnm", "install", "22"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install python + dependencies
	u.PrintSeparator("Python and dependencies")
	if err := u.RunCmd("sudo", "apt", "install", "python3", "python3-dev", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "apt", "install", "python3-pip", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "apt", "install", "python3-venv", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd(
		"sudo", "apt", "install", "autoconf", "libssl-dev", "libxml2-dev", "libxslt1-dev", "libjpeg-dev", "libffi-dev",
		"libudev-dev", "zlib1g-dev", "pkg-config", "libavformat-dev", "libavcodec-dev", "libavdevice-dev", "libavutil-dev",
		"libswscale-dev", "libswresample-dev", "libavfilter-dev", "ffmpeg", "libgammu-dev", "-y",
	); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install rust
	u.PrintSeparator("Rust")
	if forceInstall || !u.IsExecutableInstalled("rustc") {
		if err := u.DownloadFile("https://sh.rustup.rs", "rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("chmod", "+x", "rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("./rustup-init.sh", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("rustup-init.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "rust")
	}

	// Install zig
	u.PrintSeparator("Zig")
	if err := u.RunCmd("sudo", "snap", "install", "zig", "--classic", "--beta"); err != nil {
		log.Fatalf("error: %v", err)
	}
	installedPackages = append(installedPackages, "zig")

	// Install docker
	if !isWSL {
		u.PrintSeparator("Docker")
		if forceInstall || !u.IsExecutableInstalled("docker") {
			if err := u.DownloadFile("https://get.docker.com", "docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.RunCmd("chmod", "+x", "docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.RunCmd("./docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.DeleteFile("docker-install.sh"); err != nil {
				log.Fatalf("error: %v", err)
			}
			installedPackages = append(installedPackages, "docker")
		}

		// Install docker compose
		u.PrintSeparator("Docker Compose")
		if err := u.RunCmd("sudo", "apt", "install", "docker-compose-plugin", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "docker-compose-plugin")
	}

	// Install homebrew
	u.PrintSeparator("Homebrew")
	if forceInstall || !u.IsExecutableInstalled("brew") {
		if err := u.DownloadFile("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", "brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("chmod", "+x", "brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("./brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("brew-install.sh"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "homebrew")
	}

	// // Setup .zshrc
	// u.PrintSeparator("Setting up .zshrc")
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
	u.PrintSeparator("Enabling Yarn")
	if err := u.RunCmd("corepack", "enable", "yarn"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Enable pnpm
	u.PrintSeparator("Enabling pnpm")
	if err := u.RunCmd("corepack", "enable", "pnpm"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install markdownlint
	u.PrintSeparator("markdownlint")
	if err := u.RunCmd("sudo", "gem", "install", "mdl"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install neovim
	u.PrintSeparator("Neovim")
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
	if err := u.RunCmd("npm", "install", "-g", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}
	// if err := deleteDir("neovim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }

	// // Install neovim config
	// u.PrintSeparator("Neovim config")
	// if err := deleteDir(home + "/.config/nvim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// if err := updateOrCloneRepo("git@github.com:timmo001/nvim-config", home+"/.config/nvim"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// if err := u.RunCmd("nvim", "+PlugInstall", "+qall"); err != nil {
	// 	log.Fatalf("error: %v", err)
	// }

	// Install ascii-image-converter
	u.PrintSeparator("ascii-image-converter")
	if err := u.RunCmd("go", "install", "github.com/TheZoraiz/ascii-image-converter@latest"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ripgrep
	u.PrintSeparator("ripgrep")
	if err := u.RunCmd("sudo", "apt", "install", "ripgrep", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install fzf
	u.PrintSeparator("fzf")
	if err := u.RunCmd("sudo", "apt", "install", "fzf", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install bat
	u.PrintSeparator("bat")
	if err := u.RunCmd("sudo", "apt", "install", "bat", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "ln", "-s", "/usr/bin/batcat", "/usr/bin/bat"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install lynx
	u.PrintSeparator("lynx")
	if err := u.RunCmd("sudo", "apt", "install", "lynx", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install lazygit if not installed
	u.PrintSeparator("lazygit")
	if forceInstall || !u.IsExecutableInstalled("lazygit") {
		if err := u.RunCmd("brew", "install", "lazygit"); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Install lazydocker if not installed
	u.PrintSeparator("lazydocker")
	if forceInstall || !u.IsExecutableInstalled("lazydocker") {
		if err := u.RunCmd("brew", "install", "lazydocker"); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Install nerd fonts
	u.PrintSeparator("Nerd Fonts")
	if err := u.RunCmd("sudo", "apt", "install", "fonts-firacode", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "apt", "install", "fonts-hack", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("https://github.com/ryanoasis/nerd-fonts", "nerd-fonts"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir("nerd-fonts", "bash", "install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir("nerd-fonts", "sudo", "bash", "install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", "'FiraMono Nerd Font Medium 13'"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Bun
	u.PrintSeparator("Bun")
	if err := u.DownloadFile("https://bun.sh/install", "bun-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("chmod", "+x", "bun-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("./bun-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.DeleteFile("bun-install.sh"); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("isDesktop: %v", isDesktop)

	// Install desktop environment packages
	if isDesktop {
		// Install gnome-tweaks and gnome-shell-extensions
		u.PrintSeparator("gnome-tweaks and gnome-shell-extensions")
		if err := u.RunCmd("sudo", "apt", "install", "gnome-tweaks", "gnome-shell-extensions", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Setup flatpak and flathub
		u.PrintSeparator("Setting up flatpak and flathub")
		if err := u.RunCmd("sudo", "apt", "install", "flatpak", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "gnome-software-plugin-flatpak", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install zen browser
		u.PrintSeparator("Zen Browser")
		if err := u.RunCmd("flatpak", "install", "flathub", "io.github.zen_browser.zen", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install vs*ode
		if forceInstall || !u.IsExecutableInstalled("code") {
			u.PrintSeparator("VS C*de")
			if err := u.DownloadFile("https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64", "vscode.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.RunCmd("sudo", "apt", "install", "./vscode.deb", "-y"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.DeleteFile("vscode.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
		}

		// Install postman
		u.PrintSeparator("Postman")
		if err := u.DownloadFile("https://dl.pstmn.io/download/latest/linux_64", "postman.tar.gz"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "rm", "-rf", "/usr/bin/postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "rm", "-rf", "/opt/Postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "tar", "-xzf", "postman.tar.gz", "-C", "/opt"); err != nil {

		}
		if err := u.RunCmd("sudo", "ln", "-s", "/opt/Postman/Postman", "/usr/bin/postman"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("postman.tar.gz"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install ghostty
		u.PrintSeparator("Ghostty")
		if err := u.RunCmd("sudo", "apt", "install", "libgtk-4-dev", "libadwaita-1-dev", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "libgtk-4-dev", "libadwaita-1-dev")
		if err := u.UpdateOrCloneRepo("https://github.com/ghostty-org/ghostty", "ghostty"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmdInDir("ghostty", "sudo", "zig", "build", "-p", "/usr", "-Doptimize=ReleaseFast"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "ghostty")
		// Set CTRL+ALT+T to open ghostty
		if err := u.RunCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "name", "'Open Ghostty'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "binding", "'<Primary><Alt>t'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("gsettings", "set", "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/", "command", "'/usr/bin/ghostty'"); err != nil {
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
		u.PrintSeparator("Google Chrome")
		if err := u.DownloadFile("https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb", "chrome.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "./chrome.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("chrome.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install slack
		u.PrintSeparator("Slack")
		if err := u.RunCmd("sudo", "snap", "install", "slack", "--classic"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install discord
		u.PrintSeparator("Discord")
		if err := u.DownloadFile("https://discord.com/api/download?platform=linux&format=deb", "discord.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "./discord.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install steam
		u.PrintSeparator("Steam")
		if forceInstall || !u.IsExecutableInstalled("steam") {
			if err := u.DownloadFile("https://cdn.fastly.steamstatic.com/client/installer/steam.deb", "steam.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.RunCmd("sudo", "apt", "install", "./steam.deb", "-y"); err != nil {
				log.Fatalf("error: %v", err)
			}
			if err := u.DeleteFile("steam.deb"); err != nil {
				log.Fatalf("error: %v", err)
			}
		}

		// Install sunshine
		u.PrintSeparator("Sunshine")
		if err := u.DownloadFile("https://github.com/LizardByte/Sunshine/releases/download/v0.23.1/sunshine-ubuntu-24.04-amd64.deb", "sunshine.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "./sunshine.deb", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.DeleteFile("sunshine.deb"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install moonlight
		u.PrintSeparator("Moonlight")
		if err := u.RunCmd("flatpak", "install", "flathub", "com.moonlight_stream.Moonlight", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// // Check for GoXLR and install GoXLR-Utility
		// if _, err := os.Stat("/dev/snd/GoXLR"); err == nil {
		// 	if forceInstall || !isExecutableInstalled("goxlr-utility") {
		// 		u.PrintSeparator("GoXLR-Utility")
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
		u.PrintSeparator("Hyprland")
		if err := u.RunCmd("sudo", "apt", "install", "hyprland", "hyprland-backgrounds", "wofi", "wofi-pass", "wl-clipboard", "pseudo", "libgtk-4-dev", "waybar", "fonts-font-awesome", "clang-tidy", "gobject-introspection", "libdbusmenu-gtk3-dev", "libevdev-dev", "libfmt-dev", "libgirepository1.0-dev", "libgtk-3-dev", "libgtkmm-3.0-dev", "libinput-dev", "libjsoncpp-dev", "libmpdclient-dev", "libnl-3-dev", "libnl-genl-3-dev", "libpulse-dev", "libsigc++-2.0-dev", "libspdlog-dev", "libwayland-dev", "scdoc", "upower", "libxkbregistry-dev", "sway-notification-center", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// // Install flameshot
		// u.PrintSeparator("Flameshot")
		// if forceInstall || !isExecutableInstalled("flameshot") {
		// 	if err := runCmd("sudo", "apt", "install", "flameshot", "-y"); err != nil {
		// 		log.Fatalf("error: %v", err)
		// 	}
		// }

		// Install grimblast
		u.PrintSeparator("Grimblast")
		if err := u.UpdateOrCloneRepo("git@github.com:hyprwm/contrib", "hyprwm-contrib"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmdInDir("hyprwm-contrib/grimblast", "sudo", "make", "install"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "apt", "install", "grim", "slurp", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install swaybg
		u.PrintSeparator("swaybg")
		if err := u.RunCmd("sudo", "apt", "install", "swaybg", "-y"); err != nil {
			log.Fatalf("error: %v", err)
		}

		// Install hyprlock
		// u.PrintSeparator("Hyprlock")
		// if err := downloadFile("https://github.com/JaKooLit/Ubuntu-Hyprland/raw/refs/heads/24.10/install-scripts/hyprlock.sh", "hyprlock.sh"); err != nil {
		// 	log.Fatalf("error: %v", err)
		// }
		// if err := runCmd("chmod", "+x", "hyprlock.sh"); err != nil {
		// 	log.Fatalf("error: %v", err)
		// }
		// if err := runCmd("./hyprlock.sh"); err != nil {
		// 	log.Fatalf("error: %v", err)
		// }

	}

	log.Info("Bootstrapping complete.")
	log.Infof("Installed packages: %v", installedPackages)
}
