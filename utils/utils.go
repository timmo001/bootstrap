package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

func IsExecutableInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func DeleteDir(dir string) error {
	log.Infof("Deleting directory: %s", dir)

	// Delete the directory
	return os.RemoveAll(dir)
}

func DeleteFile(file string) error {
	log.Infof("Deleting file: %s", file)

	// Delete the file
	return os.Remove(file)
}

func DownloadFile(url, dest string) error {
	log.Infof("Downloading file: %s", url)

	// Download the file
	cmd := exec.Command("curl", "-L", "-o", dest, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExistsDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func RunCmdNoInput(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCmd(name string, arg ...string) error {
	log.Infof("Running command: %s %v", name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCmdInDir(dir, name string, arg ...string) error {
	log.Infof("Running command in directory: %s %s %v", dir, name, arg)

	// Run the command
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func IsLineInFile(file, line string) (bool, error) {
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

func AddIfMissingToFile(file, line string) error {
	if exists, err := IsLineInFile(file, line); err != nil {
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

func UpdateOrCloneRepo(repoURL, destDir string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		// Directory does not exist, clone the repo
		return RunCmd("git", "clone", "--depth", "1", repoURL, destDir)
	} else {
		// Directory exists, pull the latest changes
		return RunCmdInDir(destDir, "git", "pull")
	}
}

func PrintSeparator(msg ...string) {
	log.Print("================================================================================")
	log.Print("")
	log.Print("")
	log.Printf("== %s ==", strings.Join(msg, " "))
	log.Print("")
	log.Print("")
	log.Print("================================================================================")
}
