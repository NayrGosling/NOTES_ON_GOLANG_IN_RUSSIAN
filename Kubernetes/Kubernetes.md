# Лекция: Kubernetes

## Введение
Kubernetes (K8s) – это мощная система оркестрации контейнеров, позволяющая автоматизировать развертывание, управление и масштабирование контейнерных приложений. На уровне разработчика важно понимать не только базовые компоненты, но и внутреннее устройство, best practices и advanced-техники управления кластером.

---

## Архитектура Kubernetes

### Основные компоненты

- **Control Plane** (Плоскость управления):
  - **API Server** – главный компонент, принимающий все запросы и обеспечивающий взаимодействие с etcd.
  - **etcd** – распределенное хранилище для всех конфигурационных данных кластера.
  - **Controller Manager** – управляет контроллерами (репликации, узлов, сервисов и т.д.).
  - **Scheduler** – отвечает за распределение Pod'ов по узлам.

- **Worker Nodes** (Рабочие узлы):
  - **Kubelet** – агент, запускающий контейнеры и следящий за их состоянием.
  - **Container Runtime** – среда исполнения контейнеров (например, containerd, CRI-O, Docker).
  - **Kube Proxy** – сетевой прокси, обеспечивающий маршрутизацию трафика между сервисами.

![Картинка-1](./images/kubernetes-1.png)
---

## Cети Kubernetes

### CNI (Container Network Interface)
Kubernetes использует CNI-плагины для сетевого взаимодействия между контейнерами. Популярные реализации:
- **Calico** – обеспечивает Network Policy и BGP.
- **Flannel** – простая overlay-сеть.
- **Cilium** – eBPF-решение для высокопроизводительных сетей.

### Service Mesh (Istio, Linkerd)
Используются для продвинутого управления трафиком, балансировки нагрузки, мониторинга и безопасности.

### Network Policies
Позволяют ограничивать взаимодействие между Pod'ами, обеспечивая дополнительную безопасность.

---

## [Управление состоянием приложений](https://github.com/NayrGosling/NOTES_ON_GOLANG_IN_RUSSIAN/blob/main/Kubernetes/Безопасность%20в%20Kubernetes.md)

### StatefulSets
Используется для приложений, требующих сохранения состояния (например, базы данных). В отличие от Deployment, каждый Pod в StatefulSet получает уникальный хостнейм и сохраняет данные при перезапуске.

### Custom Resource Definitions (CRD) и операторы
Позволяют расширять функциональность Kubernetes, создавая новые типы ресурсов и автоматизируя управление сложными сервисами (например, PostgreSQL-оператор).

### Helm и Kustomize
- **Helm** – пакетный менеджер для Kubernetes, позволяющий управлять чартами.
- **Kustomize** – инструмент для декларативного описания Kubernetes-манифестов без шаблонов.

---

## [Масштабирование и отказоустойчивость](https://github.com/NayrGosling/NOTES_ON_GOLANG_IN_RUSSIAN/blob/main/Kubernetes/Масштабирование%20и%20отказоустойчивость%20в%20Kubernetes.md)

### Horizontal Pod Autoscaler (HPA)
Автоматически масштабирует Pods на основе метрик (CPU, RAM, custom metrics).

### Vertical Pod Autoscaler (VPA)
Динамически изменяет ресурсы для Pod'ов.

### Cluster Autoscaler
Добавляет/удаляет узлы в кластере в зависимости от нагрузки.

---

## [Логирование, мониторинг и отладка](https://github.com/NayrGosling/NOTES_ON_GOLANG_IN_RUSSIAN/blob/main/Kubernetes/Логирование%2C%20мониторинг%20и%20отладка%20в%20Kubernetes.md)

### Prometheus + Grafana
Используется для сбора и визуализации метрик Kubernetes.

### Loki + Fluentd
Система централизованного логирования.

### Debugging tools
- **kubectl logs** – просмотр логов контейнеров.
- **kubectl exec** – выполнение команд в контейнере.
- **kubectl port-forward** – проброс портов для отладки.

---

## [Безопасность в Kubernetes](https://github.com/NayrGosling/NOTES_ON_GOLANG_IN_RUSSIAN/blob/main/Kubernetes/Безопасность%20в%20Kubernetes.md)

### RBAC (Role-Based Access Control)
Ограничивает права пользователей и сервисов в кластере.

### Pod Security Standards
Политики безопасности для Pod'ов (restricted, baseline, privileged).

### Secrets и ConfigMaps
Используются для управления конфиденциальными данными.

### Network Policies и PSP (Pod Security Policies)
Контроль сетевых взаимодействий и ограничение привилегий контейнеров.

---

## [CI/CD в Kubernetes](https://github.com/NayrGosling/NOTES_ON_GOLANG_IN_RUSSIAN/blob/main/Kubernetes/CI%20CD%20в%20Kubernetes.md)

### ArgoCD и Flux
GitOps-подход для автоматизированного развертывания.

### Jenkins, Tekton
Инструменты CI/CD для управления жизненным циклом приложений.

---

## Kubernetes в продакшене: best practices

1. **Использование namespaces** для логической изоляции сервисов.
2. **Resource Requests & Limits** для предотвращения «захвата» ресурсов.
3. **Liveness & Readiness probes** для мониторинга здоровья приложений.
4. **PDB (Pod Disruption Budget)** для обеспечения высокой доступности.
5. **Использование distroless-образов** для повышения безопасности контейнеров.
6. **Контроль версий API** – избегайте deprecated API.
7. **Правильная организация логирования и мониторинга**.

---

## Заключение
Kubernetes – это мощный инструмент, который требует глубокого понимания для эффективного использования в продакшене. Разработчику важно владеть не только базовыми концепциями, но и продвинутыми техниками для масштабирования, отказоустойчивости, безопасности и автоматизации управления кластером.
