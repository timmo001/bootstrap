package main

import (
	"flag"
	"os"

	"github.com/charmbracelet/log"

	u "github.com/timmo001/bootstrap/utils"
)

var forceInstall bool

func init() {
	flag.BoolVar(&forceInstall, "force", false, "Force install all packages")
	flag.Parse()
}

func main() {
	home := os.Getenv("HOME")

	log.Info("Bootstrapping...")

	var installedPackages []string

	// Ask if the user is running on a desktop environment
	u.PrintSeparator("Checking if running on a desktop environment")
	isDesktop := true
	email := "aidan@timmo.dev"
	name := "Aidan Timson"

	// Update system: Use pacman update for Arch
	u.PrintSeparator("Updating system")
	if err := u.RunCmd("sudo", "pacman", "-Syu", "--noconfirm"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Copy .editorconfig
	u.PrintSeparator("Copying .editorconfig")
	if err := u.RunCmd("cp", ".editorconfig", home); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install wget
	u.PrintSeparator("wget")
	if forceInstall || !u.IsExecutableInstalled("wget") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "wget"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "wget")
	}

	// Install curl
	u.PrintSeparator("curl")
	if forceInstall || !u.IsExecutableInstalled("curl") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "curl"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "curl")
	}

	// Install git
	u.PrintSeparator("git")
	if forceInstall || !u.IsExecutableInstalled("git") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "git"); err != nil {
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
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "github-cli"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "gh")
	}

	// Install stow
	u.PrintSeparator("stow")
	if forceInstall || !u.IsExecutableInstalled("stow") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "stow"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "stow")
	}

	// Install ruby
	u.PrintSeparator("ruby")
	if forceInstall || !u.IsExecutableInstalled("ruby") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "ruby"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "ruby")
	}

	// Install zsh-autosuggestions
	u.PrintSeparator("zsh-autosuggestions")
	if forceInstall || !u.IsExecutableInstalled("zsh-autosuggestions") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "zsh-autosuggestions"); err != nil {
			log.Fatalf("error: %v", err)
		}
		installedPackages = append(installedPackages, "zsh-autosuggestions")
	}

	// Install zsh-syntax-highlighting
	u.PrintSeparator("zsh-syntax-highlighting")
	if forceInstall || !u.IsExecutableInstalled("zsh-syntax-highlighting") {
		if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "zsh-syntax-highlighting"); err != nil {
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
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "nodejs", "npm"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("fnm", "install", "22"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install python + dependencies
	u.PrintSeparator("Python and dependencies")
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "python", "python-pip"); err != nil {
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
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "zig"); err != nil {
		log.Fatalf("error: %v", err)
	}
	installedPackages = append(installedPackages, "zig")

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
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "ninja", "gettext", "cmake", "unzip", "curl", "base-devel", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}
  if err := u.RunCmd("npm", "install", "-g", "neovim"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ascii-image-converter
	u.PrintSeparator("ascii-image-converter")
	if err := u.RunCmd("go", "install", "github.com/TheZoraiz/ascii-image-converter@latest"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ripgrep
	u.PrintSeparator("ripgrep")
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "ripgrep"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install fzf
	u.PrintSeparator("fzf")
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "fzf"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install bat
	u.PrintSeparator("bat")
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "bat"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "ln", "-sf", "/usr/bin/bat", "/usr/bin/bat"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install lynx
	u.PrintSeparator("lynx")
	if err := u.RunCmd("sudo", "pacman", "-S", "--noconfirm", "--needed", "lynx"); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("isDesktop: %v", isDesktop)

	// Install desktop environment packages
	if isDesktop {
		// Install catppuccin cursor
		u.PrintSeparator("Catppuccin Cursor")
		if err := u.DownloadFile("https://github.com/catppuccin/cursors/releases/download/v1.0.2/catppuccin-mocha-dark-cursors.zip", "catppuccin-mocha-dark-cursors.zip"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "mkdir", "-p", "/usr/share/icons"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("sudo", "unzip", "catppuccin-mocha-dark-cursors.zip", "-d", "/usr/share/icons"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("gsettings", "set", "org.gnome.desktop.interface", "cursor-theme", "'catppuccin-mocha-dark-cursors'"); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := u.RunCmd("gsettings", "set", "org.gnome.desktop.interface", "cursor-size", "24"); err != nil {
			log.Fatalf("error: %v", err)
		}

	}

	log.Info("Bootstrapping complete.")
	log.Infof("Installed packages: %v", installedPackages)
}
