/*
Copyright The Kubernetes Authors.

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

	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsrbacv1alpha1 "k8s.io/client-go/applyconfigurations/rbac/v1alpha1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// RoleBindingsGetter has a method to return a RoleBindingInterface.
// A group's client should implement this interface.
type RoleBindingsGetter interface {
	RoleBindings(namespace string) RoleBindingInterface
}

// RoleBindingInterface has methods to work with RoleBinding resources.
type RoleBindingInterface interface {
	Create(ctx context.Context, roleBinding *rbacv1alpha1.RoleBinding, opts v1.CreateOptions) (*rbacv1alpha1.RoleBinding, error)
	Update(ctx context.Context, roleBinding *rbacv1alpha1.RoleBinding, opts v1.UpdateOptions) (*rbacv1alpha1.RoleBinding, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*rbacv1alpha1.RoleBinding, error)
	List(ctx context.Context, opts v1.ListOptions) (*rbacv1alpha1.RoleBindingList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *rbacv1alpha1.RoleBinding, err error)
	Apply(ctx context.Context, roleBinding *applyconfigurationsrbacv1alpha1.RoleBindingApplyConfiguration, opts v1.ApplyOptions) (result *rbacv1alpha1.RoleBinding, err error)
	RoleBindingExpansion
}

// roleBindings implements RoleBindingInterface
type roleBindings struct {
	*gentype.ClientWithListAndApply[*rbacv1alpha1.RoleBinding, *rbacv1alpha1.RoleBindingList, *applyconfigurationsrbacv1alpha1.RoleBindingApplyConfiguration]
}

// newRoleBindings returns a RoleBindings
func newRoleBindings(c *RbacV1alpha1Client, namespace string) *roleBindings {
	return &roleBindings{
		gentype.NewClientWithListAndApply[*rbacv1alpha1.RoleBinding, *rbacv1alpha1.RoleBindingList, *applyconfigurationsrbacv1alpha1.RoleBindingApplyConfiguration](
			"rolebindings",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *rbacv1alpha1.RoleBinding { return &rbacv1alpha1.RoleBinding{} },
			func() *rbacv1alpha1.RoleBindingList { return &rbacv1alpha1.RoleBindingList{} },
			gentype.PrefersProtobuf[*rbacv1alpha1.RoleBinding](),
		),
	}
}
