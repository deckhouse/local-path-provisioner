---
title: "Модуль local-path-provisioner: примеры"
description: "Примеры ресурсов LocalPathProvisioner, которые создают StorageClass для локальных томов на узлах Deckhouse."
weight: 52
---

## Пример кастомного ресурса `LocalPathProvisioner`

Политика освобождения тома (`reclaimPolicy`) по умолчанию — `Retain`.

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

## Пример кастомного ресурса `LocalPathProvisioner` с установленным `reclaimPolicy`

Политика освобождения тома — `Delete`.

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
