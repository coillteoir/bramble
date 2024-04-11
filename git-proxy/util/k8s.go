package util

import (
	"context"
	"fmt"
	"os"
	"path"

	"bramble-git-proxy/v1alpha1"

	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func generateUnstructured(execution *v1alpha1.Execution) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "pipelines.bramble.dev/v1alpha1",
			"kind":       "Execution",
			"metadata": map[string]any{
				"generateName": execution.ObjectMeta.GenerateName,
				"namespace":    execution.ObjectMeta.Namespace,
			},
			"spec": map[string]any{
				"repo":     execution.Spec.Repo,
				"branch":   execution.Spec.Branch,
				"pipeline": execution.Spec.Pipeline,
			},
		},
	}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "pipelines.bramble.dev",
		Version: "v1alpha1",
		Kind:    "Execution",
	})
	return obj
}

func ExecutePipeline(spec *v1alpha1.ExecutionSpec, k8sClient *dynamic.DynamicClient, sugar *zap.SugaredLogger) error {
	execution := &v1alpha1.Execution{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%v-execution", spec.Pipeline),
			Namespace:    "default",
		},
		Spec: *spec,
	}

	unstructuredExecution := generateUnstructured(execution)
	_, err := k8sClient.Resource(schema.GroupVersionResource{
		Group:    "pipelines.bramble.dev",
		Version:  "v1alpha1",
		Resource: "executions",
	}).Namespace(execution.ObjectMeta.Namespace).Create(context.TODO(), unstructuredExecution, metav1.CreateOptions{})

	sugar.Infof("created execution %v", execution.ObjectMeta.Name)
	return err
}

func InitClient() (*dynamic.DynamicClient, error) {
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

	kubeconfig := path.Join(os.Getenv("HOME"), ".kube", "config")

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
