# Установка и использование Log Linter

## Способ 1: Standalone использование

### Установка

```bash
go install github.com/dontpanicw/log-linter/cmd/loglinter@latest
```

### Использование

```bash
# Проверить текущую директорию
loglinter ./...

# Проверить конкретный пакет
loglinter ./internal/...

# Проверить конкретный файл
loglinter main.go
```

## Способ 2: Интеграция с golangci-lint

### Требования

- Go 1.22+
- golangci-lint v1.50+

### Шаг 1: Сборка плагина

```bash
# Клонировать репозиторий
git clone https://github.com/dontpanicw/log-linter.git
cd log-linter

# Собрать плагин
make plugin
```

Это создаст файл `loglinter.so` в корне проекта.

### Шаг 2: Настройка golangci-lint

Создайте или обновите файл `.golangci.yml` в корне вашего проекта:

```yaml
linters-settings:
  custom:
    loglinter:
      path: /path/to/loglinter.so
      description: Checks log messages for compliance with logging standards
      original-url: github.com/dontpanicw/log-linter

linters:
  enable:
    - loglinter
```

### Шаг 3: Запуск

```bash
golangci-lint run
```

## Способ 3: Использование в CI/CD

### GitHub Actions

```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Install loglinter
        run: go install github.com/dontpanicw/log-linter/cmd/loglinter@latest
      
      - name: Run loglinter
        run: loglinter ./...
```

### GitLab CI

```yaml
lint:
  image: golang:1.22
  stage: test
  script:
    - go install github.com/dontpanicw/log-linter/cmd/loglinter@latest
    - loglinter ./...
```

## Правила проверки

Линтер проверяет следующие правила:

### 1. Сообщения должны начинаться со строчной буквы

❌ Неправильно:
```go
log.Info("Starting server")
slog.Error("Failed to connect")
```

✅ Правильно:
```go
log.Info("starting server")
slog.Error("failed to connect")
```

### 2. Сообщения должны быть на английском языке

❌ Неправильно:
```go
log.Info("запуск сервера")
```

✅ Правильно:
```go
log.Info("starting server")
```

### 3. Без спецсимволов и эмодзи

❌ Неправильно:
```go
log.Info("server started!🚀")
log.Error("connection failed!!!")
```

✅ Правильно:
```go
log.Info("server started")
log.Error("connection failed")
```

### 4. Без чувствительных данных

❌ Неправильно:
```go
log.Info("user password: " + password)
log.Debug("api_key=" + apiKey)
```

✅ Правильно:
```go
log.Info("user authenticated successfully")
log.Debug("api request completed")
```

## Поддерживаемые логгеры

- `log/slog` (стандартная библиотека Go)
- `go.uber.org/zap`

## Troubleshooting

### Плагин не загружается в golangci-lint

Убедитесь, что:
1. Версия Go совпадает при сборке плагина и golangci-lint
2. Путь к `.so` файлу указан правильно
3. У файла есть права на выполнение

### Ложные срабатывания

Если линтер выдает ложные срабатывания, вы можете:
1. Использовать `//nolint:loglinter` для игнорирования конкретной строки
2. Настроить исключения в `.golangci.yml`

```yaml
issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - loglinter
```
