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

type NoValueForHeader struct {
	Header string
}

func (n *NoValueForHeader) Error() string {
	return fmt.Sprintf("No value for header '%s' was found", n.Header)
}

type ComputeQuota struct {
	CPU kubequota.CPUMilicore
	Mem kubequota.MemBytes
}

type StorageQuota struct {
	StorageClasses []*StorageClassQuota
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

	return nil, &NoValueForHeader{Header: hdr}
}

type WorkloadQuota struct {
	Name    string
	Request *ComputeQuota
	Limit   *ComputeQuota
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

	return nil, &NoValueForHeader{Header: hdr}
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

	return nil, &NoValueForHeader{Header: hdr}
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
