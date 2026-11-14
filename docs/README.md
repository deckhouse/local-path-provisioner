---
title: "The local-path-provisioner module"
description: "Management of local storage on Kubernetes nodes in Deckhouse Kubernetes Platform."
weight: 10
---

Local Path Provisioner provides a way for the Kubernetes users to utilize the local storage on each node.

## How does it work?

For each custom resource [LocalPathProvisioner](cr.html#localpathprovisioner), a corresponding `StorageClass` is created.

The allowed topology for `StorageClass` is calculated based on the list of `nodeGroup` names from the custom resource.
The topology is used for scheduling Pods.

When a Pod orders a disk:
- a `HostPath` PV is created
- the provisioner creates a local directory on the desired node along the path consisting of the `path` custom resource parameter, the PV name and the PVC name.

  Example of a path:

  ```shell
  /opt/local-path-provisioner/pvc-d9bd3878-f710-417b-a4b3-38811aa8aac1_d8-monitoring_prometheus-main-db-prometheus-main-0
  ```

## Requirements

- Deckhouse Kubernetes Platform 1.78 or newer.
- The module is available in the CE, BE, SE, SE+ and EE editions of Deckhouse Kubernetes Platform.

## Limitations

- The disk size limit is not supported for the local path provisioned volumes.
