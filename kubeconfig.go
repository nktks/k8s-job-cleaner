package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type kubeConfig struct {
	Clusters []struct {
		Name    string
		Cluster struct {
			Server string
		}
	}
}

func GetMasterURLFromKubeConfig(kubeconfig []byte, name string) (string, error) {
	v := kubeConfig{}
	if err := yaml.Unmarshal(kubeconfig, &v); err != nil {
		return "", err
	}

	for _, c := range v.Clusters {
		if c.Name == name {
			return c.Cluster.Server, nil
		}
	}

	return "", fmt.Errorf("could not found %s name from kubeconfig", name)
}
