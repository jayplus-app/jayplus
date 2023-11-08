# Default target to run when just executing `make`
.DEFAULT_GOAL := up

up:
	docker-compose up

up-detached:
	docker-compose up -d

up-backend:
	docker-compose up backend

up-frontend:
	docker-compose up frontend
