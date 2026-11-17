---
title: "Модуль local-path-provisioner: примеры"
description: "Примеры манифестов LocalPathProvisioner для модуля local-path-provisioner в Deckhouse Platform, которые сохраняют или удаляют данные томов."
weight: 52
---

## StorageClass, который сохраняет данные томов

Следующий LocalPathProvisioner создаёт StorageClass `localpath-system` для узлов группы `system`.
Тома создаются в `/opt/local-path-provisioner` на этих узлах.
После удаления PersistentVolumeClaim его PersistentVolume и данные в директории остаются, потому что `reclaimPolicy` по умолчанию равен `Retain`:

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

## StorageClass, который удаляет данные томов

Следующий LocalPathProvisioner задаёт `reclaimPolicy` равным `Delete`.
После удаления PersistentVolumeClaim модуль удаляет его PersistentVolume и директорию с данными:

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
