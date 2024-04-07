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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a git repo to be used with bramble",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("please select a directory")
		}

		fileInfo, err := os.Stat(args[0])
		if err != nil {
			if os.IsNotExist(err) {
				return err
			}
		}

		if !fileInfo.IsDir() {
			return errors.New("file is not a directory")
		}

		if _, err = os.Stat(filepath.Join(args[0], ".git")); os.IsNotExist(err) {
			return err
		}

		err = os.Mkdir(filepath.Join(args[0], ".bramble"), os.ModePerm)
		if err != nil {
			return err
		}

		fmt.Println("Bramble directory initialized")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
