package kubernetes

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gojektech/proctor-engine/config"
	"github.com/gojektech/proctor-engine/logger"
)

func KubeConfig() string {
	kubeConfig := new(string)
	if config.Environment() == "development" {
		home := os.Getenv("HOME")

		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		logger.Info("kubeconfig is provided", kubeConfig)
	} else {
		logger.Info("using in cluster kubeconfig")
	}
	return *kubeConfig
}
