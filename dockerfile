FROM  golang:alpine

COPY . /go/src/github.com/okex/infura-service

WORKDIR /go/src/github.com/okex/infura-service

RUN go build -o infura

CMD ["./infura"]



