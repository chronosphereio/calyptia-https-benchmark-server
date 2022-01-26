build:
	docker build -t https-benchmark-server .

run-json: build
	docker run -ti -p8443:8443 https-benchmark-server https-benchmark-server

run-avro: build
	docker run -ti -p8443:8443 https-benchmark-server https-benchmark-server -avro-schema-path avro-schema.json

.PHONY: build run-json run-avro
