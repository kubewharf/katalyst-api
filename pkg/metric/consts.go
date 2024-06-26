/*
Copyright 2022 The Katalyst Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metric

// aggregated functions that'd supported by katalyst now,
// and client should add those as suffix in metric-name when referring to kcmas.
// for instance, `pod_cpu_load_1min_agg_max` means to return
// a single metric item to represent the max value of all collected items,
// and we will put the corresponding time-window in response.
const (
	AggregateFunctionAvg    = "_agg_avg"
	AggregateFunctionMax    = "_agg_max"
	AggregateFunctionMin    = "_agg_min"
	AggregateFunctionP99    = "_agg_p99"
	AggregateFunctionP95    = "_agg_p95"
	AggregateFunctionP90    = "_agg_p90"
	AggregateFunctionLatest = "_agg_latest"
)

// MetricSelectorKeyGroupBy is the key of groupBy in metric selector. It's value should be a set of the real metric
// selector keys which will be used to group the metrics. MetricSelectorKeyGroupBy should only be used in aggregated
// metrics.
// For example, if we want to get the max cpu load of each container,we can query the `pod_cpu_load_1min_agg_max` with
// following metric selector: `groupBy=container`.
const MetricSelectorKeyGroupBy = "groupBy"

// MetricNameSPDAggMetrics represents the metric name provided to the API Server when exposing
// the metrics in SPD in the form of Extermal Metric.
const MetricNameSPDAggMetrics = "spd_agg_metrics"

// MetricSelectorKeySPD represents a series of External Metric labels, used by the SPD Metric Store to
// filter and return specific workload metric data.
const (
	MetricSelectorKeySPDName          = "metric.katalyst.kubewharf.io/spd-name"
	MetricSelectorKeySPDResourceName  = "metric.katalyst.kubewharf.io/spd-resource-name"
	MetricSelectorKeySPDScopeName     = "metric.katalyst.kubewharf.io/spd-scope-name"
	MetricSelectorKeySPDContainerName = "metric.katalyst.kubewharf.io/spd-container-name"
)
