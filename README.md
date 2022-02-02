# Fluent HTTPs Benchmark Server

This is a HTTPs Server that measure the numbers of JSON records received per request. The metrics are printed to the terminal  interface and also exposed through a Prometheus endpoint ```/metrics```.

The service works in the TCP/TLS port ```8443```

## Install dependencies

Install Go dependencies

```bash
go mod download
```

## Build and Run the Server

```sh
go build
./https-benchmark-server
```

## Scrape Metrics with Prometheus

This repository provides a Prometheus configuration file ([prometheus-config.yml](prometheus-config.yml)) that can be used to scrape the metrics from the service.

Just run Prometheus from the command line:

```
prometheus --config.file=prometheus-config.yaml
```

To query the data in Prometheus dashboard:

- Go to http://127.0.0.1:9090
- Write the query: ```irate(fluent_records_total[1m])```
- Click on the Graph Tab

## Ingesting data with Fluent Bit

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
$ ./https-benchmark-server
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
