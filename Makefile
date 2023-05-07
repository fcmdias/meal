run: 
	rm main 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
	docker build -t my-go-app .
	docker run -p 8080:8080 my-go-app