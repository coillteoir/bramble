/*
Copyright 2024 David Lynch davite3@protonmail.com
*/

package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"
	"os"
	"time"

	"bramble-git-proxy/util"
	"bramble-git-proxy/v1alpha1"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	defaultPort = 9999
	defaultPath = "config/samples/test_repo.yaml"
	headPrefix  = "refs/heads/"
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

		path, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		config, err := util.LoadConfig(path)
		if err != nil {
			return err
		}

		k8sClient, err := util.InitClient()
		if err != nil {
			return err
		}

		http.HandleFunc("/webhook", func(writer http.ResponseWriter, request *http.Request) {
			executionSpec, err := util.ProcessPushEvent(request, config, sugar)
			if err != nil {
				_, err = fmt.Fprintf(writer, "ERROR: %v", err)
				if err != nil {
					sugar.Error("Writer failed")
				}
			}

			util.ExecutePipeline(executionSpec, k8sClient, sugar)
			message := fmt.Sprintf("executing pipeline %v on branch %v", executionSpec.Pipeline, executionSpec.Branch)
			_, err = writer.Write([]byte(message))
			if err != nil {
				sugar.Error("Writer failed")
			}
		})

		http.HandleFunc("/run/{pipeline}/{owner}/{repo}/{branch}", func(writer http.ResponseWriter, request *http.Request) {
			pipeline := request.PathValue("pipeline")
			owner := request.PathValue("owner")
			repo := request.PathValue("repo")
			branch := request.PathValue("branch")
			sugar.Infof("pipeline tested is: %s", pipeline)
			spec := &v1alpha1.ExecutionSpec{
				Pipeline: pipeline,
				Repo:     filepath.Join("https://github.com",owner, repo),
				Branch:   branch,
			}

			err := util.ExecutePipeline(spec, k8sClient, sugar)
			if err != nil {
				sugar.Error(err)
			}
		})

		server := http.Server{
			Addr:              fmt.Sprintf(":%v", port),
			ReadHeaderTimeout: 3 * time.Second,
		}

		sugar.Infof("GIT PROXY RUNNING ON PORT: %v", port)
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
	rootCmd.Flags().BoolP(
		"dry-run", "d", false,
		"Just generate execution resources without running them.",
	)
	rootCmd.Flags().IntP(
		"port", "p", defaultPort,
		"Port for http server to listen to.",
	)
	rootCmd.Flags().StringP(
		"config", "c", defaultPath,
		"Path to the main config file",
	)
}
