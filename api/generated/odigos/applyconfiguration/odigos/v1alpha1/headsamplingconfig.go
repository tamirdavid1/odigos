/*
Copyright 2022.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// HeadSamplingConfigApplyConfiguration represents an declarative configuration of the HeadSamplingConfig type for use
// with apply.
type HeadSamplingConfigApplyConfiguration struct {
	AttributesAndSamplerRules []AttributesAndSamplerRuleApplyConfiguration `json:"rules,omitempty"`
	FallbackFraction          *float64                                     `json:"fallbackFraction,omitempty"`
}

// HeadSamplingConfigApplyConfiguration constructs an declarative configuration of the HeadSamplingConfig type for use with
// apply.
func HeadSamplingConfig() *HeadSamplingConfigApplyConfiguration {
	return &HeadSamplingConfigApplyConfiguration{}
}

// WithAttributesAndSamplerRules adds the given value to the AttributesAndSamplerRules field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AttributesAndSamplerRules field.
func (b *HeadSamplingConfigApplyConfiguration) WithAttributesAndSamplerRules(values ...*AttributesAndSamplerRuleApplyConfiguration) *HeadSamplingConfigApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithAttributesAndSamplerRules")
		}
		b.AttributesAndSamplerRules = append(b.AttributesAndSamplerRules, *values[i])
	}
	return b
}

// WithFallbackFraction sets the FallbackFraction field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FallbackFraction field is set to the value of the last call.
func (b *HeadSamplingConfigApplyConfiguration) WithFallbackFraction(value float64) *HeadSamplingConfigApplyConfiguration {
	b.FallbackFraction = &value
	return b
}
