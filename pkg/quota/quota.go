package quota

import v1 "k8s.io/api/core/v1"

func ForKubeQuota(kubeq *v1.ResourceQuota) *KubeQuota {
	kq := KubeQuota{}
	kq.WorkloadQuota = ConvertK8sHardToWorkload(kubeq.Spec.Hard)
	kq.StorageQuota = ConvertK8sHardToStorage(kubeq.Spec.Hard)
	return &kq
}
