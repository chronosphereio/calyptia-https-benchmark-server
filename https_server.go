package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var c metrics.Counter

var printRecords bool

type handler struct{}

var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fluent_records_total",
			Help: "Number of received records",
		})
)

func defaultURI(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	for {
		var doc map[string]interface{}

		err := dec.Decode(&doc)
		if err == io.EOF {
			break
		}
		if printRecords {
			fmt.Println("parsed record:", doc)
		}
		counter.Inc()
		c.Inc(1)
	}
	w.Write([]byte("{\"status\":\"ok\",\"errors\":false}"))
}

func main() {
	printMetrics := flag.Bool("printmetrics", false, "print metrics to the console")
	flag.BoolVar(&printRecords, "printrecords", true, "print request records")
	flag.Parse()

	c = metrics.NewCounter()
	metrics.Register("records", c)

	if *printMetrics {
		go metrics.Log(metrics.DefaultRegistry, 1*time.Second,
			log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}

	caCert, err := ioutil.ReadFile("client.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		//ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: caCertPool,
	}
	srv := &http.Server{
		Addr:      ":8443",
		TLSConfig: cfg,
	}

	prometheus.MustRegister(counter)

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/").HandlerFunc(defaultURI)
	http.Handle("/", r)

	log.Print("Starting server at https://127.0.0.1:8443")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
