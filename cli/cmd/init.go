/*
Copyright © 2024 David Lynch davite3@protonmail.com
*/
package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func checkInstallation() error {
	config, err := (func() (*rest.Config, error) {
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			config, err := rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
			return config, nil
		}

		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil
	}())
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ns := range namespaces.Items {
		if ns.ObjectMeta.Name == "bramble" {
			return nil
		}
	}
	return errors.New("could not find bramble namespace")
}

func initRepository(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	fmt.Printf("Initializing repository '%v'\n", path)
	if !fileInfo.IsDir() {
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

	fmt.Printf("Git repo found at '%v'\n", path)
	err = os.MkdirAll(filepath.Join(path, ".bramble", "pipelines"), os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Created directory: '%v'\n", filepath.Join(path, ".bramble", "pipelines"))
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

		if err := checkInstallation(); err != nil {
			return err
		}

		// Once we create the .bramble/pipelines directory,
		// Right now the only Custom Resource the user needs to care about is pipelines.
		// Every other operation is automated.
		// Maybe creating a demo style "Hello world" pipeline is the move.
		// Or template pipelines depending on a given language
		// Like a C-Build-Release etc.

		if args[0] == "." {
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			return initRepository(path)
		}

		return initRepository(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
