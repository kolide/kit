// package instrumentation providies utilities for instrumenting Go code.
package instrumentation

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// NewNopCensusExporter creates a NoOp exporter for the OpenCensus Trace package.
// This exporter can be used for tests/local development and does not require the user to provide
// authentication to a remote tracing API.
func NewNopCensusExporter() *exporter { return &exporter{} }

type exporter struct{}

// ExportView logs the view data.
func (e *exporter) ExportView(vd *view.Data) {}

// ExportSpan logs the trace span.
func (e *exporter) ExportSpan(vd *trace.SpanData) {}
