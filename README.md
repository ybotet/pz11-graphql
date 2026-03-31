# Практическая работа №8: CI/CD Pipeline с GitHub Actions

## Содержание

1. [Описание работы](#1-описание-работы)
2. [Установка и настройка](#2-установка-и-настройка)
3. [Структура проекта](#3-структура-проекта)
4. [Pipeline CI/CD](#4-pipeline-cicd)
   - 4.1. Файл пайплайна
   - 4.2. Описание Jobs
5. [Публикация в реестре](#5-публикация-в-реестре)
6. [Стратегия версионирования](#6-стратегия-версионирования)
7. [Скриншоты выполнения](#7-скриншоты-выполнения)
8. [Контрольные вопросы](#8-контрольные-вопросы)

---

## 1. Описание работы

Данная практическая работа заключается в настройке **CI/CD пайплайна** с использованием **GitHub Actions**. Цель — автоматизировать процесс проверки, сборки и публикации Docker-образов микросервисов `auth` и `tasks`.

### Достигнутые цели

| Цель                                                  |
| ----------------------------------------------------- |
| Настройка пайплайна с проверкой компиляции            |
| Сборка Docker-образов в CI                            |
| Публикация образов в GitHub Container Registry (GHCR) |
| Реализация автоматического версионирования            |
| Использование секретов CI для аутентификации          |

### Используемые технологии

- **CI/CD:** GitHub Actions
- **Язык:** Go 1.25
- **Контейнеризация:** Docker
- **Реестр:** GitHub Container Registry (ghcr.io)

---

## 2. Установка и настройка

### Требования

- Репозиторий на GitHub
- Docker (для локального тестирования)
- Go 1.25 (для разработки)

### Настройка пайплайна

Пайплайн активируется автоматически при каждом **push** или **pull request** в ветки `main`, `master` и `develop`.

Для публикации в GHCR не требуются дополнительные секреты, так как GitHub Actions автоматически предоставляет токен `GITHUB_TOKEN` с правами на запись.

### Разрешения

```yaml
permissions:
  contents: read
  packages: write
```

## 3. Структура проекта

![alt text](<public/Снимок экрана 2026-03-31 152009.png>)

## 4. Pipeline CI/CD

### 4.1. Файл пайплайна

Путь: .github/workflows/ci.yml

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, master, develop]
  pull_request:
    branches: [main, master]

env:
  GO_VERSION: "1.25"
  REGISTRY: ghcr.io
  IMAGE_NAME_AUTH: ${{ github.repository }}/auth
  IMAGE_NAME_TASKS: ${{ github.repository }}/tasks

jobs:
  # ========== JOB 1: VERIFICACIÓN DE COMPILACIÓN (CI) ==========
  verify-build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Download dependencies (root)
        run: go mod download

      - name: Download dependencies (gen)
        run: go mod download
        working-directory: ./gen

      - name: Verify auth service compiles
        run: go build -o /dev/null ./services/auth/cmd/auth/...

      - name: Verify tasks service compiles
        run: go build -o /dev/null ./services/tasks/cmd/tasks/...

      - name: Verify shared package compiles
        run: go build -o /dev/null ./shared/...

  # ========== JOB 2: DOCKER BUILD Y PUBLISH (CD) ==========
  docker-build-push:
    needs: verify-build
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # ===== Build and push AUTH service =====
      - name: Extract metadata for auth
        id: meta_auth
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_AUTH }}
          tags: |
            type=sha,prefix=,format=short
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master' }}
            type=ref,event=branch

      - name: Build and push auth image
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./services/auth/Dockerfile
          push: true
          tags: ${{ steps.meta_auth.outputs.tags }}
          labels: ${{ steps.meta_auth.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      # ===== Build and push TASKS service =====
      - name: Extract metadata for tasks
        id: meta_tasks
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_TASKS }}
          tags: |
            type=sha,prefix=,format=short
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master' }}
            type=ref,event=branch

      - name: Build and push tasks image
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./services/tasks/Dockerfile
          push: true
          tags: ${{ steps.meta_tasks.outputs.tags }}
          labels: ${{ steps.meta_tasks.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Show published images info
        run: |
          echo " Auth image published: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_AUTH }}"
          echo "   Tags: ${{ steps.meta_auth.outputs.tags }}"
          echo ""
          echo " Tasks image published: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_TASKS }}"
          echo "   Tags: ${{ steps.meta_tasks.outputs.tags }}"
```

### 4.2. Описание Jobs

Job 1: verify-build (Непрерывная интеграция)

| Шаг                    | Описание                                                |
| ---------------------- | ------------------------------------------------------- |
| checkout               | Получение кода из репозитория                           |
| setup-go               | Установка Go версии 1.25 с кэшированием                 |
| go mod download (root) | Загрузка зависимостей главного модуля                   |
| go mod download (gen)  | Загрузка зависимостей сгенерированного модуля           |
| verify auth            | Компиляция сервиса auth для проверки отсутствия ошибок  |
| verify tasks           | Компиляция сервиса tasks для проверки отсутствия ошибок |
| verify shared          | Компиляция общего пакета                                |

Job 2: docker-build-push (Непрерывная доставка)

| Шаг                | Описание                                                                 |
| ------------------ | ------------------------------------------------------------------------ |
| checkout           | Получение кода из репозитория                                            |
| setup-buildx       | Настройка Docker Buildx для оптимизированной сборки                      |
| login to GHCR      | Аутентификация в GitHub Container Registry с использованием GITHUB_TOKEN |
| metadata auth      | Генерация тегов для образа auth (latest, sha, branch)                    |
| build & push auth  | Сборка и публикация образа auth в GHCR                                   |
| metadata tasks     | Генерация тегов для образа tasks                                         |
| build & push tasks | Сборка и публикация образа tasks в GHCR                                  |

## 5. Публикация в реестре

Образы публикуются в **GitHub Container Registry (GHCR)**.

### Опубликованные образы

| Сервис | Образ                                      |
| ------ | ------------------------------------------ |
| Auth   | `ghcr.io/ybotet/pz8-pipelinecicd-go/auth`  |
| Tasks  | `ghcr.io/ybotet/pz8-pipelinecicd-go/tasks` |

### Загрузка образов

```bash
# Auth
docker pull ghcr.io/ybotet/pz8-pipelinecicd-go/auth:latest
docker pull ghcr.io/ybotet/pz8-pipelinecicd-go/auth:29fc91a

# Tasks
docker pull ghcr.io/ybotet/pz8-pipelinecicd-go/tasks:latest
docker pull ghcr.io/ybotet/pz8-pipelinecicd-go/tasks:29fc91a
```

## 6. Стратегия версионирования

Каждый образ получает несколько тегов, автоматически генерируемых **docker/metadata-action**:

| Тег                      | Условие              | Назначение                                                                                          |
| ------------------------ | -------------------- | --------------------------------------------------------------------------------------------------- |
| `<short-hash>` (29fc91a) | Всегда               | Точная версия коммита. Обеспечивает прослеживаемость и возможность развертывания конкретной версии. |
| `latest`                 | Только в main/master | Последняя стабильная версия для сред разработки и тестирования.                                     |
| `<branch-name>` (main)   | Всегда               | Например: main или develop. Позволяет идентифицировать ветку происхождения.                         |

### Преимущества

- **Прослеживаемость**: Каждый образ привязан к конкретному коммиту.
- **Воспроизводимость**: Возможность развернуть конкретную версию по хешу.
- **Простота**: Тег `latest` упрощает разработку.

## 7. Скриншоты выполнения

Успешное выполнение пайплайна в **GitHub Actions**.

![alt text](<public/Снимок экрана 2026-03-31 145510.png>)

Детали job verify-build

![alt text](<public/Снимок экрана 2026-03-31 151055.png>)

Детали job docker-build-push

![alt text](<public/Снимок экрана 2026-03-31 151111.png>)

Опубликованные пакеты в GHCR

![alt text](<public/Снимок экрана 2026-03-31 145803.png>)

## 8. Контрольные вопросы

### 1. Чем CI отличается от CD?

CI (Непрерывная интеграция) — это практика частого объединения изменений кода в общем репозитории. Каждое объединение автоматически проверяется с помощью тестов и анализа. В данном пайплайне job `verify-build` представляет часть CI, так как проверяет, что код компилируется без ошибок.

CD (Непрерывная доставка / развертывание) — это расширение CI, автоматизирующее доставку программного обеспечения. Различают:

- Непрерывная доставка: код готов к развертыванию, но требуется ручное утверждение.
- Непрерывное развертывание: каждое изменение, прошедшее проверки, автоматически развертывается в production.

В данном пайплайне job `docker-build-push` представляет часть CD, так как автоматически собирает и публикует Docker-образы в реестре.

Ключевое отличие: CI фокусируется на проверке кода (тесты, компиляция), а CD — на доставке артефакта (образа, бинарного файла) в реестр или среду выполнения.

### 2. Почему `go test` должен запускаться в пайплайне?

Хотя в данном проекте не реализованы unit-тесты, в профессиональной среде их выполнение критически важно по следующим причинам:

| Причина                     | Объяснение                                                             |
| --------------------------- | ---------------------------------------------------------------------- |
| Раннее обнаружение ошибок   | Баги выявляются в момент коммита, а не спустя дни в production.        |
| Предотвращение регрессий    | Гарантирует, что новый код не нарушает существующую функциональность.  |
| Живая документация          | Тесты служат исполняемой документацией ожидаемого поведения.           |
| Уверенность в развертывании | Позволяет выполнять автоматические развертывания с гарантией качества. |

В данном пайплайне используется `go build` как минимальная проверка компиляции, но в идеале необходимо включать `go test ./...` с проверкой покрытия.

### 3. Что такое секреты CI и почему их нельзя хранить в репозитории?

Секреты CI — это конфиденциальные переменные, которые хранятся в зашифрованном виде в платформе CI/CD (GitHub Secrets, GitLab CI/CD Variables). Примеры: пароли, токены API, SSH-ключи, JWT-ключи.

Их нельзя хранить в репозитории, потому что:

| Риск                       | Последствие                                                                   |
| -------------------------- | ----------------------------------------------------------------------------- |
| Раскрытие в Git            | Секреты остаются видимыми в истории коммитов навсегда.                        |
| Несанкционированный доступ | Любой с доступом к репозиторию может увидеть секреты.                         |
| Docker-образы              | Если секреты скопированы в образ, их можно извлечь командой `docker history`. |
| Случайная утечка           | Ошибочный коммит может сделать секреты публичными.                            |

Правильный подход: использовать встроенную систему секретов платформы CI. В GitHub Actions доступ осуществляется через `${{ secrets.ИМЯ_СЕКРЕТА }}`.

### 4. Почему важно версионировать Docker-образы?

Версионирование Docker-образов критически важно для:

| Причина           | Объяснение                                                                               |
| ----------------- | ---------------------------------------------------------------------------------------- |
| Прослеживаемость  | Возможность точно знать, какая версия кода выполняется в каждой среде.                   |
| Откат (Rollback)  | Возможность вернуться к предыдущей версии при неудачном развертывании.                   |
| Воспроизводимость | Гарантия, что один и тот же код ведет себя одинаково в разработке, staging и production. |
| Аудит             | Возможность определить, какие изменения привели к ошибке.                                |

Используемая стратегия: хеш коммита (`sha-xxxxx`) + `latest`. Это позволяет как развертывать конкретные версии (по хешу), так и быстро тестировать (`latest`).

### 5. Какие риски существуют при автоматическом развертывании без ручного контроля?

Автоматическое развертывание (непрерывное развертывание) имеет риски, которые необходимо снижать:

| Риск                   | Описание                                                | Способы снижения                                              |
| ---------------------- | ------------------------------------------------------- | ------------------------------------------------------------- |
| Сбой в production      | Баг, прошедший тесты, может нарушить работу сервиса.    | Канареечные развертывания (canary), blue-green deployment.    |
| Зависимость от CI      | При сбое пайплайна развертывание не выполняется.        | Мониторинг и алерты состояния пайплайна.                      |
| Отсутствие утверждения | Критические изменения могут развернуться без проверки.  | Ручные утверждения для защищенных веток.                      |
| Сложность отката       | Без версионирования откат затруднителен.                | Семантическое версионирование, хранение истории образов.      |
| Безопасность           | Скомпрометированный секрет может дать доступ к серверу. | Регулярная ротация секретов, политика минимальных привилегий. |

В данном пайплайне автоматическое развертывание на сервер не реализовано, но при его добавлении рекомендуется:

- Сначала развертывать в staging-среду.
- Использовать ручные утверждения для production.
- Сохранять версионирование для быстрого отката.

## 9. Заключение

Успешно реализован пайплайн CI/CD с использованием GitHub Actions, который:

- Проверяет компиляцию сервисов `auth` и `tasks` при каждом push.
- Собирает оптимизированные Docker-образы с использованием multi-stage сборки.
- Публикует образы в GitHub Container Registry.
- Реализует стратегию версионирования с тегами `latest`, `<short-hash>` и `<branch-name>`.

Пайплайн является воспроизводимым, безопасным (использует секреты GitHub) и полностью автоматизированным, что соответствует целям практической работы.
#   p z 9 - r e d i s  
 