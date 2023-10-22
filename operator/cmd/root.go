/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"context"
	"fmt"

	"github.com/spf13/cobra"
	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lugh",
	Short: "A brief description of your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := rest.InClusterConfig()

		if err != nil {
			return err
		}

		fmt.Println("Config created")

		clientset, err := kubernetes.NewForConfig(config)

		if err != nil {
			return err
		}

		fmt.Println("Clientset created")

		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			fmt.Println("Pods can't be found")
			return err
		}

		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
