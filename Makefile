run:
	go run main.go

build: 
	go build -o marketWatch main.go

lint:
	golangci-lint run --fix

serve:
	marketWatch serve

fetch-trends: 
	marketWatch fetch-trends
