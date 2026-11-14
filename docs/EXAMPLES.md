---
title: "The local-path-provisioner module: examples"
description: "Examples of LocalPathProvisioner resources that create StorageClasses for local volumes on Deckhouse nodes."
weight: 52
---

## Example of a `LocalPathProvisioner` custom resource

Reclaim policy set by default to `Retain`.

```yaml
apiVersion: deckhouse.io/v1alpha1
kind: LocalPathProvisioner
metadata:
  name: localpath-system
spec:
  nodeGroups:
  - system
  path: "/opt/local-path-provisioner"
```

## Example of a `LocalPathProvisioner` custom resource with `reclaimPolicy` set

Reclaim policy set to `Delete`.

```yaml
apiVersion: deckhouse.io/v1alpha1
kind: LocalPathProvisioner
metadata:
  name: localpath-system
spec:
  nodeGroups:
  - system
  path: "/opt/local-path-provisioner"
  reclaimPolicy: "Delete"
```
