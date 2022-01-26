package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/linkedin/goavro"
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

var codec *goavro.Codec
var printRecords bool

type handler struct{}

var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fluent_records_total",
			Help: "Number of received records",
		})
)

func avroHandler(w http.ResponseWriter, r *http.Request) {
	bodyData, _ := ioutil.ReadAll(r.Body)
	record, _, err := codec.NativeFromBinary(bodyData)
	if err != nil {
		fmt.Println("error parsing avro data:", err)
	}
	if printRecords {
		fmt.Println("parsed record:", record)
	}
	counter.Inc()
	c.Inc(1)
	w.Write([]byte("{\"status\":\"ok\",\"errors\":false}"))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
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
	avroSchema := flag.String("avro-schema-path", "", "specify file path containing avro schema. this disables json parsing.")
	flag.BoolVar(&printRecords, "printrecords", true, "print request records")
	flag.Parse()

	if *avroSchema != "" {
		schemaData, err := ioutil.ReadFile(*avroSchema)
		if err != nil {
			log.Fatal(err)
		}
		avroCodec, err := goavro.NewCodec(string(schemaData))
		if err != nil {
			log.Fatal(err)
		}
		codec = avroCodec
	}

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
	if codec != nil {
		r.PathPrefix("/").HandlerFunc(avroHandler)
	} else {
		r.PathPrefix("/").HandlerFunc(jsonHandler)
	}
	http.Handle("/", r)

	log.Print("Starting server at https://127.0.0.1:8443")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
