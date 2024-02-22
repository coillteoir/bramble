/*
Copyright 2024 David Lynch davite3@protonmail.com
*/

package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func hello(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, fmt.Sprintf("%v", req))
		}


var rootCmd = &cobra.Command{
	Use:   "bramble-git-proxy",
	Short: "A proxy between git providers and the Bramble CI/CD system",
	RunE: func(cmd *cobra.Command, args []string) error {

        port, err := cmd.Flags().GetInt("port")
        if err != nil {
            return nil
        }


		http.HandleFunc("/hello", hello)

        fmt.Printf("GIT PROXY RUNNING ON PORT: %v", port)
		err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle.")
    rootCmd.Flags().BoolP("dry-run", "d", false, "Just generate execution resources without running them.")
    rootCmd.Flags().IntP("port", "p", 9999, "Port to listen to");
}
