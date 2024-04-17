package util

import (
	"fmt"
)

const (
	resources = "https://raw.githubusercontent.com/coillteoir/bramble/master/resources.yaml"
)

func Install(install bool) error {
	installed, err := CheckInstallation()
	if err != nil {
		return err
	}
	action := ""
	if install && installed {
		fmt.Println("Bramble is already installed on this cluster, enjoy!")
		return nil
	}
	if !install && !installed {
		fmt.Println("Bramble is not installed on this cluster")
		return nil
	}
	if install && !installed {
		action = "apply"
	}
	if !install && installed {
		action = "delete"
	}

	err = ExecKubectl([]string{action, "-f", resources})
	if err != nil {
		return err
	}

	if install {
		fmt.Println("\n\nBramble successfully installed to cluster. Deployments may need time to stabilize")
	} else {
		fmt.Println("\n\nBramble successfully uninstalled from cluster.")
	}
	return nil
}
