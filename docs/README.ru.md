---
title: "Модуль local-path-provisioner"
description: "Модуль local-path-provisioner в Deckhouse Platform предоставляет тома в директориях на локальных дисках узлов кластера."
weight: 10
---

Модуль `local-path-provisioner` предоставляет тома в директориях на локальных дисках узлов кластера.
Под, который использует такой том, запускается на узле, где хранятся его данные.

Каждый объект LocalPathProvisioner описывает StorageClass: группы узлов, которые он обслуживает, и директорию на их узлах, в которой создаются тома.

## Основные возможности

Модуль даёт следующие возможности:

- StorageClass для локальных томов, который создаётся из LocalPathProvisioner.
- Размещение томов только на узлах выбранных групп узлов.
- Выбор, сохраняются ли данные тома после удаления его PersistentVolumeClaim.

## Как работает модуль

Модуль создаёт и удаляет тома так:

1. Для каждого объекта LocalPathProvisioner модуль создаёт StorageClass с тем же именем.
   Если задан `nodeGroups`, StorageClass допускает только узлы этих групп.
1. PersistentVolumeClaim с этим StorageClass получает том после того, как использующий его под назначен на узел.
1. Провижинер создаёт директорию на этом узле и PersistentVolume типа `hostPath` для неё.
   Директория создаётся в `path` объекта LocalPathProvisioner и называется `<PV_NAME>_<NAMESPACE>_<PVC_NAME>`:

   ```text
   /opt/local-path-provisioner/pvc-d9bd3878-f710-417b-a4b3-38811aa8aac1_d8-monitoring_prometheus-main-db-prometheus-main-0
   ```

1. После удаления PersistentVolumeClaim судьбу тома определяет `reclaimPolicy` объекта LocalPathProvisioner.
   При `Retain` PersistentVolume и директория с данными остаются.
   При `Delete` модуль удаляет директорию и PersistentVolume.

## Ограничения

Размер тома не ограничивается: том может занять всё свободное место на своём диске, какой бы размер ни запросил PersistentVolumeClaim.
Данные тома хранятся только на одном узле и теряются вместе с ним.
