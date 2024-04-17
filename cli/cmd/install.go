/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"bramble/util"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the most recent version of bramble to your kubernetes cluster.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Intalling bramble...")
		err := util.Install(true)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
