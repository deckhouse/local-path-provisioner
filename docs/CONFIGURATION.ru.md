---
title: "Модуль local-path-provisioner: настройки"
force_searchable: true
description: "Как включить модуль local-path-provisioner; настроек у модуля нет."
weight: 30
---

У модуля нет настроек.

## Включение модуля

Модуль включён по умолчанию в наборах модулей `Default` и `Managed`. В остальных наборах включите его через ModuleConfig:

```yaml
apiVersion: deckhouse.io/v1alpha1
kind: ModuleConfig
metadata:
  name: local-path-provisioner
spec:
  enabled: true
  version: 1
```

Дождитесь перехода модуля в состояние `Ready`:

```shell
d8 k get module local-path-provisioner -w
```

Чтобы выключить модуль, укажите `spec.enabled: false` в том же ModuleConfig.

## Переход со встроенного модуля

До Deckhouse 1.79 модуль входил в состав Deckhouse. При обновлении на 1.79 кластер, где работает встроенный модуль, переходит на этот модуль автоматически: модуль скачивается из источника модулей `deckhouse`, и обновление ждёт, пока он там не станет доступен. Ресурсы `LocalPathProvisioner`, StorageClass и PersistentVolume сохраняются. Поды модуля один раз перезапускаются, потому что меняются образ и Secret для загрузки образов; пока новые поды не станут `Ready`, тома не создаются и не удаляются.
