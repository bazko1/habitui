default: build
run:
	@go run cmd/habitui/main.go
run-dev:
	@cp testdata/.habitui.json .
	@go run cmd/habitui/main.go
run-remote:
	go run cmd/habitui/main.go -remote-password "test" -remote-user "foo" -remote-server "http://localhost:3000" -enable-remote
serve:
	go run cmd/habitui-server/main.go
serve-sqlite:
	go run cmd/habitui-server//main.go -engine sqlite
test:
	go test -count=1 ./...
test-race:
	go test -count=1 -race ./...
gen-markdown:
	pandoc README.md > README.html
gen-testdata:
	cd ./testdata/ && go run main.go
build-client:
	go build -o habitui cmd/habitui/main.go
build-server:
	go build -o habitui-server cmd/habitui-server/main.go
build: build-server build-client
clean:
	rm ./habitui
	rm ./habitui-server
