package cli

import (
	"errors"
	"os"

	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/aauren/kube-quota/pkg/unit"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	allHeaderOrder = []string{quota.HeaderCPUReq, quota.HeaderMemReq, quota.HeaderCPULim, quota.HeaderMemLim,
		quota.HeaderEphemeralStorageReq, quota.HeaderEphemeralStorageLim}
)

type TableHeaderer interface {
	TableHeader() []string
}

type OrderedTableWriter interface {
	table.Writer
	OrderedHeaders() []string
}

type HeaderValuer interface {
	ValueForHeader(string) (unit.UnitWriter, error)
}

type TableWriterHeaderTracker struct {
	*table.Table
	uniqueHeaders  map[string]bool
	orderedHeaders []string
}

func (t *TableWriterHeaderTracker) Render() string {
	// For tables with less than 5 rows disable the alternating row color which is overall distracting for small tables
	if t.Length() < 5 {
		t.Table.Style().Color.RowAlternate = table.ColorOptionsBright.Row
	}

	return t.Table.Render()
}

func (t *TableWriterHeaderTracker) OrderedHeaders() []string {
	return t.orderedHeaders
}

func CreateTableWriter() *TableWriterHeaderTracker {
	tbl := table.Table{}

	// Set output to be STDOUT
	tbl.SetOutputMirror(os.Stdout)

	// Set basic overall formatting to StyleColoredBright
	tbl.SetStyle(table.StyleColoredBright)

	// Change StyleColoredBright to output with borders and column separators
	tbl.Style().Options = table.OptionsDefault

	// Set OptionsDefault to also include row separators
	tbl.Style().Options.SeparateRows = true

	// Set the title to be aligned in the center
	tbl.Style().Title.Align = text.AlignCenter

	// Wrap the table in the TableWriterHeaderTracker so that we can track the headers that we added to the table for future reference
	twht := TableWriterHeaderTracker{
		Table:          &tbl,
		uniqueHeaders:  make(map[string]bool, 0),
		orderedHeaders: make([]string, 0),
	}

	return &twht
}

//nolint:unused // We don't care if this function is currently unused
func configureColumns(tbl *TableWriterHeaderTracker) {
	clmCfg := make([]table.ColumnConfig, len(tbl.OrderedHeaders()))

	// Set the same config options on all columns on the table.
	for i := range tbl.OrderedHeaders() {
		clmCfg[i] = table.ColumnConfig{Number: i, AutoMerge: true}
	}

	tbl.SetColumnConfigs(clmCfg)
}

func AddTableHeader(tbl *TableWriterHeaderTracker, headerPrefixs []string, headerGenerators ...TableHeaderer) {

	// Create a list of unique headers that are used by all header generators that are passed
	for _, hdrGen := range headerGenerators {
		for _, hdr := range hdrGen.TableHeader() {
			tbl.uniqueHeaders[hdr] = true
		}
	}

	// Add custom header prefixes if there are any
	if len(headerPrefixs) > 0 {
		tbl.orderedHeaders = append(tbl.orderedHeaders, headerPrefixs...)
	}

	// Loop over all known headers in order and ensure that we have a consistent order with the headers we know we have
	for _, hdr := range allHeaderOrder {
		if _, ok := tbl.uniqueHeaders[hdr]; ok {
			tbl.orderedHeaders = append(tbl.orderedHeaders, hdr)
		}
	}

	headerRow := make(table.Row, len(tbl.orderedHeaders))
	for i, hdr := range tbl.orderedHeaders {
		headerRow[i] = table.Row{hdr}
	}

	tbl.AppendHeader(headerRow)

	// Set the column configs to automerge and align center
	// configureColumns(tbl)
}

func AddRow(tbl OrderedTableWriter, hv HeaderValuer, prefixes []string) error {
	values := make([]interface{}, len(tbl.OrderedHeaders())+len(prefixes)-1)

	if len(prefixes) > 0 {
		for i, pref := range prefixes {
			values[i] = pref
		}
	}

	for i, hdr := range tbl.OrderedHeaders()[len(prefixes):] {
		var err error
		notFound := quota.NoValueForHeader{Header: hdr}
		values[i+len(prefixes)], err = hv.ValueForHeader(hdr)
		if err != nil {
			if errors.Is(err, &notFound) {
				values[i] = ""
				continue
			}
			return err
		}
	}

	tbl.AppendRow(values)

	return nil
}
