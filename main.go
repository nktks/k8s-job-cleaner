package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "local kubeconfig if you run from outside cluster.")
	masterURL  = flag.String("url", "", "master url if you run from outside cluster.")
	masterName = flag.String("name", "", "master name if you run from outside cluster.")
	namespace  = flag.String("namespace", "", "namespace that you want watch jobs.")
	ttl        = flag.Int("ttl", 10, "ttl of completed or failed job deletion minutes.")
)

func main() {
	flag.Parse()
	var config *rest.Config
	if *kubeconfig != "" {
		if *masterURL == "" && *masterName == "" {
			log.Fatal("need at least one -url or -name")
		}
		var m string
		if *masterName != "" {
			var err error
			kb, err := read(*kubeconfig)
			if err != nil {
				log.Fatal(err)
			}
			m, err = GetMasterURLFromKubeConfig(kb, *masterName)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			m = *masterURL
		}

		c, err := clientcmd.BuildConfigFromFlags(m, *kubeconfig)
		if err != nil {
			log.Fatal(err.Error())
		}
		config = c
	} else {
		c, err := rest.InClusterConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
		config = c
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	jc := NewJobCleaner(clientset.BatchV1().Jobs(*namespace), time.Duration(*ttl)*time.Minute)
	if err := jc.Watch(context.Background()); err != nil {
		log.Fatal(err)
	}
}
func read(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
