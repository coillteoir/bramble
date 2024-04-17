/*
Copyright © 2024 David Lynch davite3@protonmail.com
*/
package cmd

import (
	"errors"
	"fmt"

	"bramble/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <REPOSITORY>",
	Short: "Initializes a git repo to be used with bramble",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("can only initialize one directory")
		}
		if len(args) == 0 {
			return errors.New("no path specified")
		}
		fmt.Println("Checking if system is installed...")
		installed, err := util.CheckInstallation()
		if err != nil {
			return err
		}
		if !installed {
			return errors.New("system not installed on cluster")
		}
		fmt.Println("Bramble installed on cluster.")
		// Once we create the .bramble/pipelines directory,
		// Right now the only Custom Resource the user needs to care about is pipelines.
		// Every other operation is automated.
		// Maybe creating a demo style "Hello world" pipeline is the move.
		// Or template pipelines depending on a given language
		// Like a C-Build-Release etc.
		return util.InitRepository(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
