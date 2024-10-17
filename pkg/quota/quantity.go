package quota

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

const (
	storageClassSuffix = ".storageclass.storage.k8s.io/"
)

func ConvertK8sResourceList(rl v1.ResourceList) *ComputeQuota {
	cQuota := ComputeQuota{}

	cpu := rl[v1.ResourceCPU]
	cQuota.CPU = CPUMilicore(cpu.Value())
	mem := rl[v1.ResourceMemory]
	cQuota.Mem = MemBytes(mem.Value())

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
			wq.Request.CPU = CPUMilicore(val.Value())
		case v1.ResourceRequestsMemory, v1.ResourceMemory:
			wq.Request.Mem = MemBytes(val.Value())
		case v1.ResourceLimitsCPU:
			wq.Limit.CPU = CPUMilicore(val.Value())
		case v1.ResourceLimitsMemory:
			wq.Limit.Mem = MemBytes(val.Value())
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
			sq.Ephemeral.Requests = StorageBytes(val.Value())
			continue
		case v1.ResourceLimitsEphemeralStorage:
			if sq.Ephemeral == nil {
				sq.Ephemeral = &EphemeralQuota{}
			}
			sq.Ephemeral.Limits = StorageBytes(val.Value())
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
