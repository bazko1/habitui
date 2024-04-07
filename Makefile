run:
	@go run cmd/main.go -debug
run-dev:
	@cp testdata/.habitui.json .
	@go run cmd/main.go -debug
test:
	go test ./...
gen-markdown:
	pandoc README.md > README.html
