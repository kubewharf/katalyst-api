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

package utils

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/kubewharf/katalyst-api/pkg/consts"
)

const (
	// cnrConditionQoSLevelEffectSeparator separates CNR Condition effects into QoS Level and CNR effects.
	cnrConditionQoSLevelEffectSeparator = "/"
)

// ParseConditionEffect parses the condition effect into QoS level and CNR effect.
func ParseConditionEffect(effect string) (consts.QoSLevel, v1.TaintEffect, error) {
	parts := strings.Split(effect, cnrConditionQoSLevelEffectSeparator)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid condition effect: %s", effect)
	}
	return consts.QoSLevel(parts[0]), v1.TaintEffect(parts[1]), nil
}

// GenerateConditionEffect generates the condition effect from QoS level and CNR effect.
func GenerateConditionEffect(level consts.QoSLevel, effect v1.TaintEffect) string {
	return fmt.Sprintf("%s%s%s", level, cnrConditionQoSLevelEffectSeparator, effect)
}
