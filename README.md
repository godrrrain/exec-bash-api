# Тестовое задание на стажировку PGStart 2024

## Описание
Сервис представляет собой REST API для параллельного запуска и хранения команд.

В рамках данного проекта используются следующие определения
- Команда - bash скрипт (экранированный)
- Обычная команда - команда, вывод которой записывается в БД после ее завершения
- Долгая команда - команда, вывод которой записывается в БД по мере ее выполнения

SQL скрипты находятся в директории `postgres`.

Коллекция интеграционных тестов - [Script_Execution_System.postman_collection.json](postman/[pgstart]Script_Execution_System.postman_collection.json)

Спецификация Open API - [open-api.yml](docs/open-api.yml)

Операционная система `Ubuntu 22.04`

## Запуск проекта
1. 💻 Скачай проект
2. ✅ Запусти проект: `docker compose up`

## Решения, принятые во время разработки

### 1. Остановка команды

#### Контекст
Необходимо реализовать метод остановки команды, который должен поддерживать параллельные вызовы.

#### Решение
Было принято решение хранить информацию о запущенных процессах in-memory. Добавлять информацию о созданных процессах в момент их создания и удалять эту информацию после их завершения.

Для этого был создан в `aprocess.go` тип, представляющий map-у активных процессов, использующий RWMutex.

#### Последствия
Это решение позволило реализовать метод остановки команды, который можно безопасно вызывать параллельно.


### 2. Обычные и долгие команды

#### Контекст
Необходимо реализовать поддержку создания долгих команд, вывод которых записывается в БД по мере ее выполнения.

Долгие команды записывают в БД каждую новую строку вывода, следовательно, такие команды гораздо чаще обращаются к БД, чем обычные команды, тем самым увеличивая нагрузку.

#### Решение
Было принято решение создать отдельные методы для запуска обычных и долгих команд.

#### Последствия
Это решение позволило снизить общую нагрузку на систему, так как часто имеет смысл вызывать обычную команду, промежуточный вывод которой нет необходимости хранить.


## Описание методов

### Создание обычной команды
Запускает переданную bash-команду, сохраняет результат выполнения в БД.

#### POST `http://localhost:8080/api/v1/commands`
#### cURL
```
curl -X 'POST' \
  'http://localhost:8080/api/v1/commands' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Тест успешной команды",
  "script": "echo Test"
}'
```
#### Body
```
{
  "description" : "Тест успешной команды",
  "script": "echo Test"
}
```
#### Response
```
{
    "uuid": "9bc12250-78d1-4cc1-9ba2-a36d1685a007"
}
```

### Получение одной команды
По переданному в path `uuid` возвращает информацию о команде.

#### GET `http://localhost:8080/api/v1/commands/:uuid`
#### Response
```
{
    "command_uuid": "9bc12250-78d1-4cc1-9ba2-a36d1685a007",
    "description": "Тест 1",
    "script": "echo Test",
    "status": "EXECUTED",
    "output": "Test\n"
}
```

### Получение списка команд
Возвращает информацию о всех командах.

Поддерживает следующие параметры в query:

- `status` - фильтрация по статусу команды

- `limit` - ограничение на максимальное количество элементов

- `offset` - пропустить указанное число элементов

#### GET `http://localhost:8080/api/v1/commands`
#### cURL
```
curl -X 'GET' \
  'http://localhost:8080/api/v1/commands' \
  -H 'accept: application/json'
```
#### Response
```
[
    {
        "command_uuid": "9bc12250-78d1-4cc1-9ba2-a36d1685a007",
        "description": "Тест 1",
        "script": "echo Test",
        "status": "EXECUTED",
        "output": "Test\n"
    },
    {
        "command_uuid": "4fa5e938-3d95-4051-9376-dbbb9510a49c",
        "description": "Тест долгой команды",
        "script": "for i in $(seq 1 5); do echo Test; sleep 3; done",
        "status": "EXECUTING",
        "output": "Test\nTest\nTest\nTest\n"
    }
]
```

### Остановка команды
По переданному в path `uuid` останавливает команду.

#### PATCH `http://localhost:8080/api/v1/commands/:uuid`
#### cURL
```
curl -X 'PATCH' \
  'http://localhost:8080/api/v1/commands/4fa5e938-3d95-4051-9376-dbbb9510a49c' \
  -H 'accept: application/json'
```
#### Response
```
{
    "message": "command stopped successfully"
}
```

### Создание долгой команды
Запускает переданную bash-команду, сохраняет вывод команды в БД по мере ее выполнения.

#### POST `http://localhost:8080/api/v1/durables/commands`
#### cURL
```
curl -X 'POST' \
  'http://localhost:8080/api/v1/durables/commands' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Тест долгой команды",
  "script": "for i in $(seq 1 5); do echo Test; sleep 3; done"
}'
```
#### Body
```
{
  "description" : "Тест долгой команды",
  "script": "for i in $(seq 1 5); do echo Test; sleep 3; done"
}
```
#### Response
```
{
    "uuid": "4fa5e938-3d95-4051-9376-dbbb9510a49c"
}
```

### Удаление команды
По переданному в path `uuid` удаляет информацию о команде.

#### DELETE `http://localhost:8080/api/v1/commands/:uuid`
#### cURL
```
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/commands/9bc12250-78d1-4cc1-9ba2-a36d1685a007' \
  -H 'accept: */*'
```

### Проверка состояния сервера
Возвращает `Status OK`

#### GET `http://localhost:8080/manage/health`
#### cURL
```
curl -X 'GET' \
  'http://localhost:8080/manage/health' \
  -H 'accept: */*'
```