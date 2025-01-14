package quota

import (
	"strings"

	kubequota "github.com/aauren/kube-quota/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	storageClassSuffix = ".storageclass.storage.k8s.io/"
)

func asInt64(quantity resource.Quantity) int64 {
	var by []byte
	_, suffix := quantity.CanonicalizeBytes(by)
	if string(suffix) == "" {
		return quantity.MilliValue()
	}
	// klog.Infof("Suffix was: %s", suffix)
	return quantity.Value()
}

func ConvertK8sResourceList(rl v1.ResourceList) *ComputeQuota {
	cQuota := ComputeQuota{}

	cpu := rl[v1.ResourceCPU]
	cQuota.CPU = kubequota.CPUMilicore(asInt64(cpu))
	mem := rl[v1.ResourceMemory]
	cQuota.Mem = kubequota.MemBytes(mem.Value())

	return &cQuota
}

func ConvertK8sHardToWorkload(rl v1.ResourceList) *WorkloadQuota {
	wq := WorkloadQuota{
		Limit:   &ComputeQuota{},
		Request: &ComputeQuota{},
	}
	for key, val := range rl {
		//nolint:exhaustive // We don't care to be exhaustive here
		switch key {
		case v1.ResourceRequestsCPU, v1.ResourceCPU:
			wq.Request.CPU = kubequota.CPUMilicore(asInt64(val))
		case v1.ResourceRequestsMemory, v1.ResourceMemory:
			wq.Request.Mem = kubequota.MemBytes(val.Value())
		case v1.ResourceLimitsCPU:
			wq.Limit.CPU = kubequota.CPUMilicore(asInt64(val))
		case v1.ResourceLimitsMemory:
			wq.Limit.Mem = kubequota.MemBytes(val.Value())
		}
	}

	return &wq
}

func ConvertK8sHardToStorage(rl v1.ResourceList) *StorageQuota {
	sq := StorageQuota{}
	for key, val := range rl {
		//nolint:exhaustive // We don't care to be exhaustive here
		switch key {
		case v1.ResourceRequestsEphemeralStorage, v1.ResourceEphemeralStorage:
			if sq.Ephemeral == nil {
				sq.Ephemeral = &EphemeralQuota{}
			}
			sq.Ephemeral.Requests = kubequota.StorageBytes(val.Value())
			continue
		case v1.ResourceLimitsEphemeralStorage:
			if sq.Ephemeral == nil {
				sq.Ephemeral = &EphemeralQuota{}
			}
			sq.Ephemeral.Limits = kubequota.StorageBytes(val.Value())
			continue
		}

		//nolint:staticcheck // It's ok for this to be empty for now
		if strings.Contains(string(key), storageClassSuffix) {
			// To be implemented at some point in the future
			// sq.AddNewClass(key, val)
		}
	}

	return &sq
}
