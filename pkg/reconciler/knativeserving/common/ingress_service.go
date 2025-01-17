/*
Copyright 2020 The Knative Authors

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

package common

import (
	"fmt"
	"strings"

	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"knative.dev/operator/pkg/apis/operator/v1beta1"
)

// IngressServiceTransform pins the namespace to istio-system for the service named knative-local-gateway.
// It also removes the OwnerReference to the operator, as they are in different namespaces, which is
// invalid in Kubernetes 1.20+.
func IngressServiceTransform(ks *v1beta1.KnativeServing) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetAPIVersion() == "v1" && u.GetKind() == "Service" && u.GetName() == "knative-local-gateway" {
			u.SetNamespace("istio-system")
			u.SetOwnerReferences(nil)
			config := ks.GetSpec().GetConfig()
			if data, ok := config["istio"]; ok {
				UpdateNamespace(u, data, ks.GetNamespace())
			}

			// The "config-" prefix is optional
			if data, ok := config["config-istio"]; ok {
				UpdateNamespace(u, data, ks.GetNamespace())
			}
		}
		return nil
	}
}

// UpdateNamespace set correct namespace of istio to the service knative-local-gateway
func UpdateNamespace(u *unstructured.Unstructured, data map[string]string, ns string) {
	key := fmt.Sprintf("%s.%s.%s", "local-gateway", ns, "knative-local-gateway")
	if val, ok := data[key]; ok {
		fields := strings.Split(val, ".")
		// The value is in the format of knative-local-gateway.{istio-namespace}.svc.cluster.local
		// The second item is the istio namespace
		if len(fields) >= 2 {
			u.SetNamespace(fields[1])
		}
	}
}
