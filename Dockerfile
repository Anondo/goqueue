FROM golang:alpine

RUN apk add --no-cache --update git

ENV GOPATH=/go

RUN go get -u github.com/golang/dep/cmd/dep

ADD . $GOPATH/src/goqueue

RUN cd $GOPATH/src/goqueue/ && dep init && dep ensure -v

RUN go install -v $GOPATH/src/goqueue


FROM alpine:latest

WORKDIR app

COPY --from=0 /go/bin/goqueue /app

COPY --from=0 /go/src/goqueue/config.yml /app

ENTRYPOINT ["./goqueue"]

CMD ["server"]
