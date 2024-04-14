/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"bramble/util"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall bramble from your kubernetes cluster.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := util.Install(false)
		if err != nil {
			return err
		}
		fmt.Println("\n\nBramble successfully uninstalled from cluster.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
