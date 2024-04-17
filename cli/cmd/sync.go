/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"bramble/util"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync <REPOSITORY>",
	Short: "Synchronize pipelines described in .bramble directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Synchronizing pipleine resources with the cluster.")
		err := util.ExecKubectl([]string{"apply", "-k", filepath.Join(args[0], ".bramble")})
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
