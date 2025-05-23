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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	actionsv1alpha1 "github.com/odigos-io/odigos/api/actions/v1alpha1"
	applyconfigurationactionsv1alpha1 "github.com/odigos-io/odigos/api/generated/actions/applyconfiguration/actions/v1alpha1"
	scheme "github.com/odigos-io/odigos/api/generated/actions/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ErrorSamplersGetter has a method to return a ErrorSamplerInterface.
// A group's client should implement this interface.
type ErrorSamplersGetter interface {
	ErrorSamplers(namespace string) ErrorSamplerInterface
}

// ErrorSamplerInterface has methods to work with ErrorSampler resources.
type ErrorSamplerInterface interface {
	Create(ctx context.Context, errorSampler *actionsv1alpha1.ErrorSampler, opts v1.CreateOptions) (*actionsv1alpha1.ErrorSampler, error)
	Update(ctx context.Context, errorSampler *actionsv1alpha1.ErrorSampler, opts v1.UpdateOptions) (*actionsv1alpha1.ErrorSampler, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, errorSampler *actionsv1alpha1.ErrorSampler, opts v1.UpdateOptions) (*actionsv1alpha1.ErrorSampler, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*actionsv1alpha1.ErrorSampler, error)
	List(ctx context.Context, opts v1.ListOptions) (*actionsv1alpha1.ErrorSamplerList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *actionsv1alpha1.ErrorSampler, err error)
	Apply(ctx context.Context, errorSampler *applyconfigurationactionsv1alpha1.ErrorSamplerApplyConfiguration, opts v1.ApplyOptions) (result *actionsv1alpha1.ErrorSampler, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, errorSampler *applyconfigurationactionsv1alpha1.ErrorSamplerApplyConfiguration, opts v1.ApplyOptions) (result *actionsv1alpha1.ErrorSampler, err error)
	ErrorSamplerExpansion
}

// errorSamplers implements ErrorSamplerInterface
type errorSamplers struct {
	*gentype.ClientWithListAndApply[*actionsv1alpha1.ErrorSampler, *actionsv1alpha1.ErrorSamplerList, *applyconfigurationactionsv1alpha1.ErrorSamplerApplyConfiguration]
}

// newErrorSamplers returns a ErrorSamplers
func newErrorSamplers(c *ActionsV1alpha1Client, namespace string) *errorSamplers {
	return &errorSamplers{
		gentype.NewClientWithListAndApply[*actionsv1alpha1.ErrorSampler, *actionsv1alpha1.ErrorSamplerList, *applyconfigurationactionsv1alpha1.ErrorSamplerApplyConfiguration](
			"errorsamplers",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *actionsv1alpha1.ErrorSampler { return &actionsv1alpha1.ErrorSampler{} },
			func() *actionsv1alpha1.ErrorSamplerList { return &actionsv1alpha1.ErrorSamplerList{} },
		),
	}
}
