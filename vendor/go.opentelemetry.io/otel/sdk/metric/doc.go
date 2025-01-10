// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package metric provides an implementation of the OpenTelemetry metrics SDK.
//
// See https://opentelemetry.io/docs/concepts/signals/metrics/ for information
// about the concept of OpenTelemetry metrics and
// https://opentelemetry.io/docs/concepts/components/ for more information
// about OpenTelemetry SDKs.
//
// The entry point for the metric package is the MeterProvider. It is the
// object that all API calls use to create Meters, instruments, and ultimately
// make metric measurements. Also, it is an object that should be used to
// control the life-cycle (start, flush, and shutdown) of the SDK.
//
// A MeterProvider needs to be configured to export the measured data, this is
// done by configuring it with a Reader implementation (using the WithReader
// MeterProviderOption). Readers take two forms: ones that push to an endpoint
// (NewPeriodicReader), and ones that an endpoint pulls from. See
// [go.opentelemetry.io/otel/exporters] for exporters that can be used as
// or with these Readers.
//
// Each Reader, when registered with the MeterProvider, can be augmented with a
// View. Views allow users that run OpenTelemetry instrumented code to modify
// the generated data of that instrumentation.
//
// The data generated by a MeterProvider needs to include information about its
// origin. A MeterProvider needs to be configured with a Resource, using the
// WithResource MeterProviderOption, to include this information. This Resource
// should be used to describe the unique runtime environment instrumented code
// is being run on. That way when multiple instances of the code are collected
// at a single endpoint their origin is decipherable.
//
// To avoid leaking memory, the SDK returns the same instrument for calls to
// create new instruments with the same Name, Unit, and Description.
// Importantly, callbacks provided using metric.WithFloat64Callback or
// metric.WithInt64Callback will only apply for the first instrument created
// with a given Name, Unit, and Description. Instead, use
// Meter.RegisterCallback and Registration.Unregister to add and remove
// callbacks without leaking memory.
//
// See [go.opentelemetry.io/otel/metric] for more information about
// the metric API.
//
// See [go.opentelemetry.io/otel/sdk/metric/internal/x] for information about
// the experimental features.
package metric
