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

import (
	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	common "github.com/odigos-io/odigos/common"
)

// RuntimeDetailsByContainerApplyConfiguration represents a declarative configuration of the RuntimeDetailsByContainer type for use
// with apply.
type RuntimeDetailsByContainerApplyConfiguration struct {
	ContainerName         *string                         `json:"containerName,omitempty"`
	Language              *common.ProgrammingLanguage     `json:"language,omitempty"`
	RuntimeVersion        *string                         `json:"runtimeVersion,omitempty"`
	EnvVars               []EnvVarApplyConfiguration      `json:"envVars,omitempty"`
	OtherAgent            *OtherAgentApplyConfiguration   `json:"otherAgent,omitempty"`
	LibCType              *common.LibCType                `json:"libCType,omitempty"`
	CriErrorMessage       *string                         `json:"criErrorMessage,omitempty"`
	EnvVarsFromDockerFile []EnvVarApplyConfiguration      `json:"envVarsFromDockerFile,omitempty"`
	RuntimeUpdateState    *odigosv1alpha1.ProcessingState `json:"runtimeUpdateState,omitempty"`
}

// RuntimeDetailsByContainerApplyConfiguration constructs a declarative configuration of the RuntimeDetailsByContainer type for use with
// apply.
func RuntimeDetailsByContainer() *RuntimeDetailsByContainerApplyConfiguration {
	return &RuntimeDetailsByContainerApplyConfiguration{}
}

// WithContainerName sets the ContainerName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ContainerName field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithContainerName(value string) *RuntimeDetailsByContainerApplyConfiguration {
	b.ContainerName = &value
	return b
}

// WithLanguage sets the Language field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Language field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithLanguage(value common.ProgrammingLanguage) *RuntimeDetailsByContainerApplyConfiguration {
	b.Language = &value
	return b
}

// WithRuntimeVersion sets the RuntimeVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RuntimeVersion field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithRuntimeVersion(value string) *RuntimeDetailsByContainerApplyConfiguration {
	b.RuntimeVersion = &value
	return b
}

// WithEnvVars adds the given value to the EnvVars field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EnvVars field.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithEnvVars(values ...*EnvVarApplyConfiguration) *RuntimeDetailsByContainerApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEnvVars")
		}
		b.EnvVars = append(b.EnvVars, *values[i])
	}
	return b
}

// WithOtherAgent sets the OtherAgent field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OtherAgent field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithOtherAgent(value *OtherAgentApplyConfiguration) *RuntimeDetailsByContainerApplyConfiguration {
	b.OtherAgent = value
	return b
}

// WithLibCType sets the LibCType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LibCType field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithLibCType(value common.LibCType) *RuntimeDetailsByContainerApplyConfiguration {
	b.LibCType = &value
	return b
}

// WithCriErrorMessage sets the CriErrorMessage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CriErrorMessage field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithCriErrorMessage(value string) *RuntimeDetailsByContainerApplyConfiguration {
	b.CriErrorMessage = &value
	return b
}

// WithEnvVarsFromDockerFile adds the given value to the EnvVarsFromDockerFile field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EnvVarsFromDockerFile field.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithEnvVarsFromDockerFile(values ...*EnvVarApplyConfiguration) *RuntimeDetailsByContainerApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEnvVarsFromDockerFile")
		}
		b.EnvVarsFromDockerFile = append(b.EnvVarsFromDockerFile, *values[i])
	}
	return b
}

// WithRuntimeUpdateState sets the RuntimeUpdateState field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RuntimeUpdateState field is set to the value of the last call.
func (b *RuntimeDetailsByContainerApplyConfiguration) WithRuntimeUpdateState(value odigosv1alpha1.ProcessingState) *RuntimeDetailsByContainerApplyConfiguration {
	b.RuntimeUpdateState = &value
	return b
}
