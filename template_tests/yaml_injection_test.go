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

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/yaml"
)

const namespace = "d8-local-path-provisioner"

type object = map[string]any

func render(t *testing.T, provisioners ...any) []object {
	t.Helper()

	chart, err := loader.LoadDir("..")
	if err != nil {
		t.Fatalf("load the chart: %v", err)
	}

	values := map[string]any{
		"global": map[string]any{
			"enabledModules": []any{"vertical-pod-autoscaler", "admission-policy-engine-crd"},
			"modulesImages": map[string]any{
				"registry": map[string]any{"base": "registry.example.com/deckhouse", "dockercfg": "e30="},
				"digests": map[string]any{
					"localPathProvisioner": map[string]any{
						"localPathProvisioner": "sha256:0000000000000000000000000000000000000000000000000000000000000001",
						"helper":               "sha256:0000000000000000000000000000000000000000000000000000000000000002",
					},
				},
			},
			"discovery": map[string]any{
				"d8SpecificNodeCountByRole": map[string]any{"master": 3},
			},
			"internal": map[string]any{
				"modules": map[string]any{"kubeRBACProxyCA": map[string]any{"cert": "CERT"}},
			},
		},
		"localPathProvisioner": map[string]any{
			"registry": map[string]any{"base": "registry.example.com/modules/local-path-provisioner", "dockercfg": "e30="},
			"internal": map[string]any{"localPathProvisioners": provisioners},
		},
	}

	renderValues, err := chartutil.ToRenderValues(chart, values, chartutil.ReleaseOptions{
		Name:      chart.Name(),
		Namespace: namespace,
		IsInstall: true,
	}, chartutil.DefaultCapabilities)
	if err != nil {
		t.Fatalf("build the render values: %v", err)
	}

	files, err := engine.Render(chart, renderValues)
	if err != nil {
		t.Fatalf("render the chart: %v", err)
	}

	var objects []object
	for name, content := range files {
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		for _, manifest := range releaseutil.SplitManifests(content) {
			var obj object
			if err := yaml.Unmarshal([]byte(manifest), &obj); err != nil {
				t.Fatalf("parse a document of %s: %v\n%s", name, err, manifest)
			}
			if obj != nil {
				objects = append(objects, obj)
			}
		}
	}
	return objects
}

func field(obj object, path ...string) any {
	var current any = obj
	for _, key := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[key]
	}
	return current
}

func find(objects []object, kind, name string) object {
	for _, obj := range objects {
		if obj["kind"] == kind && field(obj, "metadata", "name") == name {
			return obj
		}
	}
	return nil
}

func TestPathCannotInjectDocuments(t *testing.T) {
	for _, tc := range []struct{ name, payload string }{
		{"newline", "/opt/pwned\"]}\n---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: pwned\n  namespace: " + namespace},
		{"NEL", "/opt/pwned\u0085kind:Secret"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			objects := render(t, map[string]any{
				"name": "evil",
				"spec": map[string]any{"path": tc.payload, "reclaimPolicy": "Delete"},
			})

			for _, obj := range objects {
				if field(obj, "metadata", "name") == "pwned" {
					t.Fatalf("the payload rendered a document of its own: %v", obj)
				}
			}

			cm := find(objects, "ConfigMap", "local-path-provisioner-evil")
			if cm == nil {
				t.Fatal("ConfigMap local-path-provisioner-evil is not rendered")
			}
			var config struct {
				NodePathMap []struct {
					Paths []string `json:"paths"`
				} `json:"nodePathMap"`
			}
			rawConfig, _ := field(cm, "data", "config.json").(string)
			if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
				t.Fatalf("config.json is not valid JSON: %v\n%s", err, rawConfig)
			}
			if len(config.NodePathMap) == 0 || len(config.NodePathMap[0].Paths) != 1 || config.NodePathMap[0].Paths[0] != tc.payload {
				t.Fatalf("config.json carries %+v, want the single path %q", config.NodePathMap, tc.payload)
			}

			var helperPod object
			rawHelperPod, _ := field(cm, "data", "helperPod.yaml").(string)
			if err := yaml.Unmarshal([]byte(rawHelperPod), &helperPod); err != nil || helperPod["kind"] != "Pod" {
				t.Fatalf("helperPod.yaml is not a Pod (%v):\n%s", err, rawHelperPod)
			}

			spe := find(objects, "SecurityPolicyException", "helper-evil")
			if spe == nil {
				t.Fatal("SecurityPolicyException helper-evil is not rendered")
			}
			description, _ := field(spe, "spec", "securityContext", "runAsNonRoot", "metadata", "description").(string)
			if !strings.Contains(description, fmt.Sprintf("%q", tc.payload)) {
				t.Fatalf("the description %q does not quote the path %q", description, tc.payload)
			}
			allowed, _ := field(spe, "spec", "volumes", "hostPath", "allowedValues").([]any)
			if len(allowed) == 0 || field(allowed[0].(object), "path") != tc.payload+"/" {
				t.Fatalf("allowedValues %v, want the path %q", allowed, tc.payload+"/")
			}
		})
	}
}
