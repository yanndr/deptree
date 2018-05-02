build:
	go build ./cmd/deptree

install:
	go install ./cmd/deptree

docker:
	docker build -t deptree .

test:
	go test -race -v  ./...
