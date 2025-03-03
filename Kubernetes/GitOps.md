# GitOps: Что это и как использовать с GitHub Actions?

## Введение
**GitOps** — это как "умный повар" 🍲, который использует Git (твой "рецепт") для управления приложениями и серверами. А **GitHub Actions** — это как "кухонный робот" 🤖, который помогает всё это автоматизировать. Давай узнаем, что это такое и как их подружить! 🚀

---

## Что такое GitOps?

**GitOps** — это способ управлять приложениями и инфраструктурой (серверами, кластерами), используя Git как "главную книгу рецептов". Всё, что нужно твоему приложению — код, настройки, количество серверов, — хранится в Git-репозитории. Если что-то меняется, ты обновляешь Git, и система сама "готовит" всё по новому рецепту.

### Как это работает?
1. Ты пишешь "рецепт" (например, в YAML-файлах) — что должно быть на сервере.
2. Этот рецепт лежит в Git-репозитории.
3. Специальный "повар" (например, ArgoCD или Flux) смотрит на Git и делает так, чтобы реальный сервер соответствовал рецепту.
4. Если что-то ломается или меняется вручную, "повар" возвращает всё к тому, что в Git.

### Пример из жизни
У тебя есть сайт магазина игрушек. Ты хочешь, чтобы он работал на трёх серверах. В Git ты пишешь: "Запустите мой сайт в трёх копиях". "Повар" (GitOps-инструмент) видит это и запускает три контейнера. Если один сервер "упал", он автоматически добавляет новый, чтобы было ровно три — как в рецепте.

---

## Зачем нужен GitOps?

1. **Всё в одном месте**: код и настройки живут в Git, как единая "книга правды".
2. **Легко чинить**: если что-то сломалось, просто смотришь в Git и возвращаешь всё обратно.
3. **Автоматика**: никаких ручных команд на сервере — всё само обновляется.
4. **История**: Git хранит все изменения, как дневник.

---

## Как использовать GitOps с GitHub Actions?

**GitHub Actions** — это инструмент для автоматизации задач прямо в GitHub. Он идеально подходит для **CI** (Continuous Integration — сборка и проверка кода), а с GitOps можно добавить и **CD** (Continuous Deployment — развёртывание). Вот как это сделать шаг за шагом.

### Простая схема
1. Пишешь код приложения (например, сайт).
2. GitHub Actions собирает его в контейнер (Docker-образ).
3. Actions обновляет "рецепт" в Git-репозитории (например, меняет версию образа).
4. GitOps-инструмент (например, ArgoCD) видит изменения в Git и обновляет сервер.

---

### Пример: Сайт магазина игрушек

Давай настроим GitOps с GitHub Actions для сайта. У нас будет два репозитория:
- **toy-shop-app**: код сайта.
- **toy-shop-infra**: "рецепты" для Kubernetes (инфраструктура).

#### Шаг 1: Код сайта и Dockerfile
В `toy-shop-app` у тебя есть простой сайт и `Dockerfile`:
```dockerfile
FROM node:18
WORKDIR /app
COPY . .
RUN npm install
CMD ["npm", "start"]
EXPOSE 3000
```

#### Шаг 2: Настрой GitHub Actions
Создай файл `.github/workflows/build.yml` в `toy-shop-app`:
```yaml
name: Build and Deploy
on: 
  push:
    branches: [ main ]  # Срабатывает при пуше в main

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4  # Берём код
      - name: Build Docker image
        run: docker build -t ghcr.io/мой-логин/toy-shop:${{ github.sha }} .
      - name: Login to GitHub Container Registry
        run: echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin
      - name: Push Docker image
        run: docker push ghcr.io/мой-логин/toy-shop:${{ github.sha }}
      - name: Update infra repo
        uses: actions/checkout@v4
        with:
          repository: мой-логин/toy-shop-infra  # Второй репозиторий
          token: ${{ secrets.INFRA_TOKEN }}  # Токен для доступа
          path: infra
      - name: Update image tag
        run: |
          sed -i "s|image: ghcr.io/мой-логин/toy-shop:.*|image: ghcr.io/мой-логин/toy-shop:${{ github.sha }}|" infra/deployment.yaml
          cd infra
          git config user.name "GitHub Actions Bot"
          git config user.email "bot@github.com"
          git add .
          git commit -m "Update to ${{ github.sha }}"
          git push
```

- **Что тут происходит?**
  - При пуше в `main` Actions собирает Docker-образ.
  - Образ пушится в GitHub Container Registry.
  - Потом Actions обновляет файл `deployment.yaml` во втором репозитории (`toy-shop-infra`) с новой версией образа.

*(Примечание: Создай токен в GitHub -> Settings -> Developer settings -> Personal access tokens, дай ему доступ к `repo`, и добавь в секреты как `INFRA_TOKEN`.)*

#### Шаг 3: "Рецепт" в `toy-shop-infra`
В `toy-shop-infra` лежит файл `deployment.yaml` для Kubernetes:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: toy-shop
spec:
  replicas: 3  # Три копии сайта
  selector:
    matchLabels:
      app: toy-shop
  template:
    metadata:
      labels:
        app: toy-shop
    spec:
      containers:
      - name: toy-shop
        image: ghcr.io/мой-логин/toy-shop:sha  # Тут меняется версия
        ports:
        - containerPort: 3000
```

#### Шаг 4: Подключаем GitOps-инструмент
Установи **ArgoCD** в свой Kubernetes-кластер:
1. Добавь репозиторий `toy-shop-infra` в ArgoCD:
   ```
   argocd repo add https://github.com/мой-логин/toy-shop-infra --username мой-логин --password твой-токен
   ```
2. Создай приложение:
   ```
   argocd app create toy-shop --repo https://github.com/мой-логин/toy-shop-infra --path . --dest-server https://kubernetes.default.svc --dest-namespace default --sync-policy automated
   ```
3. ArgoCD будет следить за `toy-shop-infra` и обновлять кластер, когда там меняется `deployment.yaml`.

---

## Как это всё работает вместе?

1. Ты пушаешь изменения в `toy-shop-app`.
2. GitHub Actions собирает новый Docker-образ и пушит его.
3. Actions обновляет `deployment.yaml` в `toy-shop-infra`.
4. ArgoCD видит изменения в Git и обновляет твой Kubernetes-кластер.

В итоге: новая версия сайта автоматически появляется у пользователей за пару минут! 🎉

---

## Почему это круто?

- **Автоматика**: ты только пишешь код, остальное само.
- **Надёжность**: всё в Git, легко откатить изменения.
- **Прозрачность**: вся команда видит, что происходит.

---

## Итог

**GitOps** — это как "кулинарная книга" для серверов, где Git — источник правды. **GitHub Actions** — это "робот", который готовит контейнеры и обновляет рецепты. Вместе они делают твою жизнь проще:
- Пишешь код → Actions собирает → GitOps разворачивает.
