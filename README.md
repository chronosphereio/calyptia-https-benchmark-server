# HTTPs Benchmark Server

This is a simple HTTPs Server that meassure the numbers of JSON records received per second.
Upon start it listens on TCP/TLS port ```8443``` .


## Running the Server

```sh
go run https_server.go
```

## Testing with Fluent Bit

Start the HTTPs Server in one terminal, in the other start Fluent Bit

```
fluent-bit -i dummy                  \
              -p rate=100000         \
           -o http://127.0.0.1:8443  \
              -p tls=on              \
              -p tls.verify=off      \
              -p workers=2           \
           -f 1
```
