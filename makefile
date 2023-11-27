build:
	go build -o cloud-platform-label-pods-bin .

run:
	go run cloud-platform-label-pods-bin .

test:
	go test -race -covermode=atomic -coverprofile=c.out -v ./...

coverage:
	make test
	go tool cover -html=c.out

fmt:
	go fmt ./...

