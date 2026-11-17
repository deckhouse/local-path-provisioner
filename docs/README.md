---
title: "Local-path-provisioner module"
description: "The local-path-provisioner module of Deckhouse Platform provides volumes in directories on the local disks of the cluster nodes."
weight: 10
---

The `local-path-provisioner` module provides volumes in directories on the local disks of the cluster nodes.
A Pod that uses such a volume runs on the node that holds its data.

Each LocalPathProvisioner describes a StorageClass: the node groups it serves and the directory on their nodes where volumes are created.

## Main features

The module offers the following:

- A StorageClass for local volumes, created from a LocalPathProvisioner.
- Volumes placed only on the nodes of the chosen node groups.
- A choice of whether the volume data is kept or removed after its PersistentVolumeClaim is deleted.

## How the module works

The module creates and removes volumes as follows:

1. For each LocalPathProvisioner, the module creates a StorageClass with the same name.
   If `nodeGroups` is set, the StorageClass allows only the nodes of these node groups.
1. A PersistentVolumeClaim of this StorageClass gets a volume after the Pod that uses it is scheduled to a node.
1. The provisioner creates a directory on that node and a `hostPath` PersistentVolume for it.
   The directory is created in `path` of the LocalPathProvisioner and is named `<PV_NAME>_<NAMESPACE>_<PVC_NAME>`:

   ```text
   /opt/local-path-provisioner/pvc-d9bd3878-f710-417b-a4b3-38811aa8aac1_d8-monitoring_prometheus-main-db-prometheus-main-0
   ```

1. After the PersistentVolumeClaim is deleted, `reclaimPolicy` of the LocalPathProvisioner decides what happens to the volume.
   With `Retain`, the PersistentVolume and the directory with the data stay.
   With `Delete`, the module removes the directory and the PersistentVolume.

## Limitations

The size of a volume is not limited: a volume can use all free space of the disk it is on, whatever size the PersistentVolumeClaim requests.
The data of a volume is kept on one node only and is lost together with the node.
