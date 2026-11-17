---
title: "Local-path-provisioner module: examples"
description: "Examples of LocalPathProvisioner manifests for the local-path-provisioner module of Deckhouse Platform that keep or remove the volume data."
weight: 52
---

## StorageClass that keeps the volume data

The following LocalPathProvisioner creates the StorageClass `localpath-system` for the nodes of the `system` node group.
Volumes are created in `/opt/local-path-provisioner` on these nodes.
After a PersistentVolumeClaim is deleted, its PersistentVolume and the data in the directory stay, because `reclaimPolicy` is `Retain` by default:

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

## StorageClass that removes the volume data

The following LocalPathProvisioner sets `reclaimPolicy` to `Delete`.
After a PersistentVolumeClaim is deleted, the module removes its PersistentVolume and the directory with the data:

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
