# 1C

Потапов Даниил Аликович

Задача номер 3

Запуск.

Поднимаем базу данных:

docker-compose up -d

Выполняем запросы queries.sql

Инициализируем переменные окружение:

export SERVER_ADDRESS=localhost:8080

export POSTGRES_CONN=postgres://postgres:postgres@localhost:5432/postgres

Запускаем программу:

go run main.go
