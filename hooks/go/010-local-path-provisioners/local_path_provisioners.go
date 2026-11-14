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

package localpathprovisioners

import (
	"context"
	"fmt"

	"github.com/deckhouse/module-sdk/pkg"
	objectpatch "github.com/deckhouse/module-sdk/pkg/object-patch"
	"github.com/deckhouse/module-sdk/pkg/registry"

	"github.com/deckhouse/local-path-provisioner/hooks/go/internal/v1alpha1"
)

const (
	snapshotName = "lpp"
	valuesPath   = "localPathProvisioner.internal.localPathProvisioners"
)

type LocalPathProvisionerInfo struct {
	Name string                            `json:"name"`
	Spec v1alpha1.LocalPathProvisionerSpec `json:"spec"`
}

var _ = registry.RegisterFunc(&pkg.HookConfig{
	Queue: "/modules/local-path-provisioner",
	Kubernetes: []pkg.KubernetesConfig{
		{
			Name:       snapshotName,
			APIVersion: "deckhouse.io/v1alpha1",
			Kind:       "LocalPathProvisioner",
			JqFilter:   `{name: .metadata.name, spec: .spec}`,
		},
	},
}, handle)

func handle(_ context.Context, input *pkg.HookInput) error {
	localPathProvisioners, err := objectpatch.UnmarshalToStruct[LocalPathProvisionerInfo](input.Snapshots, snapshotName)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s snapshot: %w", snapshotName, err)
	}
	if len(localPathProvisioners) == 0 {
		localPathProvisioners = make([]LocalPathProvisionerInfo, 0)
	}

	input.Values.Set(valuesPath, localPathProvisioners)
	return nil
}
