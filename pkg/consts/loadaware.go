// Copyright 2022 The Katalyst Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consts

const (
	PodAnnotationLoadAwareEnableTrue  = "true"
	PodAnnotationLoadAwareEnableFalse = "false"

	DefaultPodAnnotationEnableLoadAware = "katalyst.kubewharf.io/enable_load_aware"
)

const (
	// Usage5MinAvgKey ...
	Usage5MinAvgKey = "avg_5min"
	// Usage15MinAvgKey ...
	Usage15MinAvgKey = "avg_15min"
	// Usage1HourMaxKey ...
	Usage1HourMaxKey = "max_1hour"
	// Usage1DayMaxKey ...
	Usage1DayMaxKey = "max_1day"
	// Annotation15MinAvgKey ...
	Annotation15MinAvgKey = "avg_15min_metadata"
	// Annotation1HourMaxKey ...
	Annotation1HourMaxKey = "max_1hour_metadata"
	// Annotation1DayMaxKey ...
	Annotation1DayMaxKey = "max_1day_metadata"
)
