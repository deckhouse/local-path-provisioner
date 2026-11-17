---
title: "Local-path-provisioner module: FAQ"
description: "Answers to common questions about the local-path-provisioner module of Deckhouse Platform: Prometheus on local disks, changing the volume path and the update to Deckhouse Platform 1.79."
weight: 55
---

## How to store Prometheus data on local disks?

To keep the Prometheus data on the local disks of the nodes where Prometheus runs, do the following:

1. Create a LocalPathProvisioner.
   In `nodeGroups`, list the node group where the Prometheus Pods run; in `path`, set the directory on the nodes for the data:

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

1. In the settings of the `prometheus` module, set [`storageClass`](/modules/prometheus/configuration.html#parameters-storageclass)
   and [`longtermStorageClass`](/modules/prometheus/configuration.html#parameters-longtermstorageclass) to the name of the new StorageClass:

   ```yaml
   storageClass: localpath-system
   longtermStorageClass: localpath-system
   ```

1. Wait until the Prometheus Pods restart.

## What happens to existing volumes when the volume path changes?

New volumes are created in the new `path`, and volumes created earlier stay in the old directory.

With `reclaimPolicy: Delete`, a helper Pod removes the directory of a deleted volume and mounts the old path to do so.
When Pod Security Standards in system namespaces reject Pods instead of only warning about them,
the SecurityPolicyException `helper-<NAME>` allows only the current path, where `<NAME>` is the name of the LocalPathProvisioner.
The helper Pod for an old volume is then rejected: the PersistentVolume stays `Released`, and its data stays on the node.
With `reclaimPolicy: Retain`, the default, no helper Pod is started on deletion, and this limitation does not apply.

Before changing `path`, delete the PersistentVolumeClaims of the volumes created in the old path.
If such a volume stays `Released` after the change, delete its directory on the node and then the PersistentVolume.

## What happens to the module on the update to Deckhouse Platform 1.79?

Before Deckhouse Platform 1.79, the module was part of Deckhouse Platform.
On the update to 1.79, a cluster that runs the embedded module switches to this module automatically:
the module is downloaded from the `deckhouse` ModuleSource, and the update waits until the module is available there.

The LocalPathProvisioners, the StorageClasses and the PersistentVolumes are kept.
The module Pods restart once, because the image and the image pull Secret change.
Until the new Pods are `Ready`, volumes are not created or removed.
