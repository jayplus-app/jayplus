# Default target to run when just executing `make`
.DEFAULT_GOAL := up

up:
	docker-compose up --build

up-detached:
	docker-compose up -d --build

down:
	docker-compose down && docker-compose rm -f

up-backend:
	docker-compose up backend --build

up-frontend:
	docker-compose up frontend --build

up-website:
	docker-compose up website --build

up-nginx:
	docker-compose up nginx --build