package quota

import (
	"fmt"

	kubequota "github.com/aauren/kube-quota/pkg"
	"github.com/aauren/kube-quota/pkg/unit"
)

const (
	HeaderEphemeralStorageReq = "Ephemeral Storage Request"
	HeaderEphemeralStorageLim = "Ephemeral Storage Limit"
	HeaderCPUReq              = "CPU Request"
	HeaderMemReq              = "Mem Request"
	HeaderCPULim              = "CPU Limit"
	HeaderMemLim              = "Mem Limit"
)

type NoValueForHeaderError struct {
	Header string
}

func (n *NoValueForHeaderError) Error() string {
	return fmt.Sprintf("No value for header '%s' was found", n.Header)
}

type ComputeQuota struct {
	CPU kubequota.CPUMilicore
	Mem kubequota.MemBytes
}

type StorageQuota struct {
	StorageClasses map[string]*StorageClassQuota
	Ephemeral      *EphemeralQuota
}

func (s *StorageQuota) HasEphemeralQuota() bool {
	return s.Ephemeral != nil
}

func (s *StorageQuota) HasStorageClassQuota() bool {
	return len(s.StorageClasses) > 0
}

func (s *StorageQuota) TableHeader() []string {
	header := make([]string, 0)
	if s.HasEphemeralQuota() {
		header = append(header, HeaderEphemeralStorageReq, HeaderEphemeralStorageLim)
	}
	return header
}

type StorageClassQuota struct {
	Name     string
	Requests kubequota.StorageBytes
	Claims   kubequota.ClaimsNum
}

type EphemeralQuota struct {
	Requests kubequota.StorageBytes
	Limits   kubequota.StorageBytes
}

func (e *EphemeralQuota) ValueForHeader(hdr string) (unit.UnitWriter, error) {
	switch hdr {
	case HeaderEphemeralStorageReq:
		return unit.NewUnitWriter(e.Requests, unit.Bytes)
	case HeaderEphemeralStorageLim:
		return unit.NewUnitWriter(e.Limits, unit.Bytes)
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

func (e *EphemeralQuota) ComparativeUsage(hdr string, totalQuota *EphemeralQuota) (*kubequota.Percentage, error) {
	switch hdr {
	case HeaderEphemeralStorageReq:
		return &kubequota.Percentage{Parts: int64(e.Requests), Whole: int64(totalQuota.Requests)}, nil
	case HeaderEphemeralStorageLim:
		return &kubequota.Percentage{Parts: int64(e.Limits), Whole: int64(totalQuota.Limits)}, nil
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

// WorkloadQuota used to represent the quota reserved by a single workload. Most commonly this represents a container, but it can also be
// used to aggregate workloads when the discrete workloads are no longer needed and only the totals are important
type WorkloadQuota struct {
	Name         string
	Request      *ComputeQuota
	Limit        *ComputeQuota
	StorageQuota *StorageQuota
}

func (w *WorkloadQuota) TableHeader() []string {
	header := make([]string, 0)
	if w.Request != nil {
		header = append(header, HeaderCPUReq, HeaderMemReq)
	}
	if w.Limit != nil {
		header = append(header, HeaderCPULim, HeaderMemLim)
	}
	return header
}

func (w *WorkloadQuota) ValueForHeader(hdr string) (unit.UnitWriter, error) {
	switch hdr {
	case HeaderCPUReq:
		return unit.NewUnitWriter(w.Request.CPU, unit.Cores)
	case HeaderMemReq:
		return unit.NewUnitWriter(w.Request.Mem, unit.Bytes)
	case HeaderCPULim:
		return unit.NewUnitWriter(w.Limit.CPU, unit.Cores)
	case HeaderMemLim:
		return unit.NewUnitWriter(w.Limit.Mem, unit.Bytes)
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

func (w *WorkloadQuota) ComparativeUsage(hdr string, totalQuota *WorkloadQuota) (*kubequota.Percentage, error) {
	switch hdr {
	case HeaderCPUReq:
		return &kubequota.Percentage{Parts: int64(w.Request.CPU), Whole: int64(totalQuota.Request.CPU)}, nil
	case HeaderMemReq:
		return &kubequota.Percentage{Parts: int64(w.Request.Mem), Whole: int64(totalQuota.Request.Mem)}, nil
	case HeaderCPULim:
		return &kubequota.Percentage{Parts: int64(w.Limit.CPU), Whole: int64(totalQuota.Limit.CPU)}, nil
	case HeaderMemLim:
		return &kubequota.Percentage{Parts: int64(w.Limit.Mem), Whole: int64(totalQuota.Limit.Mem)}, nil
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

func (w *WorkloadQuota) ComparativeUsageAsWriter(hdr string, totalQuota *WorkloadQuota) (unit.UnitWriter, error) {
	p, err := w.ComparativeUsage(hdr, totalQuota)
	if err != nil {
		return nil, err
	}

	switch hdr {
	case HeaderCPUReq, HeaderCPULim:
		return unit.NewUnitWriter(p, unit.PercentCores)
	case HeaderMemReq, HeaderMemLim:
		return unit.NewUnitWriter(p, unit.PercentBytes)
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

type PodQuota struct {
	Name           string
	Namespace      string
	WorkloadQuotas []*WorkloadQuota
}

type NamespaceWorkloadQuota struct {
	Name      string
	Namespace string
	PodQuotas []*PodQuota
}

type KubeQuota struct {
	WQ *WorkloadQuota
	SQ *StorageQuota
}

func (k *KubeQuota) ValueForHeader(hdr string) (unit.UnitWriter, error) {
	switch hdr {
	case HeaderCPUReq, HeaderMemReq, HeaderCPULim, HeaderMemLim:
		return k.WQ.ValueForHeader(hdr)
	case HeaderEphemeralStorageReq, HeaderEphemeralStorageLim:
		return k.SQ.Ephemeral.ValueForHeader(hdr)
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}

func (k *KubeQuota) HasEphemeralQuota() bool {
	return k.SQ.HasEphemeralQuota()
}

func (k *KubeQuota) HasStorageClassQuota() bool {
	return k.SQ.HasStorageClassQuota()
}

func (k *KubeQuota) HasWorkloadQuota() bool {
	return k.WQ != nil
}

func (k *KubeQuota) TableHeader() []string {
	header := make([]string, 0)
	if k.HasWorkloadQuota() {
		header = append(header, k.WQ.TableHeader()...)
	}
	if k.HasEphemeralQuota() {
		header = append(header, k.SQ.TableHeader()...)
	}
	return header
}

type QuotaUsage struct {
	KQ  *KubeQuota
	NWQ *NamespaceWorkloadQuota
}

func (qu *QuotaUsage) ValueForHeader(hdr string) (unit.UnitWriter, error) {
	wq := qu.NWQ.Sum()
	switch hdr {
	case HeaderCPUReq, HeaderMemReq, HeaderCPULim, HeaderMemLim:
		return wq.ComparativeUsageAsWriter(hdr, qu.KQ.WQ)
	case HeaderEphemeralStorageReq, HeaderEphemeralStorageLim:
		p, err := wq.StorageQuota.Ephemeral.ComparativeUsage(hdr, qu.KQ.SQ.Ephemeral)
		if err != nil {
			return nil, err
		}
		return unit.NewUnitWriter(p, unit.PercentBytes)
	}

	return nil, &NoValueForHeaderError{Header: hdr}
}
