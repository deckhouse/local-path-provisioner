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

package templatetests

import "testing"

func TestStorageClassCarriesManagedByLabel(t *testing.T) {
	objects := render(t, map[string]any{
		"name": "localpath",
		"spec": map[string]any{"path": "/opt/local-path-provisioner", "reclaimPolicy": "Delete"},
	})

	sc := find(objects, "StorageClass", "localpath")
	if sc == nil {
		t.Fatal("StorageClass localpath is not rendered")
	}
	labels, _ := field(sc, "metadata", "labels").(map[string]any)
	if got := labels["storage.deckhouse.io/managed-by"]; got != "local-path-provisioner" {
		t.Fatalf("storage.deckhouse.io/managed-by is %v, want local-path-provisioner", got)
	}
	if got := labels["app"]; got != "local-path-provisioner" {
		t.Fatalf("app is %v, want local-path-provisioner", got)
	}
}
