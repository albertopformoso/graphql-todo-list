run:
	docker-compose -f ./docker-compose.yml up -d --build

rm:
	docker-compose -f ./docker-compose.yml down