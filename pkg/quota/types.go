package quota

type CPUMilicore int64
type MemBytes int64

type ComputeQuota struct {
	CPU CPUMilicore
	Mem MemBytes
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
