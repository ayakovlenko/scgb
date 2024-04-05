package internal

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Hash(appSrcPath string) (string, error) {
	files, err := listGoFiles(appSrcPath)
	if err != nil {
		return "", err
	}
	return hashFiles(files)
}

func listGoFiles(appSrcPath string) ([]string, error) {
	files := []string{}
	if err := filepath.Walk(
		appSrcPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return files, nil
}

func hashFiles(files []string) (string, error) {
	hasher := sha256.New()
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := io.Copy(hasher, f); err != nil {
			log.Fatal(err)
		}
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
