build:
	docker build -t https-benchmark-server .

run: build
	docker run -ti -p8443:8443 https-benchmark-server https-benchmark-server

.PHONY: build run
