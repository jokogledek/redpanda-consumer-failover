package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redpanda-consumer-failover/internal/consumer"
	"github.com/ujunglangit-id/redpanda-consumer-failover/internal/model/config"
	"os"
)

func main() {
	workerName := flag.String("name", "nocomit_consumer", "specify worker name")
	flag.Parse()

	cfg, err := config.InitConfig(*workerName)
	if err != nil {
		log.Fatal().Err(err).Msg("[main] failed to load config")
	}
	defer cfg.LogFile.Close()

	consumer.NewConsumer(*workerName, cfg, true).InitConsumer()
	os.Exit(1)
}
