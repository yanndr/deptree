build:
	go build ./cmd/deptree

docker:
	docker build -t deptree .

test:
	go test -race  -tags="unit integration" ./...

test-unit:
	go test -race -tags unit ./...

test-integration:
	go test -race -tags integration ./...