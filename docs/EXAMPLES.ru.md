---
title: "Модуль local-path-provisioner: примеры"
description: "Примеры ресурсов LocalPathProvisioner, которые создают StorageClass для локальных томов на узлах Deckhouse."
weight: 52
---

## Пример custom resource `LocalPathProvisioner`

Reclaim policy устанавливается по умолчанию в `Retain`.

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

## Пример custom resource `LocalPathProvisioner` с установленным `reclaimPolicy`

Reclaim policy устанавливается в `Delete`.

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
