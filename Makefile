build:
	docker-compose -f docker-compose.yml up --build

rebuild:
	docker-compose down
	docker-compose -f docker-compose.yml up --build