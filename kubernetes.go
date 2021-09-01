package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/labels"
)

type K8s struct {
	Client  *kubernetes.Clientset
	Context context.Context
}

type KubernetesClient interface {
	GetNode(string) (*v1.Node, error)
	GetNodeList(string) (*v1.NodeList, error)
	GetRegions(string) ([]int, error)
}

// NewKubernetesClient return a Kubernetes client
func NewKubernetesClient(host string, port string, namespace string, kubeConfigPath string) (k8s KubernetesClient, err error) {
	var client *kubernetes.Clientset

	if len(host) > 0 && len(port) > 0 {
		log.Info().Msg("in cluster client created")
		client, err = inClusterConfig()

		if err != nil {
			err = fmt.Errorf("Error loading incluster client:\n%v", err)
			return
		}
	} else {
		log.Info().Msg("creating out of cluster client")
		client, err = loadOutOfClusterK8sClient(kubeConfigPath)

		if err != nil {
			err = fmt.Errorf("Error loading client using kubeconfig:\n%v", err)
			return
		}
	}

	k8s = &K8s{
		Client:  client,
		Context: context.Background(),
	}

	return
}

// GetNode return the node object from given name
func (k *K8s) GetNode(name string) (node *v1.Node, err error) {
	node, err = k.Client.CoreV1().Nodes().Get(k.Context, name, metav1.GetOptions{})
	return
}

// GetNodeList return a list of nodes from a given node pool name, if name is empty all nodes are returned
func (k *K8s) GetNodeList(name string) (nodes *v1.NodeList, err error) {
	opts := metav1.ListOptions{}

	if name != "" {
		selector := map[string]string{
				"cloud.google.com/gke-nodepool": name,
			}
		ls := labels.SelectorFromSet(selector)
		opts.LabelSelector = ls.String()
	}

	nodes, err = k.Client.CoreV1().Nodes().List(k.Context, opts)
	return
}


// GetRegions returns the list of region names, useful determine the amount of regions
func (k *K8s) GetRegions(name string) (regions []int, err error) {
	regions = []int{}
	opts := metav1.ListOptions{}
	zones := []string{"europe-west1-b", "europe-west1-c", "europe-west1-d"}
	var nodes *v1.NodeList

	for _, zone := range zones {
		selector := map[string]string{
			"cloud.google.com/gke-nodepool": name,
			"failure-domain.beta.kubernetes.io/zone": zone,
		}
		ls := labels.SelectorFromSet(selector)
		opts.LabelSelector = ls.String()
		nodes, err = k.Client.CoreV1().Nodes().List(k.Context, opts)
		if err != nil {
			return
		}
		regions = append(regions, len(nodes.Items))

	}

	return
}

// inClusterConfig returns a kubernetes client for authenticating inside the cluster
func inClusterConfig() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

// loadOutOfClusterK8sClient parses a kubeconfig from a file and returns a Kubernetes
// client. It does not support extensions or client auth providers.
func loadOutOfClusterK8sClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// fmt.Printf("%#v", config)
	return clientset, nil
}
