package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redpanda-consumer-failover/pkg/model/config"
)

func main() {
	workerName := flag.String("name", "producer", "specify worker name")
	flag.Parse()

	cfg, err := config.InitConfig(*workerName)
	if err != nil {
		log.Fatal().Err(err).Msg("[main] failed to load config")
	}
	defer cfg.LogFile.Close()
}
