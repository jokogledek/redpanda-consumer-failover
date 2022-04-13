package producer

import (
	"github.com/ujunglangit-id/redpanda-consumer-failover/internal/model/config"
	"github.com/ujunglangit-id/redpanda-consumer-failover/internal/model/types"
)

type Producer struct {
	Config      *config.Config
	WorkerName  string
	PayloadData []*types.Stocks
}

func NewProducer(name string, cfg *config.Config) *Producer {
	return &Producer{
		Config:     cfg,
		WorkerName: name,
	}
}

func (p *Producer) LoadDataSource() {

}

func (p *Producer) Run() {

}
