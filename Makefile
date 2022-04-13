#!/bin/bash

build_producer:
	@echo ">> build producer binary ..."
	@go build -o ./cmd/data-loader/producer ./cmd/data-loader

build_consumer:
	@echo ">> build consumer binary ..."
	@go build -o ./cmd/consumer/consumer ./cmd/consumer

run_producer:
	@./cmd/consumer

#docker run -p 0.0.0.0:18080:8080 --network=host -e KAFKA_BROKERS=0.0.0.0:61424 quay.io/cloudhut/kowl:master