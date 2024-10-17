package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Byter interface {
	ToBytes() int64
}

func CreateTableWriter() table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.SetStyle(table.StyleColoredBright)
	return tbl
}

func AddNewSection(tbl table.Writer, sectionName string) {
	tbl.AppendSeparator()
	tbl.AppendRow(table.Row{"", strings.ToUpper(sectionName), "", ""})
	tbl.AppendSeparator()
	tbl.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, AutoMerge: true},
		{Number: 2, Align: text.AlignCenter, AutoMerge: true},
		{Number: 3, Align: text.AlignCenter, AutoMerge: true},
		{Number: 4, Align: text.AlignCenter, AutoMerge: true},
	})
}

func TabularizeTotalQuota(tbl table.Writer, tq *quota.TotalQuota) {
	tbl.AppendHeader(table.Row{"CPU Request Total", "Mem Request Total", "CPU Limit Total", "Mem Limit Total"})

	wq := tq.Sum()

	tbl.AppendRow([]interface{}{formatMilliCores(wq.Request.CPU), formatBytes(&wq.Request.Mem), formatMilliCores(wq.Limit.CPU),
		formatBytes(&wq.Limit.Mem)})
}

func TabularizeKubeQuota(tbl table.Writer, kq *quota.KubeQuota) {
	headers := make(table.Row, 0)
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

	tbl.AppendHeader(headers)

	tbl.AppendRow(values)
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
