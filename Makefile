run:
	go run main.go

build:
	go build -o marketWatch main.go

lint:
	golangci-lint run --fix

serve:
	./marketWatch serve -p 8080 -c ./config.yaml

refresh-trends:
	./marketWatch refresh-trends -c ./config.yaml
