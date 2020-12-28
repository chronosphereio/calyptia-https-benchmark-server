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
              -p format=json_lines   \
              -p tls=on              \
              -p tls.verify=off      \
              -p workers=2           \
           -f 1
```

## Expected Output

In ```https-benchmark-server``` stderr you should see some stats like this:


```
$ ./https_benchmark_server 
2020/12/27 23:11:55 Starting server at https://127.0.0.1:8443
metrics: 23:11:58.164353 counter records
metrics: 23:11:58.164377   count:           39266
metrics: 23:11:59.164347 counter records
metrics: 23:11:59.164363   count:          138395
metrics: 23:12:00.164288 counter records
metrics: 23:12:00.164307   count:          230671
metrics: 23:12:01.164367 counter records
metrics: 23:12:01.164385   count:          320348
metrics: 23:12:02.164410 counter records
metrics: 23:12:02.164424   count:          417787
metrics: 23:12:03.164350 counter records
metrics: 23:12:03.164367   count:          515907
metrics: 23:12:04.164378 counter records
metrics: 23:12:04.164392   count:          615270

```
