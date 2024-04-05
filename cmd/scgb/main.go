package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"scgb/internal"
)

//go:embed generated/source-dir
var appSrcDir string

//go:embed generated/app-name
var appName string

type exitCode int

const (
	exitCodeOk    exitCode = 0
	exitCodeError exitCode = 1
)

func runMain() exitCode {
	fmt.Println("running main...") // TODO: EDIT ME
	return exitCodeOk
}

func getHashFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "."+appName, "hash"), nil
}

func main() {
	// compare with current hash
	needsRecompiling, err := checkNeedsRecompiling(appSrcDir)
	if err != nil {
		handleExitCode(err)
	}

	if needsRecompiling {
		if err := Compile(path.Join(appSrcDir, "cmd", appName)); err != nil {
			handleExitCode(err)
		}

		log.Println("recompiling...")

		// write current hash
		if err := writeCurrentHash(appSrcDir); err != nil {
			handleExitCode(err)
		}

		executable, err := os.Executable()
		if err != nil {
			handleExitCode(err)
		}

		// execute binary in `executable`
		log.Printf("re-running %s\n", executable)
		cmd := exec.Command(executable)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			handleExitCode(err)
		}

		handleExitCode(nil)
	}

	code := runMain()
	os.Exit(int(code))
}

func checkNeedsRecompiling(appSrcPath string) (bool, error) {
	hashFilePath, err := getHashFilePath()
	if err != nil {
		return false, err
	}

	f, err := os.Open(hashFilePath)
	if os.IsNotExist(err) {
		return true, nil
	}
	defer f.Close()

	if err != nil {
		// check if error is file not found
		return false, err
	}

	// read f contents to string
	bs, err := io.ReadAll(f)

	if err != nil {
		return false, err
	}

	existingHash := string(bs)
	currentHash, err := internal.Hash(appSrcPath)
	if err != nil {
		return false, err
	}

	if currentHash != existingHash {
		return true, nil
	}

	return false, nil
}

func writeCurrentHash(appSrcPath string) error {
	hashFilePath, err := getHashFilePath()
	if err != nil {
		return err
	}

	currentHash, err := internal.Hash(appSrcPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(hashFilePath), 0666); err != nil {
		return err
	}

	f, err := os.Create(hashFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(currentHash)
	if err != nil {
		return err
	}

	return nil
}

func Compile(selfPath string) error {
	execCmd := exec.Command("go", "install")
	execCmd.Dir = selfPath
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return err
	}

	return nil
}

func toExitCode(err error) exitCode {
	if err != nil {
		log.Println(err)
		return exitCodeError
	}
	return exitCodeOk
}

func handleExitCode(err error) {
	os.Exit(int(toExitCode(err)))
}
