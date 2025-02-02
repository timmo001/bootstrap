package main

import (
	"github.com/charmbracelet/log"

	u "github.com/timmo001/bootstrap/utils"
)

func main() {
	u.PrintSeparator("Install or update ghostty")

	if err := u.RunCmd("sudo", "apt", "install", "libgtk-4-dev", "libadwaita-1-dev", "-y"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.UpdateOrCloneRepo("https://github.com/ghostty-org/ghostty", "ghostty"); err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := u.RunCmdInDir("ghostty", "sudo", "zig", "build", "-p", "/usr", "-Doptimize=ReleaseFast"); err != nil {
		log.Fatalf("error: %v", err)
	}

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
}
