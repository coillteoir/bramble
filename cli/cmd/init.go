/*
Copyright © 2024 David Lynch davite3@protonmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func checkInstallation() error {
	return nil
}

func initRepository(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	fmt.Printf("Initializing repository '%v'\n", path)

	if !fileInfo.IsDir() {
		return errors.New("file is not a directory")
	}

	if _, err = os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		return err
	}

	fmt.Printf("Git repo found at '%v'\n", path)

	err = os.MkdirAll(filepath.Join(path, ".bramble", "pipelines"), os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Printf("Created directory: '%v'\n", filepath.Join(path, ".bramble", "pipelines"))

	fmt.Printf("\nBramble directory initialized at '%v'\n.", path)
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init [REPOSITORY]",
	Short: "Initializes a git repo to be used with bramble",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("can only initialize one directory")
		}
		if len(args) == 0 {
			return errors.New("no path specified")
		}

		// Once we create the .bramble/pipelines directory,
		// Right now the only Custom Resource the user needs to care about is pipelines.
		// Every other operation is automated.
		// Maybe creating a demo style "Hello world" pipeline is the move.
		// Or template pipelines depending on a given language
		// Like a C-Build-Release etc.

		if args[0] == "." {
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			initRepository(path)
		}

		return initRepository(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
