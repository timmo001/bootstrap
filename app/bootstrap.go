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

	// Copy .editorconfig
	u.PrintSeparator("Copying .editorconfig")
	if err := u.RunCmd("cp", ".editorconfig", home); err != nil {
		log.Fatalf("error: %v", err)
	}

	// // Install oh-my-zsh
	// u.PrintSeparator("oh-my-zsh")
	// exists, err := u.ExistsDir(home + "/.oh-my-zsh")
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// if forceInstall || !exists {
	// 	if err := u.DeleteDir(home + "/.oh-my-zsh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	if err := u.DownloadFile("https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh", "omz-install.sh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	if err := u.RunCmdNoInput("sh", "omz-install.sh"); err != nil {
	// 		log.Errorf("error: %v", err)
	// 		if err := u.DeleteFile("omz-install.sh"); err != nil {
	// 			log.Fatalf("error: %v", err)
	// 		}
	// 		log.Fatal("error installing oh-my-zsh")
	// 	}
	// 	if err := u.DeleteFile("omz-install.sh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	installedPackages = append(installedPackages, "oh-my-zsh")
	// }

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

	// // Install starship
	// u.PrintSeparator("starship")
	// if forceInstall || !u.IsExecutableInstalled("starship") {
	// 	if err := u.RunCmd("curl", "-fsSL", "https://starship.rs/install.sh", "-o", "starship-install.sh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	if err := u.RunCmd("chmod", "+x", "starship-install.sh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	if err := u.RunCmd("./starship-install.sh", "--yes"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	if err := u.DeleteFile("starship-install.sh"); err != nil {
	// 		log.Fatalf("error: %v", err)
	// 	}
	// 	installedPackages = append(installedPackages, "starship")
	// }

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

	// Install neovim (update build-essential replacement)
	u.PrintSeparator("Neovim")
	if err := u.RunCmd("rpm-ostree", "install", "ninja-build", "gettext", "-y"); err != nil {
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

	// Install ascii-image-converter
	u.PrintSeparator("ascii-image-converter")
	if err := u.RunCmd("go", "install", "github.com/TheZoraiz/ascii-image-converter@latest"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install ripgrep
	u.PrintSeparator("ripgrep")
	if err := u.RunCmd("rpm-ostree", "install", "ripgrep", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Install bat
	u.PrintSeparator("bat")
	if err := u.RunCmd("rpm-ostree", "install", "bat", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("sudo", "ln", "-s", "/usr/bin/batcat", "/usr/bin/bat"); err != nil {
		log.Errorf("error: %v", err)
	}

	// Install lynx
	u.PrintSeparator("lynx")
	if err := u.RunCmd("rpm-ostree", "install", "lynx", "-y"); err != nil {
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
	if err := u.RunCmd("rpm-ostree", "install", "fira-code-fonts", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmd("rpm-ostree", "install", "hack-fonts", "-y"); err != nil {
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

	log.Info("Bootstrapping complete.")
	log.Infof("Installed packages: %v", installedPackages)
}
