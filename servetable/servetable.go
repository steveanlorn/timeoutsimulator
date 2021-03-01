// Package servetable generate report of Result in table
package servetable

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/steveanlorn/timeoutsimulator"
)

const (
	tableHeaderName    = "NAME"
	tableHeaderLatency = "LATENCY"
	tableHeaderTimeout = "IS TIMEOUT"
	tableHeaderTimeIn  = "TIME IN"
	tableHeaderTimeOut = "TIME OUT"
	tableHeaderRemark  = "NOTE"
)

type tableHeader struct {
	label    string
	dataType reflect.Kind
}

var tableHeaders = []tableHeader{
	{
		label:    tableHeaderName,
		dataType: reflect.String,
	},
	{
		label:    tableHeaderLatency,
		dataType: reflect.String,
	},
	{
		label:    tableHeaderTimeout,
		dataType: reflect.Bool,
	},
	{
		label:    tableHeaderTimeIn,
		dataType: reflect.String,
	},
	{
		label:    tableHeaderTimeOut,
		dataType: reflect.String,
	},
	{
		label:    tableHeaderRemark,
		dataType: reflect.String,
	},
}

const (
	legendSimulator = "SIMULATOR"
	legendTimeout   = "TIMEOUT DURATION"
)

const (
	separator1 = "================================"
	separator2 = "--------------------------------"
)

func generateRowFormat() string {
	b := strings.Builder{}
	for _, th := range tableHeaders {
		switch th.dataType {
		case reflect.String:
			b.WriteString("%s\t")
		case reflect.Bool:
			b.WriteString("%v\t")
		}
	}

	b.WriteString("\n")

	return b.String()
}

func generateTableHeader() string {
	b := strings.Builder{}
	for _, th := range tableHeaders {
		b.WriteString(th.label)
		b.WriteString("\t")
	}

	b.WriteString("\n")

	return b.String()
}

// Generate generates report of result to output.
func Generate(result timeoutsimulator.Result, output io.Writer) error {
	rowFormat := generateRowFormat()

	w := tabwriter.NewWriter(output, 0, 0, 1, ' ', tabwriter.Debug)

	_, _ = fmt.Fprintf(w, "%s\n", separator1)
	_, _ = fmt.Fprintf(w, "%s:%s\n", legendSimulator, result.Name)
	_, _ = fmt.Fprintf(w, "%s:%s\n", legendTimeout, result.TimeoutDuration.String())
	_, _ = fmt.Fprintf(w, "%s\n", separator2)

	_, _ = fmt.Fprint(w, generateTableHeader())
	for _, data := range result.Data {
		_, _ = fmt.Fprintf(w, rowFormat, data.Name, data.Latency.String(), data.IsDeadlineExceeded, data.TimeIn.String(), data.TimeOut.String(), data.Remark)
	}
	_, _ = fmt.Fprintf(w, "%s\n", separator1)

	return w.Flush()
}
