
run:
	@go run cmd/habitui/main.go
run-dev:
	@cp testdata/.habitui.json .
	@go run cmd/habitui/main.go
run-remote:
	go run cmd/habitui/main.go -remote-password "test" -remote-user "foo" -remote-server "http://localhost:3000" -enable-remote
serve:
	go run cmd/server/main.go
serve-sqlite:
	go run cmd/server/main.go -engine sqlite
test:
	go test -count=1 ./...
test-race:
	go test -count=1 -race ./...
gen-markdown:
	pandoc README.md > README.html
gen-testdata:
	cd ./testdata/ && go run main.go
