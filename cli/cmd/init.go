/*
Copyright © 2024 David Lynch davite3@protonmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"bramble/util"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func initRepository(path string) error {
	repoInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if !repoInfo.IsDir() {
		return errors.New("file is not a directory")
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		fmt.Println(remote.Config().URLs[0])
	}
	bramblePath := filepath.Join(path, ".bramble", "pipelines")
	_, err = os.Stat(bramblePath)
	if !os.IsNotExist(err) {
		return errors.New(".bramble directory found, remove if you wish to re-initialze this repository")
	}

	fmt.Printf("Initializing repository '%v'\n", path)
	err = os.MkdirAll(bramblePath, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Created directory: '%v'\n", bramblePath)
	fmt.Printf("\nBramble directory initialized at '%v'.\n", path)
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
		installed, err := util.CheckInstallation()
		if err != nil {
			return err
		}
		if installed {
			return errors.New("system not installed on cluster")
		}
		// Once we create the .bramble/pipelines directory,
		// Right now the only Custom Resource the user needs to care about is pipelines.
		// Every other operation is automated.
		// Maybe creating a demo style "Hello world" pipeline is the move.
		// Or template pipelines depending on a given language
		// Like a C-Build-Release etc.

		return initRepository(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
