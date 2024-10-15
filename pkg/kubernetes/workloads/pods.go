package workloads

import (
	"context"

	"github.com/aauren/kube-quota/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodsByNamespace(ctx context.Context, ns string) (*v1.PodList, error) {
	k8s, err := kubernetes.GetClientSet()
	if err != nil {
		return nil, err
	}

	return k8s.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
}
