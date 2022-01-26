FROM golang:1.15

RUN go get github.com/gorilla/mux
RUN go get github.com/rcrowley/go-metrics
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/linkedin/goavro

WORKDIR /go/src/https-benchmark-server
COPY . .

RUN go build
RUN go install

CMD ["https-benchmark-server"]
