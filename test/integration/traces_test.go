//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goccy/go-json"
	"github.com/grafana/ebpf-autoinstrument/test/integration/components/jaeger"
	"github.com/mariomac/guara/pkg/test"
	"github.com/stretchr/testify/require"
)

func testHTTPTraces(t *testing.T) {
	doHTTPGet(t, instrumentedServiceStdURL+"/create-trace?delay=10ms", 200)

	var trace jaeger.Trace
	test.Eventually(t, testTimeout, func(t require.TestingT) {
		resp, err := http.Get(jaegerQueryURL + "?service=testserver&operation=GET%20%2Fcreate-trace")
		require.NoError(t, err)
		if resp == nil {
			return
		}
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var tq jaeger.TracesQuery
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&tq))
		traces := tq.FindBySpan(jaeger.Tag{Key: "http.target", Type: "string", Value: "/create-trace"})
		if len(traces) == 1 {
			trace = traces[0]
		} else {
			require.Failf(t, "error getting traces", "expected one trace, got: %+v", traces)
		}
	}, test.Interval(100*time.Millisecond))

	// Check the information of the parent span
	res := trace.FindByOperationName("GET /create-trace")
	require.Len(t, res, 1)
	parent := res[0]
	require.NotEmpty(t, parent.TraceID)
	require.NotEmpty(t, parent.SpanID)
	// check duration is at least 10ms
	assert.Less(t, 10*time.Millisecond, parent.Duration)
	// check span attributes
	assert.Truef(t, parent.AllMatches(
		jaeger.Tag{Key: "otel.library.name", Type: "string", Value: "github.com/grafana/ebpf-autoinstrument"},
		jaeger.Tag{Key: "http.method", Type: "string", Value: "GET"},
		jaeger.Tag{Key: "http.status_code", Type: "int64", Value: float64(200)},
		jaeger.Tag{Key: "http.target", Type: "string", Value: "/create-trace"},
		jaeger.Tag{Key: "net.host.port", Type: "int64", Value: float64(8080)},
		jaeger.Tag{Key: "http.route", Type: "string", Value: "/create-trace"},
		jaeger.Tag{Key: "span.kind", Type: "string", Value: "server"},
	), "not all tags matched in %+v", parent.Tags)

	// Check the information of the "in queue" span
	res = trace.FindByOperationName("in queue")
	require.Len(t, res, 1)
	queue := res[0]
	// Check parenthood
	p, ok := trace.ParentOf(&queue)
	require.True(t, ok)
	assert.Equal(t, parent.TraceID, p.TraceID)
	assert.Equal(t, parent.SpanID, p.SpanID)
	// check reasonable start and end times
	assert.GreaterOrEqual(t, queue.StartTime, parent.StartTime)
	assert.LessOrEqual(t, queue.StartTime+queue.Duration, parent.StartTime+parent.Duration)
	// check span attributes
	assert.Truef(t, queue.AllMatches(
		jaeger.Tag{Key: "otel.library.name", Type: "string", Value: "github.com/grafana/ebpf-autoinstrument"},
		jaeger.Tag{Key: "span.kind", Type: "string", Value: "internal"},
	), "not all tags matched in %+v", queue.Tags)

	// Check the information of the "processing" span
	res = trace.FindByOperationName("processing")
	require.Len(t, res, 1)
	processing := res[0]
	// Check parenthood
	p, ok = trace.ParentOf(&queue)
	require.True(t, ok)
	assert.Equal(t, parent.TraceID, p.TraceID)
	assert.Equal(t, parent.SpanID, p.SpanID)
	// check reasonable start and end times
	assert.GreaterOrEqual(t, processing.StartTime, queue.StartTime+queue.Duration)
	assert.LessOrEqual(t, processing.StartTime+processing.Duration, parent.StartTime+parent.Duration)
	// check span attributes
	assert.Truef(t, queue.AllMatches(
		jaeger.Tag{Key: "otel.library.name", Type: "string", Value: "github.com/grafana/ebpf-autoinstrument"},
		jaeger.Tag{Key: "span.kind", Type: "string", Value: "internal"},
	), "not all tags matched in %+v", queue.Tags)

	// check process ID
	require.Contains(t, trace.Processes, parent.ProcessID)
	assert.Equal(t, parent.ProcessID, queue.ProcessID)
	assert.Equal(t, parent.ProcessID, processing.ProcessID)
	process := trace.Processes[parent.ProcessID]
	assert.Equal(t, "testserver", process.ServiceName)
	assert.Truef(t, jaeger.AllMatches(process.Tags, []jaeger.Tag{
		{Key: "telemetry.sdk.language", Type: "string", Value: "go"},
		{Key: "service.namespace", Type: "string", Value: "integration-test"},
	}), "not all tags matched in %+v", process.Tags)

}