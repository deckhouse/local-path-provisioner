---
title: "The local-path-provisioner module: configuration"
force_searchable: true
description: "How to enable the local-path-provisioner module; the module has no settings."
weight: 30
---

The module has no settings.

## Enabling the module

The module is enabled by default in the `Default` and `Managed` bundles. In other bundles, enable it with a ModuleConfig:

```yaml
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: local-path-provisioner
spec:
  enabled: true
  version: 1
```

Wait for the module to become `Ready`:

```shell
d8 k get module local-path-provisioner -w
```

To disable the module, set `spec.enabled: false` in the same ModuleConfig.

## Moving from the embedded module

Before Deckhouse 1.79 the module was part of Deckhouse. On the update to 1.79, a cluster that runs the embedded module switches to this one automatically: the module is downloaded from the `deckhouse` ModuleSource, and the update waits until it is available there. The `LocalPathProvisioner` resources, the StorageClasses and the PersistentVolumes are kept. The module Pods restart once, because the image and the image pull Secret change; volumes are not created or removed until the new Pods are `Ready`.
