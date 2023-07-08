# emailTask

## Заполнение конфигов:
1. Использовать configs/configs.go для заполнения конфигов

## Запуск сервиса локально:
1. Запуск сервера
```shell
  make createdb
  go run cmd/main.go
```
2. Миграции в базу данных
```shell
  make migrate-up
```

## Запуск сервиса через Docker:
1. Запуск Docker-compose
```shell
  make docker-up
```
2. Миграции в базу данных
```shell
  make migrate-up
```

По умолчанию запуск на localhost:8080

## Пример запросов к сервису:

POST /add-user-mail?mail=muhammed_mail@gmail.com
Response:
```
    e160ea36-396a-48bb-b662-b3093bf5ce5e
```

GET /get-user-mail?id=e160ea36-396a-48bb-b662-b3093bf5ce5e
Response:
```
    --возвращает zip файл с файлами--
```

