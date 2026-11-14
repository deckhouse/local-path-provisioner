---
title: "Модуль local-path-provisioner"
description: "Управление локальным хранилищем на узлах Kubernetes в Deckhouse Kubernetes Platform."
weight: 10
---

Позволяет пользователям Kubernetes использовать локальное хранилище на узлах.

## Как это работает?

Для каждого кастомного ресурса [LocalPathProvisioner](cr.html#localpathprovisioner) создаётся соответствующий `StorageClass`.

Допустимая топология для `StorageClass` вычисляется на основе списка `nodeGroup` из кастомного ресурса. Топология используется при планировании подов.

Когда под заказывает диск, то:
- создаётся `HostPath` PV;
- провижинер создаёт на нужном узле локальный каталог по пути, состоящем из параметра `path` кастомного ресурса, имени PV и имени PVC.

  Пример пути:

  ```shell
  /opt/local-path-provisioner/pvc-d9bd3878-f710-417b-a4b3-38811aa8aac1_d8-monitoring_prometheus-main-db-prometheus-main-0
  ```

## Требования

- Deckhouse Kubernetes Platform 1.78 или новее.
- Модуль доступен в редакциях CE, BE, SE, SE+ и EE Deckhouse Kubernetes Platform.

## Ограничения

- Ограничение на размер диска не поддерживается для локальных томов.
