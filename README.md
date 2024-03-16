# VK GO assignment

## TODO

## !!!ПОМЕНЯТЬ НАСТРОЙКИ НА ПРОД!!!

### API:
- [X] Добавить актеров в /add-movie
- [X] Добавить сортировку в /movie-list
- [X] Получение списка актеров /actor-list
- [X] Добавление актера /add-actor
- [X] Удаление Актера /delete-actor
- [ ] Удаление поля у Актера /delete-actor-field
- [X] Удаление поля у Фильма /delete-movie-field
- [ ] Добавить функции валидации
- [X] Добавить Авторизацию и Аутентификацию
- [X] Добавить роль пользователя в tokenClaims и через нее определять права доступа
- [ ] Переписать ошибки на формат json
- [ ] Убрать поле id у актера из доступа

### DevOPS:
- [ ] Описать хендлеры в swagger
- [ ] Добавить запуск миграций, если нет таблиц


### Остальное: 
- [ ] Написать тесты
- [ ] Описать функции в readme.md (toggle line)

### Total test coverage
go test -v -coverpkg=./... -coverprofile=profile.cov ./...

go tool cover -func profile.cov
