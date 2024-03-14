build:
	docker compose build filmoteka

run:
	docker compose up filmoteka

swag:
	swag init -g cmd/main.go