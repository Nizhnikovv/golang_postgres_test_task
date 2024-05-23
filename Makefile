build:
		@go build -o bin/goposg
run: build
		@./bin/goposg
test:
		@go test -v ./...