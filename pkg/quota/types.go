package quota

type CPUMilicore int64
type MemBytes int64
type StorageBytes int64
type ClaimsNum int16

func (m *MemBytes) ToBytes() int64 {
	return int64(*m)
}

func (s *StorageBytes) ToBytes() int64 {
	return int64(*s)
}

type ComputeQuota struct {
	CPU CPUMilicore
	Mem MemBytes
}

type StorageQuota struct {
	StorageClasses []*StorageClassQuota
	Ephemeral      *EphemeralQuota
}

type StorageClassQuota struct {
	Name     string
	Requests StorageBytes
	Claims   ClaimsNum
}

type EphemeralQuota struct {
	Requests StorageBytes
	Limits   StorageBytes
}

type WorkloadQuota struct {
	Name    string
	Request *ComputeQuota
	Limit   *ComputeQuota
}

type AggregateQuota struct {
	Name           string
	Namespace      string
	WorkloadQuotas []*WorkloadQuota
}

type TotalQuota struct {
	AggregateQuotas []*AggregateQuota
}

type KubeQuota struct {
	*WorkloadQuota
	*StorageQuota
}

func (k *KubeQuota) HasEphemeralQuota() bool {
	return k.Ephemeral != nil
}

func (k *KubeQuota) HasStorageClassQuota() bool {
	return len(k.StorageClasses) > 0
}

func (k *KubeQuota) HasComputeQuota() bool {
	return k.WorkloadQuota != nil
}
