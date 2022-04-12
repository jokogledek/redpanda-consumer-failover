package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ujunglangit-id/redpanda-consumer-failover/pkg/model/types"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Kafka   KafkaConfig `yaml:"kafka"`
	LogFile *os.File    `yaml:"-"`
}

type KafkaConfig struct {
	Brokers       []string `yaml:"brokers"`
	Topic         string   `yaml:"topic"`
	ConsumerGroup string   `yaml:"consumer_group"`
}

func InitConfig(logPrefix string) (cfg *Config, err error) {
	cfg = &Config{}
	cfg.initLog(logPrefix)
	err = cfg.readFile()
	return
}

func (cfg *Config) initLog(prefix string) {

	//initialize log format & log file output
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logFile, err := os.OpenFile(prefix+"_"+types.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("error opening file")
	}
	log.Logger = zerolog.New(logFile).With().Timestamp().Caller().Logger()
}

func (cfg *Config) readFile() (err error) {
	var (
		f *os.File
	)

	path := []string{
		"",
		"files/config",
		"./files/config",
		"../files/config",
		"../../files/config",
	}

	for _, val := range path {
		f, err = os.Open(fmt.Sprintf(`%s/cfg.yml`, val))
		if err == nil {
			log.Info().Msgf("[config][readFile] load config file from %s", fmt.Sprintf(`%s/cfg.yml`, val))
			decoder := yaml.NewDecoder(f)
			err = decoder.Decode(cfg)
			break
		}
	}

	if err != nil {
		return
	}

	log.Info().Msg("[config][readFile] Config load success")
	return
}
