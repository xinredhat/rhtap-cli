// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
)

// FileReferenceSourceApplyConfiguration represents a declarative configuration of the FileReferenceSource type for use
// with apply.
type FileReferenceSourceApplyConfiguration struct {
	From      *operatorv1.SourceType                    `json:"from,omitempty"`
	ConfigMap *ConfigMapFileReferenceApplyConfiguration `json:"configMap,omitempty"`
}

// FileReferenceSourceApplyConfiguration constructs a declarative configuration of the FileReferenceSource type for use with
// apply.
func FileReferenceSource() *FileReferenceSourceApplyConfiguration {
	return &FileReferenceSourceApplyConfiguration{}
}

// WithFrom sets the From field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the From field is set to the value of the last call.
func (b *FileReferenceSourceApplyConfiguration) WithFrom(value operatorv1.SourceType) *FileReferenceSourceApplyConfiguration {
	b.From = &value
	return b
}

// WithConfigMap sets the ConfigMap field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ConfigMap field is set to the value of the last call.
func (b *FileReferenceSourceApplyConfiguration) WithConfigMap(value *ConfigMapFileReferenceApplyConfiguration) *FileReferenceSourceApplyConfiguration {
	b.ConfigMap = value
	return b
}
