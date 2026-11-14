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
	"slices"
	"testing"

	"github.com/deckhouse/deckhouse/pkg/log"
	"github.com/deckhouse/module-sdk/pkg"
	"github.com/deckhouse/module-sdk/testing/mock"

	"github.com/deckhouse/local-path-provisioner/hooks/go/internal/testsnapshot"
)

var clusterStorageClasses = []string{
	`{"name":"localpath-worker","reclaimPolicy":"Retain"}`,
	`{"name":"localpath-worker2","reclaimPolicy":"Delete"}`,
}

func run(t *testing.T, crds ...string) []string {
	t.Helper()

	var deleted []string
	patches := mock.NewPatchCollectorMock(t)
	patches.DeleteMock.Optional().Set(func(apiVersion, kind, namespace, name string) {
		if apiVersion != "storage.k8s.io/v1" || kind != "StorageClass" || namespace != "" {
			t.Fatalf("unexpected delete of %s %s %s/%s", apiVersion, kind, namespace, name)
		}
		deleted = append(deleted, name)
	})

	input := &pkg.HookInput{
		Snapshots: testsnapshot.Snapshots{
			storageClassesSnapshot: clusterStorageClasses,
			crdsSnapshot:           crds,
		},
		PatchCollector: patches,
		Logger:         log.NewNop(),
	}
	if err := handle(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	return deleted
}

func TestFreshClusterDeletesNothing(t *testing.T) {
	if deleted := run(t); len(deleted) != 0 {
		t.Fatalf("deleted %v, want nothing", deleted)
	}
}

func TestDeletesStorageClassWithChangedReclaimPolicy(t *testing.T) {
	deleted := run(t, `{"name":"localpath-worker","reclaimPolicy":"Delete"}`)

	if !slices.Equal(deleted, []string{"localpath-worker"}) {
		t.Fatalf("deleted %v, want [localpath-worker]", deleted)
	}
}

func TestKeepsStorageClassWithSameReclaimPolicy(t *testing.T) {
	if deleted := run(t, `{"name":"localpath-worker2","reclaimPolicy":"Delete"}`); len(deleted) != 0 {
		t.Fatalf("deleted %v, want nothing", deleted)
	}
}
