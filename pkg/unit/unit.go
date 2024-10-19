package unit

import (
	"fmt"

	kubequota "github.com/aauren/kube-quota/pkg"
)

type FormatUnit int

const (
	Bytes = iota
	Cores
)

var (
	AllFormatters = []FormatUnit{Bytes, Cores}
)

type Byter interface {
	ToBytes() int64
}

type Unit struct {
	bytes Byter
	unit  FormatUnit
	cores kubequota.CPUMilicore
}

type UnitWriter interface {
	fmt.Stringer
}

func (u *Unit) String() string {
	switch u.unit {
	case Bytes:
		return formatBytes(u.bytes)
	case Cores:
		return formatMilliCores(u.cores)
	}

	return ""
}

func NewUnitWriter(value interface{}, unit FormatUnit) (UnitWriter, error) {
	switch unit {
	case Bytes:
		switch b := value.(type) {
		case kubequota.MemBytes:
			return &Unit{bytes: &b, unit: unit}, nil
		case kubequota.StorageBytes:
			return &Unit{bytes: &b, unit: unit}, nil
		default:
			return nil, fmt.Errorf("unable to cast %v to a valid bytes type, cannot continue", value)
		}
	case Cores:
		switch c := value.(type) {
		case kubequota.CPUMilicore:
			return &Unit{cores: c, unit: unit}, nil
		default:
			return nil, fmt.Errorf("unable to cast %v to a valid core type, cannot continue", value)
		}
	}

	return nil, fmt.Errorf("exhausted all Unit unit cases, cannot create NewUnitWriter")
}

func formatBytes(byter Byter) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	bytes := byter.ToBytes()

	switch {
	case bytes < KB:
		return fmt.Sprintf("%d B", bytes)
	case bytes < MB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	case bytes < GB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes < TB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	default:
		return fmt.Sprintf("%.1f TB", float64(bytes)/TB)
	}
}

func formatMilliCores(cores kubequota.CPUMilicore) string {
	const (
		Core = 1000
	)

	if cores < Core {
		return fmt.Sprintf("%d Millicores", cores)
	} else {
		return fmt.Sprintf("%.1f Cores", float64(cores)/Core)
	}
}
