/*
Copyright © 2024 David Lynch davite3@protonmail.com
*/
package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

const (
	helloPipeline = `
apiVersion: pipelines.bramble.dev/v1alpha1
kind: Pipeline
metadata:
  name: hello
  namespace: bramble
spec:
  tasks:
  - name: hello
    image: alpine
    command: ["echo", "hello", "world"]
`
	pipelinesKustomization = `
resources:
- hello.yaml`

	baseKustomization = `
resources:
- pipelines`
)

func InitRepository(path string) error {
	repoInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if !repoInfo.IsDir() {
		return errors.New("file is not a directory")
	}

	_, err = git.PlainOpen(path)
	if err != nil {
		return err
	}

	bramblePath := filepath.Join(path, ".bramble")
	_, err = os.Stat(bramblePath)
	if !os.IsNotExist(err) {
		return errors.New(".bramble directory found, remove if you wish to re-initialze this repository")
	}

	fmt.Printf("Initializing repository '%v'\n", path)

	pipelinesPath := filepath.Join(bramblePath, "pipelines")

	err = os.MkdirAll(pipelinesPath, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Created directory: '%v'\n", bramblePath)

	err = os.WriteFile(filepath.Join(bramblePath, "kustomization.yaml"), []byte(baseKustomization), 0o644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(pipelinesPath, "kustomization.yaml"), []byte(pipelinesKustomization), 0o644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(pipelinesPath, "hello.yaml"), []byte(helloPipeline), 0o644)
	if err != nil {
		return err
	}

	fmt.Printf("\nBramble directory initialized at '%v'.\n", path)
	return nil
}
