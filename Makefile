gen:
	@echo "Building gqlgen"
	@go get github.com/99designs/gqlgen
	@go run github.com/99designs/gqlgen generate


build:
	@echo "Building server"
	@go build -o bin/server cmd/main.go

run:
	@echo "Running server"
	@go run cmd/main.go