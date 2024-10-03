build:
	CGO_ENABLED=0 go build -ldflags="-s -w"

test:
	go test -v

lint:
	go vet
	go fmt

update:
	go get -u
	go mod tidy
