run:
	go run main.go

build: 
	go build -o marketWatch main.go

lint:
	golangci-lint run --fix