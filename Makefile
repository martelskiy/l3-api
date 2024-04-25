BINARY_PATH=${EXEC_DIRECTORY}/bin/l3-api
EXEC_DIRECTORY=cmd/l3-api

.PHONY: default
default: build run 

.PHONY: build
build:
	go build -o ${BINARY_PATH} ${EXEC_DIRECTORY}/main.go

.PHONY: run
run:
	./${BINARY_PATH}

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test ./...

.PHONY: test_coverage
test_coverage:
	go test ./... -coverprofile=coverage.out

.PHONY: dep
dep:
	go mod download

.PHONY: vet
vet:
	go vet ./...

.PHONY: swag-init
swag-init:
	swag init -g ${EXEC_DIRECTORY}/main.go -o ./api/docs