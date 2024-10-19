package quota

import (
	v1 "k8s.io/api/core/v1"
)

func QuotaForPod(pod *v1.Pod) *PodQuota {
	podQuota := PodQuota{
		Name:           pod.Name,
		Namespace:      pod.Namespace,
		WorkloadQuotas: make([]*WorkloadQuota, 0),
	}

	for _, cnt := range pod.Spec.Containers {
		tmpQuota := WorkloadQuota{
			StorageQuota: &StorageQuota{
				Ephemeral:      &EphemeralQuota{},
				StorageClasses: make(map[string]*StorageClassQuota),
			},
		}
		tmpQuota.Name = cnt.Name
		tmpQuota.Limit = ConvertK8sResourceList(cnt.Resources.Limits)
		tmpQuota.Request = ConvertK8sResourceList(cnt.Resources.Requests)
		podQuota.WorkloadQuotas = append(podQuota.WorkloadQuotas, &tmpQuota)
	}

	return &podQuota
}

func QuotaForPodList(pl *v1.PodList) *NamespaceWorkloadQuota {
	tq := NamespaceWorkloadQuota{
		PodQuotas: make([]*PodQuota, 0),
	}
	for idx := range pl.Items {
		tq.PodQuotas = append(tq.PodQuotas, QuotaForPod(&pl.Items[idx]))
	}

	return &tq
}
