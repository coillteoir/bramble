/*
Copyright 2024 David Lynch davite3@protonmail.com
*/

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v59/github"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	defaultPort = 9999
)

var rootCmd = &cobra.Command{
	Use:   "bramble-git-proxy",
	Short: "A proxy between git providers and the Bramble CI/CD framework",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return err
		}

		sugar := logger.Sugar()

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}

		http.HandleFunc("/webhook", func(writer http.ResponseWriter, request *http.Request) {
			payload, err := github.ValidatePayload(request, nil)
			if err != nil {
				sugar.Errorf("Invalid payload: %v", request)
				return
			}

			event, err := github.ParseWebHook(github.WebHookType(request), payload)
			if err != nil {
				sugar.Errorf("Cannot parse webhook: %v", request)
				return
			}

			switch event := event.(type) {
			case *github.PushEvent:
				sugar.Infof("Pushed!! %v", *event.Ref)
				_, err = writer.Write([]byte("Thanks for the push!\n"))
				if err != nil {
					sugar.Errorf("Writer failed: %v", request)
				}
				return

			default:
				sugar.Infof("DifferentEvent")
			}

			_, err = writer.Write([]byte("Hello from webhook!\n"))
			if err != nil {
				sugar.Errorf("Writer failed: %v", request)
				return
			}
		})

		sugar.Infof("GIT PROXY RUNNING ON PORT: %v", port)

		server := http.Server{
			Addr:              fmt.Sprintf(":%v", port),
			ReadHeaderTimeout: 3 * time.Second,
		}

		err = server.ListenAndServe()
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
	rootCmd.Flags().IntP("port", "p", defaultPort, "Port for http server to listen to.")
}
