package util

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"bramble-git-proxy/v1alpha1"
	"github.com/google/go-github/v59/github"
	"go.uber.org/zap"
)

const (
	headPrefix = "refs/heads/"
)

type ProxyConfig struct {
	Provider string
	Owner    string
	Repo     string
	Pairings map[string]string
}

func LoadConfig(configPath string) ([]ProxyConfig, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := []ProxyConfig{}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ProcessPushEvent(
	request *http.Request,
	config []ProxyConfig,
	sugar *zap.SugaredLogger,
) (*v1alpha1.ExecutionSpec, string, error) {
	payload, err := github.ValidatePayload(request, nil)
	if err != nil {
		sugar.Errorf("Invalid payload: %v", request)
		return nil, "", err
	}

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		sugar.Errorf("Cannot parse webhook: %v", request)
		return nil, "", err
	}

	switch event := event.(type) {
	case *github.PushEvent:
		sugar.Infof("Pushed!! %v", *event.Ref)
		branchName, found := strings.CutPrefix(*event.Ref, headPrefix)
		if !found {
			return nil, "", errors.New("could not parse branch name")
		}

		for _, repo := range config {
			nameMatch := (repo.Owner == *event.Repo.Owner.Name && repo.Repo == *event.Repo.Name)
			if repo.Provider != "github" || !nameMatch {
				continue
			}
			sugar.Infof("Push to branch: %v Pipeline: %v", branchName, repo.Pairings[branchName])
			if _, exists := repo.Pairings[branchName]; !exists {
				return nil, fmt.Sprintf("No pipelines are configured for branch %v", branchName), nil
			}

			spec := &v1alpha1.ExecutionSpec{
				Repo:     repo.Repo,
				Branch:   branchName,
				Pipeline: repo.Pairings[branchName],
			}

			return spec, fmt.Sprintf("executing pipeline %v on branch %v", repo.Pairings[branchName], branchName), nil
		}
	default:
		sugar.Infof("DifferentEvent")
		return nil, "Hello from webhook!\n", nil
	}
	return nil, "", nil
}
