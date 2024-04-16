run:
	@go run cmd/tui/main.go
run-dev:
	@cp testdata/.habitui.json .
	@go run cmd/tui/main.go
serve:
	go run cmd/server/main.go
test:
	go test ./...
gen-markdown:
	pandoc README.md > README.html
