package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/linkedin/goavro"
	"github.com/rcrowley/go-metrics"
	"io/ioutil"
	"log"
	"net/http"
)

var c metrics.Counter

var logEventCodec *goavro.Codec
var metaCodec *goavro.Codec

type handler struct{}

const (
	logEventAvroSchema = `{  "type": "record",
                               "name": "LogEvent",
                               "fields": [{"name": "metadata", "type": "bytes"},
                                          {"name": "avro_schema", "type": "string"},
                                          {"name": "max_size", "type": "int"},
                                          {"name": "payload", "type":{"type":"array","items":"bytes"}}]
                             }`
	metaSchema = `
        { "type":"record",
          "name":"metadata",
          "fields": [
  	         {"name": "wd_platform", "type":"string"},{"name": "wd_env_physical", "type":"string"},{"name": "wd_dc_physical", "type":"string"},
			 {"name": "wd_env_logical", "type":"string"},{"name": "wd_service", "type":"string"},{"name": "wd_owner", "type":"string"},
             {"name": "wd_datatype", "type":"string"},{"name": "wd_objectname", "type":"string"},{"name": "wd_solas", "type":"string"},
			 {"name": "swh_server", "type":"string"},{"name": "wd_service_instance", "type":"string"}
	      ]
	    } `
)

func avroHandler(w http.ResponseWriter, r *http.Request) {
	bodyData, _ := ioutil.ReadAll(r.Body)
	logEvent, _, err := logEventCodec.NativeFromBinary(bodyData)
	if err != nil {
		log.Fatal(err)
	}
	logEv := logEvent.(map[string]interface{})
	recordSchemaData, ok := logEv["avro_schema"]
	if !ok {
		w.Write([]byte("{\"status\":\"fail\",\"errors\":true}"))
		return
	}
	recordCodec, err := goavro.NewCodec(string(recordSchemaData.(string)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n\n==========logevent===========")
	fmt.Println("==========headers===========")
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	fmt.Println("==========payload===========")
	payloads := logEv["payload"].([]interface{})
	fmt.Println("records:")
	for _, p := range payloads {
		payload := p.([]uint8)
		record, _, _ := recordCodec.NativeFromBinary(payload)
		fmt.Println(record)
	}
	metadata, _, err := metaCodec.NativeFromBinary(logEv["metadata"].([]uint8))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("metadata:")
	fmt.Println(metadata)
	w.Write([]byte("{\"status\":\"ok\",\"errors\":false}"))
}

func main() {
	if avroCodec, err := goavro.NewCodec(string(logEventAvroSchema)); err == nil {
		logEventCodec = avroCodec
	} else {
		log.Fatal(err)
	}

	if avroCodec, err := goavro.NewCodec(string(metaSchema)); err == nil {
		metaCodec = avroCodec
	} else {
		log.Fatal(err)
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

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(avroHandler)
	http.Handle("/", r)

	log.Print("Starting server at https://127.0.0.1:8443")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
