up: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
	docker-compose up --build -V

down: 
	docker-compose down
	rm main

restart: down up
