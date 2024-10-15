package quota

import (
	v1 "k8s.io/api/core/v1"
)

func ConvertK8sResourceList(rl v1.ResourceList) *ComputeQuota {
	cQuota := ComputeQuota{}

	cpu := rl[v1.ResourceCPU]
	cQuota.CPU = CPUMilicore(cpu.Value())
	mem := rl[v1.ResourceMemory]
	cQuota.Mem = MemBytes(mem.Value())

	return &cQuota
}
