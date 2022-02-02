FROM golang:1.17

WORKDIR /go/src/https-benchmark-server
COPY . .

RUN go mod download
RUN go build
RUN go install

CMD ["https-benchmark-server"]
