obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

agg:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

docker-up:
	@sudo docker-compose up -d

docker-down:
	@sudo docker-compose down

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
	@go mod tidy

gate:
	@go build -o bin/gate gateway/main.go
	@./bin/gate


.PHONY: obu invoicer