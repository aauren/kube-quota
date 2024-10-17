package quota

import (
	"context"
	"fmt"

	"github.com/aauren/kube-quota/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FindByNSAndName(ctx context.Context, ns, name string) (*v1.ResourceQuota, error) {
	k8s, err := kubernetes.GetClientSet()
	if err != nil {
		return nil, err
	}

	if name == "" {
		rql, err := k8s.CoreV1().ResourceQuotas(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		if len(rql.Items) > 1 {
			return nil, fmt.Errorf("more than 1 resource quota exists in namespace %s, please add a valid name", ns)
		} else if len(rql.Items) < 1 {
			return nil, fmt.Errorf("no resource quotas existed in namespace %s, please try a different namespace", ns)
		}

		return &rql.Items[0], nil
	}

	return k8s.CoreV1().ResourceQuotas(ns).Get(ctx, name, metav1.GetOptions{})
}
