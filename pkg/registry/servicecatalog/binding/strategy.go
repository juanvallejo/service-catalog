/*
Copyright 2017 The Kubernetes Authors.

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

package binding

// this was copied from where else and edited to fit our objects

import (
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/validation/field"

	"github.com/golang/glog"
	sc "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog"
	scv "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/validation"
)

// implements interfaces RESTCreateStrategy, RESTUpdateStrategy, RESTDeleteStrategy
type bindingRESTStrategy struct {
	runtime.ObjectTyper // inherit ObjectKinds method
	kapi.NameGenerator  // GenerateName method for CreateStrategy
}

var (
	bindingRESTStrategies = bindingRESTStrategy{
		// embeds to pull in existing code behavior from upstream

		ObjectTyper: kapi.Scheme,
		// use the generator from upstream k8s, or implement method
		// `GenerateName(base string) string`
		NameGenerator: kapi.SimpleNameGenerator,
	}
	_ rest.RESTCreateStrategy = bindingRESTStrategies
	_ rest.RESTUpdateStrategy = bindingRESTStrategies
	_ rest.RESTDeleteStrategy = bindingRESTStrategies
)

// Canonicalize does not transform a binding.
func (bindingRESTStrategy) Canonicalize(obj runtime.Object) {
	_, ok := obj.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to create")
	}
}

// NamespaceScoped returns false as bindings are not scoped to a namespace.
func (bindingRESTStrategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate receives a the incoming Binding and clears it's
// Status. Status is not a user settable field.
func (bindingRESTStrategy) PrepareForCreate(ctx kapi.Context, obj runtime.Object) {
	binding, ok := obj.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to create")
	}
	// Is there anything to pull out of the context `ctx`?

	// Creating a brand new object, thus it must have no
	// status. We can't fail here if they passed a status in, so
	// we just wipe it clean.
	binding.Status = sc.BindingStatus{}
	// Fill in the first entry set to "creating"?
	binding.Status.Conditions = []sc.BindingCondition{}
}

func (bindingRESTStrategy) Validate(ctx kapi.Context, obj runtime.Object) field.ErrorList {
	return scv.ValidateBinding(obj.(*sc.Binding))
}

func (bindingRESTStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (bindingRESTStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (bindingRESTStrategy) PrepareForUpdate(ctx kapi.Context, new, old runtime.Object) {
	newBinding, ok := new.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to update to")
	}
	oldBinding, ok := old.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to update from")
	}
	newBinding.Spec = oldBinding.Spec
	newBinding.Status = oldBinding.Status
}

func (bindingRESTStrategy) ValidateUpdate(ctx kapi.Context, new, old runtime.Object) field.ErrorList {
	newBinding, ok := new.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to validate to")
	}
	oldBinding, ok := old.(*sc.Binding)
	if !ok {
		glog.Fatal("received a non-binding object to validate from")
	}

	return scv.ValidateBindingUpdate(newBinding, oldBinding)
}
