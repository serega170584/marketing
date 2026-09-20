.PHONY: help up down restart logs logs-all ps clean lint lint-fix fmt fmt-check

# Переменная для вызова docker compose (поддерживает как старый синтаксис, так и новый)
DC = docker compose

help: ## Показать справку по командам
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Запустить все сервисы в фоновом режиме (VictoriaMetrics, VictoriaLogs, OTel, Grafana)
	$(DC) up -d

down: ## Остановить и удалить контейнеры стека
	$(DC) down

restart: ## Перезапустить все сервисы
	$(DC) restart

logs: ## Посмотреть логи OpenTelemetry Collector в реальном времени
	$(DC) logs -f otel-collector

logs-all: ## Посмотреть логи всех сервисов в реальном времени
	$(DC) logs -f

ps: ## Статус запущенных контейнеров и их порты
	$(DC) ps

clean: ## Остановить стек и ПОЛНОСТЬЮ удалить все сохраненные данные (volumes)
	$(DC) down -v

fmt: ## Автоматически отформатировать весь проект (go fmt)
	go fmt ./...

fmt-check: ## Проверить, весь ли код отформатирован (без изменения файлов)
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Следующие файлы не отформатированы:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "Весь код успешно отформатирован!"

lint: ## Запустить проверку кода с помощью golangci-lint
	golangci-lint run ./...

lint-fix: fmt ## Запустить форматирование и исправить доступные ошибки линтера
	golangci-lint run --fix ./...
