# VK Golang assignment

## About
1. Реализованны все хендлеры
2. Написаны юнит тесты 

## Сборка и Запуск

### Сборка
```shell
make build
```

### Запуск
```shell
make run
```
Если запускаете приложение первый раз - отправьте запрос на /migrate-up, чтобы инициализировать таблицы базы данных

### Swagger UI
Чтобы просмотреть спецификацию - запустите приложение и откройте ссылку 
```URL
http://localhost:3000/swagger/index.html
```

## Usage

### Миграции

### POST: /migrate-up - создает sql таблицы

### POST: /migrate-down - удаляет sql таблицы

### POST: /migrate-force - если необходимо починить "грязные" версии миграций
```json
{
  "version": 1
}
```

### Аутентификация

### POST: /sign-up - регистрация
```json
{
    "name": "Katya",
    "username": "Katya@",
    "password": "qwerty"
}
```
```json
{
    "name": "Katya",
    "username": "admin_Katya@",
    "password": "qwerty"
}
```
#### response: 
```json
{
    "id": 1,
    "message": "User is created"
}
```
Если имя пользователя начинается на admin, то пользователь будет обладать правами администрора, иначе - пользователя

### POST: /sign-in - аутентификация
```json
{
    "username": "admin_Katya@",
    "password": "qwerty"
}
```
#### response: 
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA4MTAzMzYsImlhdCI6MTcxMDc2NzEzNiwidXNlcl9pZCI6MSwiVXNlclJvbGUiOiJhZG1pbiJ9.IstnvOXayUxlLVP8eEN1wvjVRWqTapoZsAkdGmEXvoY"
}
```
Возращается токен, котоырй нужно подставить в заголовок Authorization, Bearer Token

### Запросы обычного пользователя

### GET: /actor-list
Возвращает список всех актеров
#### response: 
```json
[
    {
        "actor_id": 1,
        "actor_name": "Alan Rickman",
        "actor_gender": "male",
        "actor_birthday": "1946-01-01",
        "movies": [
            {
                "Id": 2,
                "Name": "Сумерки 99",
                "Description": "Кино для абстрактных мужчин",
                "Date": "2012-10-10T00:00:00Z",
                "Rate": 0
            },
            {
                "Id": 3,
                "Name": "Сумерки 4",
                "Description": "Кино для настоящих мужчин",
                "Date": "2012-10-10T00:00:00Z",
                "Rate": 0
            }
        ]
    },
    {
        "actor_id": 2,
        "actor_name": "Skarlette Johnason",
        "actor_gender": "female",
        "actor_birthday": "1984-01-01",
        "movies": [
            {
                "Id": 3,
                "Name": "Сумерки 4",
                "Description": "Кино для настоящих мужчин",
                "Date": "2012-10-10T00:00:00Z",
                "Rate": 0
            }
        ]
    },
    {
        "actor_id": 3,
        "actor_name": "Jody Foster",
        "actor_gender": "female",
        "actor_birthday": "1989-01-01",
        "movies": [
            {
                "Id": 3,
                "Name": "Сумерки 4",
                "Description": "Кино для настоящих мужчин",
                "Date": "2012-10-10T00:00:00Z",
                "Rate": 0
            }
        ]
    }
]
```

### POST: /movie-list
Возвращает список всех фильмов
```json
{
  "sort_by": "a.birthday"
}
```
#### response: 
```json
[
    {
        "movie_id": 2,
        "movie_name": "Сумерки 99",
        "movie_description": "Кино для абстрактных мужчин",
        "movie_date": "2012-10-10T00:00:00Z",
        "movie_rate": 0,
        "actors": [
            {
                "id": 1,
                "name": "Alan Rickman",
                "gender": "male",
                "Birthday": "1946-01-01"
            }
        ]
    },
    {
        "movie_id": 3,
        "movie_name": "Сумерки 4",
        "movie_description": "Кино для настоящих мужчин",
        "movie_date": "2012-10-10T00:00:00Z",
        "movie_rate": 0,
        "actors": [
            {
                "id": 1,
                "name": "Alan Rickman",
                "gender": "male",
                "Birthday": "1946-01-01"
            },
            {
                "id": 2,
                "name": "Skarlette Johnason",
                "gender": "female",
                "Birthday": "1984-01-01"
            },
            {
                "id": 3,
                "name": "Jody Foster",
                "gender": "female",
                "Birthday": "1989-01-01"
            }
        ]
    }
]
```
Возвращает список фильмов со всеми полями.

### POST: /search-movies
Возвращает список фильмов, по отрывку фильма и имени актера. Оба поля могут 
быть пустыми
```json
{
    "movie_pattern": "Сум",
    "actor_pattern": "Val"
}
```
#### response: 
```json
[
    {
        "movie_id": 3,
        "movie_name": "Сумерки 4",
        "movie_description": "Кино для настоящих мужчин",
        "movie_date": "2012-10-10T00:00:00Z",
        "movie_rate": 0,
        "actors": [
            {
                "id": 4,
                "name": "Valera ZHMA",
                "gender": "male",
                "Birthday": "1970-01-01"
            },
            {
                "id": 5,
                "name": "Valera Kruglov",
                "gender": "male",
                "Birthday": "1970-01-01"
            }
        ]
    }
]
```

### Запросы администратора
Требуется токен

### POST: /add-actor
Добавляет актера
```json
{
    "name": "Valera Devyatkin",
    "gender": "male",
    "birthday": "1999-09-09"
}
```
#### response: 
```json
{
    "id": 5,
    "message": "Actor is added"
}
```
Возвращает id актера и сообщение. В случае провала - код и описание ошибки

### POST: /add-actor-to-movie
Добавляет актера в заданный фильм 
```json
{
    "actor_id": 5,
    "movie_id": 3
}
```
#### response: 
```json
{
    "message": "Actor 5 added to the movie 3"
}
```

### PUT: /change-actor-field
Изменяет заданное поле у актера
```json
{
  "actor_id": 1,
  "field": "name",
  "new_value": "Valera"
}
```
#### response: 
```json
{
    "message": "Field is changed"
}
```

### DELETE: /delete-actor
Удаляет актера из базы данных
```json
{
    "actor_id": 1
}
```
#### response: 
```json
{
    "message": "Actor is removed"
}
```

### POST: /add-movie
Добавляет фильм по названию, описанию, дате, рейтингу и списку актеров
```json
{
    "name": "Сумерки 4",
    "description": "Кино для настоящих мужчин",
    "date": "2012-10-10",
    "rate": 0,
    "actors": [1, 2, 3]
}
```
#### response: 
```json
{
    "id": 3,
    "message": "Movie is added"
}
```
Возвращает id добавленного фильма

### DELETE: /delete-movie
Удаляет фильма из базы данных
```json
{
    "movie_id": 1
}
```
#### response: 
```json
{
    "message": "Movie is removed"
}
```

### DELETE: /delete-movie-field
Удаляет поле у заданного фильма из базы данных
```json
{
    "movie_id": 2,
    "field": "date"
}
```
Удаление актера из фильма
```json
{
    "movie_id": 2,
    "actor_id": 1
}
```
### response:
```json
{
    "message": "Field is removed"
}
```

### PUT /change-movie-field
Меняет поле у заданного фильма
```json
{
    "movie_id": 1,
    "field": "name",
    "new_value": "Sumerki"
}
```
### response:
```json
{
    "message": "Field is changed"
}
```
