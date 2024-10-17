package cli

import (
	"fmt"

	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	headerForeground = color.FgGreen
	headerPadding    = 10
)

var (
	headerFmt = color.New(headerForeground, color.Underline).SprintfFunc()
)

type Byter interface {
	ToBytes() int64
}

func TabularizeTotalQuota(tq *quota.TotalQuota) table.Table {
	tbl := table.New("CPU Request Total", "Mem Request Total", "CPU Limit Total", "Mem Limit Total")
	tbl.WithHeaderFormatter(headerFmt).WithPadding(headerPadding)

	wq := tq.Sum()

	tbl.AddRow(formatMilliCores(wq.Request.CPU), formatBytes(&wq.Request.Mem), formatMilliCores(wq.Limit.CPU), formatBytes(&wq.Limit.Mem))

	return tbl
}

func TabularizeKubeQuota(kq *quota.KubeQuota) table.Table {
	headers := make([]interface{}, 0)
	values := make([]interface{}, 0)
	if kq.HasComputeQuota() {
		headers = append(headers, "CPU Request Total", "Mem Request Total", "CPU Limit Total", "Mem Limit Total")
		values = append(values, formatMilliCores(kq.Request.CPU), formatBytes(&kq.Request.Mem), formatMilliCores(kq.Limit.CPU),
			formatBytes(&kq.Limit.Mem))
	}
	if kq.HasEphemeralQuota() {
		headers = append(headers, "Ephemeral Storage Request", "Ephemeral Storage Limit")
		values = append(values, formatBytes(&kq.Ephemeral.Requests), formatBytes(&kq.Ephemeral.Limits))
	}

	tbl := table.New(headers...)
	tbl.WithHeaderFormatter(headerFmt).WithPadding(headerPadding)

	tbl.AddRow(values...)

	return tbl
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

func formatMilliCores(cores quota.CPUMilicore) string {
	const (
		Core = 1000
	)

	if cores < Core {
		return fmt.Sprintf("%d Millicores", cores)
	} else {
		return fmt.Sprintf("%.1f Cores", float64(cores)/Core)
	}
}
