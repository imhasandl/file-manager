# File Manager API

## Описание

Этот сервис позволяет скачивать файлы (.pdf, .jpeg) по ссылкам из интернета, упаковывать их в ZIP-архив и возвращать пользователю. Система работает через HTTP API и не требует внешних баз данных.

## Запуск сервера

1. Установите Go (1.18+)
2. Клонируйте репозиторий и перейдите в папку проекта
3. Создайте файл `.env` с содержимым:
   ```
   PORT=8080
   ```
4. Установите зависимости:
   ```bash
   go mod tidy
   ```
5. Запустите сервер:
   ```bash
   go run main.go
   ```

## API Endpoints

### 1. Создать задачу
- **POST /tasks**
- Ответ:
  ```json
  {
    "id": "<uuid>",
    "files": [],
    "created_at": "2025-07-12T12:00:00Z"
  }
  ```
- Ограничение: одновременно не более 3 активных задач

### 2. Добавить файл в задачу
- **POST /tasks/{taskID}/files**
- Тело запроса:
  ```json
  { "url": "https://example.com/file.pdf" }
  ```
- Ответ:
  ```json
  {
    "id": "<uuid>",
    "files": [ { "url": "...", "filename": "..." } ],
    "created_at": "..."
  }
  ```
- Ограничения: только .pdf и .jpeg, максимум 3 файла на задачу

### 3. Получить статус задачи
- **GET /tasks/{taskID}/status**
- Ответ:
  ```json
  {
    "id": "<uuid>",
    "files": [ ... ],
    "created_at": "...",
    "archive_url": "/download/archive_<uuid>.zip"
  }
  ```
- Если архив еще не готов, поле `archive_url` отсутствует

### 4. Скачать архив
- **GET /download/archive_{taskID}.zip**
- Ответ: ZIP-файл

## Примеры использования (curl)

1. Создать задачу:
   ```bash
   curl -X POST http://localhost:8080/tasks
   ```
2. Добавить файл:
   ```bash
   curl -X POST http://localhost:8080/tasks/<uuid>/files \
     -H "Content-Type: application/json" \
     -d '{"url":"https://example.com/file.pdf"}'
   ```
3. Проверить статус:
   ```bash
   curl http://localhost:8080/tasks/<uuid>/status
   ```
4. Скачать архив:
   ```bash
   curl http://localhost:8080/download/archive_<uuid>.zip --output archive.zip
   ```

## Ограничения
- Максимум 3 задачи одновременно
- Максимум 3 файла в задаче
- Только .pdf и .jpeg

## Структура проекта
- main.go — запуск сервера
- handlers/ — обработчики HTTP
- models/ — структуры данных
- utils/ — вспомогательные функции
- archives/ — папка для архивов

## Примечания
- В проекте используется методология DRY (Don't Repeat Yourself) через файл `utils/json.go`, который содержит переиспользуемые функции `RespondWithJSON` и `RespondWithError` для стандартизированной обработки HTTP-ответов и ошибок во всех обработчиках.
