---
title: "Модуль local-path-provisioner: FAQ"
description: "Ответы на частые вопросы о модуле local-path-provisioner в Deckhouse Platform: Prometheus на локальных дисках, смена пути томов и обновление на Deckhouse Platform 1.79."
weight: 55
---

## Как хранить данные Prometheus на локальных дисках?

Чтобы хранить данные Prometheus на локальных дисках узлов, где он запущен, выполните следующие шаги:

1. Создайте LocalPathProvisioner.
   В `nodeGroups` укажите группу узлов, где запущены поды Prometheus, в `path` — директорию на узлах для данных:

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

1. В настройках модуля `prometheus` укажите имя нового StorageClass в параметрах [`storageClass`](/modules/prometheus/configuration.html#parameters-storageclass)
   и [`longtermStorageClass`](/modules/prometheus/configuration.html#parameters-longtermstorageclass):

   ```yaml
   storageClass: localpath-system
   longtermStorageClass: localpath-system
   ```

1. Дождитесь перезапуска подов Prometheus.

## Что происходит с существующими томами при изменении пути томов?

Новые тома создаются в новом `path`, а тома, созданные раньше, остаются в старой директории.

При `reclaimPolicy: Delete` директорию удалённого тома удаляет вспомогательный под, и для этого он монтирует старый путь.
Когда Pod Security Standards в системных неймспейсах не только предупреждают о нарушениях, но и отклоняют поды,
SecurityPolicyException `helper-<NAME>` разрешает только текущий путь, где `<NAME>` — имя объекта LocalPathProvisioner.
Тогда вспомогательный под для старого тома отклоняется: PersistentVolume остаётся в статусе `Released`, а его данные — на узле.
При `reclaimPolicy: Retain`, который действует по умолчанию, вспомогательный под при удалении не запускается, и это ограничение не действует.

Перед изменением `path` удалите PersistentVolumeClaim томов, созданных в старом пути.
Если после изменения такой том остался в статусе `Released`, удалите его директорию на узле, а затем сам PersistentVolume.

## Что происходит с модулем при обновлении на Deckhouse Platform 1.79?

До Deckhouse Platform 1.79 модуль входил в состав Deckhouse Platform.
При обновлении на 1.79 кластер, где работает встроенный модуль, переходит на этот модуль автоматически:
модуль скачивается из ModuleSource `deckhouse`, и обновление ждёт, пока модуль не станет там доступен.

Объекты LocalPathProvisioner, StorageClass и PersistentVolume сохраняются.
Поды модуля один раз перезапускаются, потому что меняются образ и Secret для загрузки образов.
Пока новые поды не станут `Ready`, тома не создаются и не удаляются.
