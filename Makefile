build:
	docker volume create pg_data
	docker-compose -f docker-compose.yml up --build

rebuild:
	docker-compose down
	docker-compose -f docker-compose.yml up --build

delete:
	docker-compose down
	docker volume rm pg_data
	docker image rm api
	docker image rm postgres