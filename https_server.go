package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
        "github.com/rcrowley/go-metrics"
        "time"
        "os"
        "io"
        //"fmt"
        "encoding/json"
)

var c metrics.Counter

func main() {
        c  = metrics.NewCounter()
        metrics.Register("records", c)

        go metrics.Log(metrics.DefaultRegistry, 1 * time.Second,
                       log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	caCert, err := ioutil.ReadFile("client.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		//ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}
	srv := &http.Server{
		Addr:      ":8443",
		Handler:   &handler{},
		TLSConfig: cfg,
	}
        log.Print("Starting server at https://127.0.0.1:8443")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var records int64 = 0
        dec := json.NewDecoder(req.Body)
        for {
            var doc string

            err := dec.Decode(&doc)
            if err == io.EOF {
               break
            }
            records += 1
       }
       c.Inc(records)
       w.Write([]byte("PONG"))
}
