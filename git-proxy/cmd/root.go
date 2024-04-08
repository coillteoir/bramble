/*
Copyright 2024 David Lynch davite3@protonmail.com
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	//	"bramble-git-proxy/v1alpha1"

	"github.com/google/go-github/v59/github"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	// apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	//"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	// apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultPort = 9999
	defaultPath = "config/samples/test_repo.yaml"
	headPrefix  = "refs/heads/"
)

type proxyConfig struct {
	Provider string
	Owner    string
	Repo     string
	Pairings map[string]string
}

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

		config, err := loadConfig(path)
		if err != nil {
			return err
		}

		k8sClient, err := initClient()
		if err != nil {
			return err
		}

		fmt.Println(k8sClient)

		http.HandleFunc("/webhook", func(writer http.ResponseWriter, request *http.Request) {
			response, err := processPushEvent(request, config, sugar)
			if err != nil {
				_, err = fmt.Fprintf(writer, "ERROR: %v", err)
				if err != nil {
					sugar.Error("Writer failed")
				}
			}
			_, err = writer.Write([]byte(response))
			if err != nil {
				sugar.Error("Writer failed")
			}
		})

		http.HandleFunc("/test/{pipeline}", func(writer http.ResponseWriter, request *http.Request) {
			pipeline := request.PathValue("pipeline")
			sugar.Infof("pipeline tested is: %s", pipeline)
			/*
				execution := &v1alpha1.Execution{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: pipeline,
						Namespace:    "bramble",
					},
					Spec: v1alpha1.ExecutionSpec{
						Repo:     "https://github.com/coillteoir/bramble",
						Branch:   "git-proxy",
						Pipeline: pipeline,
					},
				}
			*/
			execution := &unstructured.Unstructured{}

			k8sClient.Resource(schema.GroupVersionResource{
				Group:    "pipelines.bramble.dev",
				Version:  "v1alpha1",
				Resource: "executions",
			}).Namespace("bramble").Create(context.TODO(), execution, metav1.CreateOptions{})
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

func initClient() (*dynamic.DynamicClient, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		client, err := dynamic.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	kubeconfig := path.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func loadConfig(configPath string) ([]proxyConfig, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := []proxyConfig{}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func processPushEvent(
	request *http.Request,
	config []proxyConfig,
	sugar *zap.SugaredLogger,
) (string, error) {
	payload, err := github.ValidatePayload(request, nil)
	if err != nil {
		sugar.Errorf("Invalid payload: %v", request)
		return "", err
	}

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		sugar.Errorf("Cannot parse webhook: %v", request)
		return "", err
	}

	switch event := event.(type) {
	case *github.PushEvent:
		sugar.Infof("Pushed!! %v", *event.Ref)
		branchName, found := strings.CutPrefix(*event.Ref, headPrefix)
		if !found {
			return "", errors.New("could not parse branch name")
		}

		for _, repo := range config {
			nameMatch := (repo.Owner == *event.Repo.Owner.Name && repo.Repo == *event.Repo.Name)
			if repo.Provider != "github" || !nameMatch {
				continue
			}
			sugar.Infof("Push to branch: %v Pipeline: %v", branchName, repo.Pairings[branchName])
			if _, exists := repo.Pairings[branchName]; !exists {
				return fmt.Sprintf("No pipelines are configured for branch %v", branchName), nil
			}
			return fmt.Sprintf(
				"Executing pipeline %v on branch %v",
				repo.Pairings[branchName],
				branchName,
			), nil
		}
	default:
		sugar.Infof("DifferentEvent")
		return "Hello from webhook!\n", nil
	}
	return "", nil
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
