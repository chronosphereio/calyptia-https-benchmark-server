FROM golang:1.19.1 as builder

WORKDIR /calyptia/https-benchmark-server

# Cache go module setup
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Now copy and build source
COPY . .
RUN go build -o https-benchmark-server

CMD [ "./https-benchmark-server" ]
