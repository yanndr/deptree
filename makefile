build:
	go build ./cmd/deptree

docker:
	docker build -t deptree .