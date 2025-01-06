package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/exp/slog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout))
	slog.SetDefault(l)

	kubeconfig, err := getKubeConfig()
	if err != nil {
		slog.Error("", err)
		os.Exit(1)
	}
	slog.Info("kubeconfig loaded, using current context")

	client := kubernetes.NewForConfigOrDie(kubeconfig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.CoreV1().Pods("argocd").List(ctx, v1.ListOptions{
		LabelSelector: "app.kubernetes.io/name", // lists pods that has this key and any value
		// LabelSelector: "app.kubernetes.io/name=server", // lists pod matching key and value
	})
	if err != nil {
		slog.Error("", err)
		os.Exit(1)
	}
	slog.Info("pods fetched", "pods", len(pods.Items))
	for _, p := range pods.Items {
		slog.Info("argocd server pod", "pod", p.Name, "namespace", p.Namespace)
	}
}

func getKubeConfig() (*rest.Config, error) {
	kubeconfigFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		return nil, err
	}
	kubeconfig.Burst = 1000 // make it 100 for prod uses
	kubeconfig.QPS = 500    // make it 50 for prod uses

	return kubeconfig, err
}
