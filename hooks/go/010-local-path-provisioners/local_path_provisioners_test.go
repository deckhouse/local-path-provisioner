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
	"encoding/json"
	"testing"

	"github.com/deckhouse/module-sdk/pkg"
	"github.com/deckhouse/module-sdk/testing/mock"

	"github.com/deckhouse/local-path-provisioner/hooks/go/internal/testsnapshot"
)

func run(t *testing.T, objects ...string) string {
	t.Helper()

	var got string
	values := mock.NewOutputPatchableValuesCollectorMock(t)
	values.SetMock.Set(func(path string, value any) {
		if path != valuesPath {
			t.Fatalf("unexpected values path %q", path)
		}
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		got = string(raw)
	})

	input := &pkg.HookInput{
		Snapshots: testsnapshot.Snapshots{snapshotName: objects},
		Values:    values,
	}
	if err := handle(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	return got
}

func TestFillsInternalValues(t *testing.T) {
	got := run(t, `{"name":"local1","spec":{"nodeGroups":["master"],"path":"/local","reclaimPolicy":"Retain","unknown":"dropped"}}`)

	want := `[{"name":"local1","spec":{"nodeGroups":["master"],"path":"/local","reclaimPolicy":"Retain"}}]`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestEmptyArrayWithoutObjects(t *testing.T) {
	if got := run(t); got != "[]" {
		t.Fatalf("got %s, want []", got)
	}
}
