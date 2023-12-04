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
	AggregateFunctionP90    = "_agg_p90"
	AggregateFunctionLatest = "_agg_latest"
)
