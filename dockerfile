FROM golang:1.10 as builder


RUN apt-get update
RUN mkdir -p /go/src/bitbucket.org/yanndr/deptree
COPY . /go/src/bitbucket.org/yanndr/deptree
WORKDIR /go/src/bitbucket.org/yanndr/deptree/cmd/deptree
RUN go install

FROM debian:stretch-slim
COPY --from=builder /go/bin/deptree /bin/deptree
COPY ./cmd/deptree/data/ /data/
ENTRYPOINT ["/bin/deptree"]