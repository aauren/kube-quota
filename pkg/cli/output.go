package cli

import (
	"fmt"

	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func TabularizeTotalQuota(tq *quota.TotalQuota) table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("CPU Request Total", "Mem Request Total", "CPU Limit Total", "Mem Limit Total")
	tbl.WithHeaderFormatter(headerFmt).WithPadding(10)

	wq := tq.Sum()

	tbl.AddRow(formatMilliCores(wq.Request.CPU), formatBytes(wq.Request.Mem), formatMilliCores(wq.Limit.CPU), formatBytes(wq.Limit.Mem))

	return tbl
}

func formatBytes(bytes quota.MemBytes) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

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
