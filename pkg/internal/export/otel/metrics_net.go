package otel

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/mariomac/pipes/pipe"
	"go.opentelemetry.io/otel/attribute"
	metric2 "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"

	"github.com/grafana/beyla/pkg/internal/export/attributes"
	"github.com/grafana/beyla/pkg/internal/export/expire"
	"github.com/grafana/beyla/pkg/internal/netolly/ebpf"
	"github.com/grafana/beyla/pkg/internal/pipe/global"
)

// NetMetricsConfig extends MetricsConfig for Network Metrics
type NetMetricsConfig struct {
	Metrics            *MetricsConfig
	AttributeSelectors attributes.Selection
}

func (mc NetMetricsConfig) Enabled() bool {
	return mc.Metrics != nil && mc.Metrics.EndpointEnabled() && slices.Contains(mc.Metrics.Features, FeatureNetwork)
}

func nmlog() *slog.Logger {
	return slog.With("component", "otel.NetworkMetricsExporter")
}

func newResource() *resource.Resource {
	attrs := []attribute.KeyValue{
		semconv.ServiceName("beyla-network-flows"),
		semconv.ServiceInstanceID(uuid.New().String()),
		// SpanMetrics requires an extra attribute besides service name
		// to generate the traces_target_info metric,
		// so the service is visible in the ServicesList
		// This attribute also allows that App O11y plugin shows this app as a Go application.
		semconv.TelemetrySDKLanguageKey.String(semconv.TelemetrySDKLanguageGo.Value.AsString()),
		// We set the SDK name as Beyla, so we can distinguish beyla generated metrics from other SDKs
		semconv.TelemetrySDKNameKey.String("beyla"),
	}

	return resource.NewWithAttributes(semconv.SchemaURL, attrs...)
}

func newMeterProvider(res *resource.Resource, exporter *metric.Exporter, interval time.Duration) (*metric.MeterProvider, error) {
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(*exporter, metric.WithInterval(interval))),
	)
	return meterProvider, nil
}

type netMetricsExporter struct {
	metrics *Expirer[*ebpf.Record, metric2.Int64Observer, *IntCounter, int64]
	clock   *expire.CachedClock
}

func NetMetricsExporterProvider(ctxInfo *global.ContextInfo, cfg *NetMetricsConfig) (pipe.FinalFunc[[]*ebpf.Record], error) {
	if !cfg.Enabled() {
		// This node is not going to be instantiated. Let the pipes library just ignore it.
		return pipe.IgnoreFinal[[]*ebpf.Record](), nil
	}
	log := nmlog()
	log.Debug("instantiating network metrics exporter provider")
	exporter, err := InstantiateMetricsExporter(context.Background(), cfg.Metrics, log)
	if err != nil {
		log.Error("", "error", err)
		return nil, err
	}

	provider, err := newMeterProvider(newResource(), &exporter, cfg.Metrics.Interval)

	if err != nil {
		log.Error("", "error", err)
		return nil, err
	}

	attrProv, err := attributes.NewAttrSelector(ctxInfo.MetricAttributeGroups, cfg.AttributeSelectors)
	if err != nil {
		return nil, fmt.Errorf("network OTEL exporter attributes enable: %w", err)
	}
	attrs := attributes.OpenTelemetryGetters(
		ebpf.RecordGetters,
		attrProv.For(attributes.BeylaNetworkFlow))

	clock := expire.NewCachedClock(timeNow)
	expirer := NewExpirer[*ebpf.Record, metric2.Int64Observer](NewIntCounter, attrs, clock.Time, cfg.Metrics.TTL)
	ebpfEvents := provider.Meter("network_ebpf_events")

	_, err = ebpfEvents.Int64ObservableCounter(
		attributes.BeylaNetworkFlow.OTEL,
		metric2.WithDescription("total bytes_sent value of network flows observed by probe since its launch"),
		metric2.WithUnit("{bytes}"),
		metric2.WithInt64Callback(expirer.Collect),
	)
	if err != nil {
		log.Error("creating observable counter", "error", err)
		return nil, err
	}
	log.Debug("restricting attributes not in this list", "attributes", cfg.AttributeSelectors)
	return (&netMetricsExporter{
		metrics: expirer,
		clock:   clock,
	}).Do, nil
}

func (me *netMetricsExporter) Do(in <-chan []*ebpf.Record) {
	for i := range in {
		me.clock.Update()
		for _, v := range i {
			me.metrics.ForRecord(v).Add(int64(v.Metrics.Bytes))
		}
	}
}
