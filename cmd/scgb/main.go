package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"scgb/internal"
)

//go:embed generated/source-dir
var appSrcDir string

//go:embed generated/app-name
var appName string

//go:embed generated/app-hash
var existingHash string

type exitCode int

const (
	exitCodeOk    exitCode = 0
	exitCodeError exitCode = 1
)

func runMain() exitCode {
	fmt.Println("running main...") // TODO: EDIT ME
	fmt.Println("hash:", existingHash)
	return exitCodeOk
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
	currentHash, err := internal.Hash(appSrcPath)
	if err != nil {
		return err
	}

	// open file in write mode
	f, err := os.OpenFile(path.Join(appSrcPath, "cmd", appName, "generated", "app-hash"), os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
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
