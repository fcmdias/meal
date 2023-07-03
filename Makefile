# ========================================
# Local Deployment

up: 
	go mod tidy 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
	docker-compose up --build -V

down: 
	docker-compose down
	rm main

restart: down up

sample-data:
	curl -X POST localhost:8080/recipes/savemany -d @data/recipes.json --header "Content-Type: application/json"

query-recipes:
	curl -X GET localhost:8080/recipes

# ========================================
# AWS Deployment 