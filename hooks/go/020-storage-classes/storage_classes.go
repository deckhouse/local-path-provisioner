/*
Copyright 2026 Flant JSC

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

package storageclasses

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/deckhouse/module-sdk/pkg"
	objectpatch "github.com/deckhouse/module-sdk/pkg/object-patch"
	"github.com/deckhouse/module-sdk/pkg/registry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	storageClassesSnapshot = "module_storageclasses"
	crdsSnapshot           = "module_crds"
)

type StorageClass struct {
	Name          string `json:"name"`
	ReclaimPolicy string `json:"reclaimPolicy"`
}

var hookConfig = &pkg.HookConfig{
	Queue:        "/modules/local-path-provisioner",
	OnBeforeHelm: &pkg.OrderedConfig{Order: 20},
	Kubernetes: []pkg.KubernetesConfig{
		{
			Name:       storageClassesSnapshot,
			APIVersion: "storage.k8s.io/v1",
			Kind:       "StorageClass",
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"module": "local-path-provisioner"},
			},
			JqFilter: `{name: .metadata.name, reclaimPolicy: .reclaimPolicy}`,
		},
		{
			Name:       crdsSnapshot,
			APIVersion: "deckhouse.io/v1alpha1",
			Kind:       "LocalPathProvisioner",
			JqFilter:   `{name: .metadata.name, reclaimPolicy: .spec.reclaimPolicy}`,
		},
	},
}

var _ = registry.RegisterFunc(hookConfig, handle)

func handle(_ context.Context, input *pkg.HookInput) error {
	moduleStorageClasses, err := objectpatch.UnmarshalToStruct[StorageClass](input.Snapshots, storageClassesSnapshot)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s snapshot: %w", storageClassesSnapshot, err)
	}
	moduleCRDs, err := objectpatch.UnmarshalToStruct[StorageClass](input.Snapshots, crdsSnapshot)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s snapshot: %w", crdsSnapshot, err)
	}

	for _, name := range StorageClassesToDelete(moduleStorageClasses, moduleCRDs) {
		input.Logger.Info("Deleting storageclass because its parameters have been changed", slog.String("storage_class", name))
		input.PatchCollector.Delete("storage.k8s.io/v1", "StorageClass", "", name)
	}

	return nil
}

func StorageClassesToDelete(storageClasses, crds []StorageClass) []string {
	var names []string
	for _, crd := range crds {
		for _, storageClass := range storageClasses {
			if storageClass.Name == crd.Name {
				if storageClass.ReclaimPolicy != crd.ReclaimPolicy {
					names = append(names, storageClass.Name)
				}
				break
			}
		}
	}
	return names
}
