### Инфраструктура сбора телеметрии (Метрики и Логи)

Данный репозиторий содержит конфигурацию для развертывания локального или продакшен-ready стека мониторинга на базе **OpenTelemetry**, **VictoriaMetrics** и **VictoriaLogs**. Стек предназначен для нативного приема данных по унифицированному протоколу **OTLP (OpenTelemetry Protocol)**. 

### Архитектура системы

1. **Приложения (Applications)** — отправляют метрики и логи по протоколу OTLP (gRPC на порт 4317 или HTTP на порт 4318).
2. **OpenTelemetry Collector** — принимает OTLP-трафик, обрабатывает его пакетами (batching) и маршрутизирует в соответствующие хранилища по HTTP.
3. **VictoriaMetrics** — высокопроизводительная TSDB, нативно принимающая OTLP-метрики на эндпоинт /opentelemetry/v1/metrics.
4. **VictoriaLogs** — масштабируемая база данных для логов, нативно принимающая OTLP-логи на эндпоинт /insert/opentelemetry/v1/logs.
5. **Grafana** — единая панель визуализации для построения графиков и анализа логов.

### Структура файлов проекта

* docker-compose.yaml — описание и сетевая связка сервисов, конфигурация томов данных и Grafana.
* otel-collector-config.yaml — настройки конвейеров (pipelines) обработки данных в OTel Collector.
* Makefile — утилита автоматизации рутинных команд (запуск, остановка, логи).

### Быстрый старт

### 1. Предварительные требования

* Установленный **Docker** и **Docker Compose v2+**.
* Утилита **make** (опционально, для использования Makefile).

### 2. Запуск стека

Выполните команду для сборки и инициализации всех контейнеров в фоновом режиме: 

bash

make up

Use code with caution.

*Если утилита make недоступна, используйте:* docker compose up -d 

### 3. Доступные порты и веб-интерфейсы

После успешного запуска сервисы доступны по следующим адресам: 

СервисПорт (Хост)Описание / Веб-интерфейс
****OTel Collector (gRPC)****
4317Точка приема OTLP gRPC от приложений
****OTel Collector (HTTP)****
4318Точка приема OTLP HTTP от приложений
****VictoriaMetrics****
8428Встроенный графический интерфейс: http://localhost:8428/vmui/
****VictoriaLogs****
9428Встроенный веб-интерфейс логов: http://localhost:9428/vlogs/
****Grafana****
3000Панель мониторинга: http://localhost:3000 (Логин/Пароль: admin / admin)

### Проверка работоспособности (Smoke Tests)

Вы можете симулировать отправку данных от приложения с помощью стандартных HTTP-запросов через curl. 

### Отправка тестового лога

bash

curl -i -X POST http://localhost:4318/v1/logs \
-H "Content-Type: application/json" \
-d '{
"resourceLogs": [{
"resource": {
"attributes": [{"key": "service.name", "value": {"stringValue": "doc-test-service"}}]
},
"scopeLogs": [{
"logRecords": [{
"timeUnixNano": "'$(date +%s)00000000'",
"body": {"stringValue": "Проверка отправки логов через OTel в VictoriaLogs"},
"attributes": [{"key": "log.severity", "value": {"stringValue": "INFO"}}]
}]
}]
}]
}'

Use code with caution.

*Убедитесь в получении ответа HTTP/1.1 200 OK. Текст лога должен отобразиться в веб-интерфейсе http://localhost:9428/vlogs/ при поиске по слову doc-test-service.* 

### Отправка тестовой метрики

bash

curl -i -X POST http://localhost:4318/v1/metrics \
-H "Content-Type: application/json" \
-d '{
"resourceMetrics": [{
"resource": {
"attributes": [{"key": "service.name", "value": {"stringValue": "doc-test-service"}}]
},
"scopeMetrics": [{
"metrics": [{
"name": "doc_success_clicks",
"description": "Тестовый счетчик кликов",
"sum": {
"dataPoints": [{
"asInt": "100",
"timeUnixNano": "'$(date +%s)00000000'"
}],
"aggregationTemporality": 1,
"isMonotonic": true
}
}]
}]
}]
}'

Use code with caution.

*Метрика doc_success_clicks со значением 100 станет доступна на графиках в http://localhost:8428/vmui/.* 

### Настройка Grafana

При первом входе используйте стандартные данные admin / admin (система попросит изменить пароль). 

### 1. Подключение метрик (VictoriaMetrics)

1. Перейдите в **Connections** -> **Data Sources** -> **Add data source**.
2. Выберите тип **Prometheus**.
3. В поле **Connection URL** введите внутренний сетевой адрес Docker: http://victoria-metrics:8428.
4. Прокрутите вниз и нажмите **Save & test**.

### 2. Подключение логов (VictoriaLogs)

1. В контейнер автоматически предустановлен официальный плагин victoriametrics-logs-datasource.
2. Перейдите в **Connections** -> **Data Sources** -> **Add data source**.
3. Найдите и выберите **VictoriaLogs**.
4. В поле **URL** укажите: http://victoria-logs:9428.
5. Нажмите **Save & test**.

### Справочник команд (Makefile)

В корневой директории доступна автоматизация основных команд Docker Compose: 

* make up — запустить весь стек в бэкграунде.
* make down — остановить работу контейнеров и удалить их (данные в volumes сохраняются).
* make restart — быстрая перезагрузка всех компонентов.
* make logs — просмотр логов только для OpenTelemetry Collector (удобно для отладки маршрутизации).
* make logs-all — вывод логов всех контейнеров в один поток.
* make ps — отображение статуса контейнеров и маппинга портов.
* make clean — **Внимание:** полная остановка стека и удаление всех Docker Volumes (базы данных будут полностью очищены).