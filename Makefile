#!/bin/bash

build_producer:
	@echo ">> build producer binary ..."
	@GOOS=linux go build -o ./cmd/data-loader/producer ./cmd/data-loader

build_consumer:
	@echo ">> build consumer binary ..."
	@GOOS=linux  go build -o ./cmd/consumer/nix_consumer ./cmd/consumer
	@GOOS=linux  go build -o ./cmd/uncomitted_consumer/nix_c_consumer ./cmd/uncomitted_consumer
	@go build -o ./cmd/consumer/consumer ./cmd/consumer
	@go build -o ./cmd/uncomitted_consumer/c_consumer ./cmd/uncomitted_consumer

run_producer:
	@./cmd/consumer

upload_gcp:
	@gcloud compute scp  ./cmd/consumer/consumer benchmark-openmessaging-1:/home/m.azwarnurrosat/sandbox/failover/consumer --tunnel-through-iap   --zone asia-southeast2-a
    @gcloud compute scp  ./cmd/uncomitted_consumer/c_consumer benchmark-openmessaging-1:/home/m.azwarnurrosat/sandbox/failover/c_consumer --tunnel-through-iap   --zone asia-southeast2-a
    @gcloud compute scp  ./files/config/cfg.yml benchmark-openmessaging-1:/home/m.azwarnurrosat/sandbox/failover/cfg.yml --tunnel-through-iap   --zone asia-southeast2-a