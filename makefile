build:
	go build ./cmd/deptree

install:
	go install ./cmd/deptree

docker:
	docker build -t deptree .

test:
	go test -race -v  ./...
benchmark-trace:
	go test -trace trace.out -benchmem -run=^$ bitbucket.org/yanndr/deptree -bench ^BenchmarkResolveScale$ 
