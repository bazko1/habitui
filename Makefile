run:
	@go run cmd/main.go
run-dev:
	@cp testdata/.habitui.json .
	@go run cmd/main.go
test:
	go test ./...
