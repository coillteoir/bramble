package util

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func generateConfig() (*rest.Config, error) {
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
}

func CreateK8sClient() (*kubernetes.Clientset, error) {
	config, err := generateConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func CheckInstallation() (bool, error) {
	clientset, err := CreateK8sClient()
	if err != nil {
		return false, err
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	return (slices.IndexFunc(namespaces.Items, func(ns corev1.Namespace) bool {
		return ns.ObjectMeta.Name == "bramble"
	}) != -1), nil
}

func Install(install bool) error {
	program := "kubectl"
	command := exec.Command("which", program)
	if err := command.Run(); err != nil {
		return fmt.Errorf("program %v not found on system, please install to install Bramble", program)
	}

	installed, err := CheckInstallation()
	if err != nil {
		return err
	}
	action := ""
	if install && !installed {
		action = "apply"
	} else {
		action = "delete"
	}

	command = exec.Command(program, action, "-f", "https://raw.githubusercontent.com/coillteoir/bramble/master/resources.yaml")
	pipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	if err := command.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}
	if err != nil {
		if !errors.Is(io.EOF, err) {
			return err
		}
	}
	return nil
}
