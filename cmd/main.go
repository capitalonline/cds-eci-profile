package main

import (
	"context"
	"github.com/capitalonline/cds-eci-profile/profile"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/home/cck/.kube/config")
	if err != nil {
		panic(err)
	}
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	profileConfig := &profile.Config{
		K8sClient:  k8sClient,
		CACertPath: "/etc/kubernetes/pki/ca.crt",
		CAKeyPath:  "/etc/kubernetes/pki/ca.key",
	}
	manager, err := profile.NewManager(profileConfig)
	if err != nil {
		klog.Fatalf("failed to create eci-profile manager: %q", err)
	}
	klog.Infof("ready to start eci-profile manager service")
	if err := manager.Run(context.TODO()); err != nil {
		klog.Fatalf("run profile service failed: %q", err)
	}
}
