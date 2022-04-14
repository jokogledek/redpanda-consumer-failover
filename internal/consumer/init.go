package consumer

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/ujunglangit-id/redpanda-consumer-failover/internal/model/config"
	"math/rand"
	"sync"
	"time"
)

type tp struct {
	t string
	p int32
}

type Consumer struct {
	Config     *config.Config
	WorkerName string
	Client     *kgo.Client
	consumers  map[tp]*pconsumer
	noCommit   bool
}

type pconsumer struct {
	cl        *kgo.Client
	topic     string
	partition int32

	quit    chan struct{}
	done    chan struct{}
	recs    chan []*kgo.Record
	noComit bool
}

func NewConsumer(name string, cfg *config.Config, noCommit bool) *Consumer {
	return &Consumer{
		Config:     cfg,
		WorkerName: name,
		consumers:  make(map[tp]*pconsumer),
		noCommit:   noCommit,
	}
}

func (c *Consumer) InitConsumer() {
	opts := []kgo.Opt{
		kgo.SeedBrokers(c.Config.Kafka.Brokers...),
		kgo.ConsumerGroup(c.Config.Kafka.ConsumerGroup),
		kgo.ConsumeTopics(c.Config.Kafka.Topic),
		kgo.OnPartitionsAssigned(c.assigned),
		kgo.OnPartitionsRevoked(c.lost),
		kgo.OnPartitionsLost(c.lost),
		kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
		kgo.SessionTimeout(time.Second * 6),
		kgo.RebalanceTimeout(time.Second * 10),
		kgo.HeartbeatInterval(time.Second * 2),
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	if err = cl.Ping(context.Background()); err != nil { // check connectivity to cluster
		panic(err)
	}

	c.poll(cl)
}

func (c *Consumer) assigned(_ context.Context, cl *kgo.Client, assigned map[string][]int32) {
	for topic, partitions := range assigned {
		for _, partition := range partitions {
			pc := &pconsumer{
				cl:        cl,
				topic:     topic,
				partition: partition,

				quit:    make(chan struct{}),
				done:    make(chan struct{}),
				recs:    make(chan []*kgo.Record, 10),
				noComit: c.noCommit,
			}
			c.consumers[tp{topic, partition}] = pc
			go pc.consume()
		}
	}
}

// In this example, each partition consumer commits itself. Those commits will
// fail if partitions are lost, but will succeed if partitions are revoked. We
// only need one revoked or lost function (and we name it "lost").
func (c *Consumer) lost(_ context.Context, cl *kgo.Client, lost map[string][]int32) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for topic, partitions := range lost {
		for _, partition := range partitions {
			tp := tp{topic, partition}
			pc := c.consumers[tp]
			delete(c.consumers, tp)
			close(pc.quit)
			fmt.Printf("waiting for work to finish t %s p %d\n", topic, partition)
			wg.Add(1)
			go func() { <-pc.done; wg.Done() }()
		}
	}
}

func (c *Consumer) poll(cl *kgo.Client) {
	for {
		// PollRecords is strongly recommended when using
		// BlockRebalanceOnPoll. You can tune how many records to
		// process at once (upper bound -- could all be on one
		// partition), ensuring that your processor loops complete fast
		// enough to not block a rebalance too long.
		fetches := cl.PollRecords(context.Background(), 10000)
		if fetches.IsClientClosed() {
			return
		}
		fetches.EachError(func(_ string, _ int32, err error) {
			// Note: you can delete this block, which will result
			// in these errors being sent to the partition
			// consumers, and then you can handle the errors there.
			panic(err)
		})
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			tp := tp{p.Topic, p.Partition}

			// Since we are using BlockRebalanceOnPoll, we can be
			// sure this partition consumer exists:
			//
			// * onAssigned is guaranteed to be called before we
			// fetch offsets for newly added partitions
			//
			// * onRevoked waits for partition consumers to quit
			// and be deleted before re-allowing polling.
			c.consumers[tp].recs <- p.Records
		})
		cl.AllowRebalance()
	}
}
func (pc *pconsumer) consume() {
	defer close(pc.done)
	fmt.Printf("Starting consume for  t %s p %d\n", pc.topic, pc.partition)
	defer fmt.Printf("Closing consume for t %s p %d\n", pc.topic, pc.partition)
	for {
		select {
		case <-pc.quit:
			return
		case recs := <-pc.recs:
			time.Sleep(time.Duration(rand.Intn(150)+100) * time.Millisecond) // simulate work
			fmt.Printf("Some sort of work done, about to commit t %s p %d\n", pc.topic, pc.partition)
			if pc.noComit {
				for _, v := range recs {
					fmt.Printf("ignore commit for data : %s\n", string(v.Value))
				}
			} else {
				err := pc.cl.CommitRecords(context.Background(), recs...)
				if err != nil {
					fmt.Printf("Error when committing offsets to kafka err: %v t: %s p: %d offset %d\n", err, pc.topic, pc.partition, recs[len(recs)-1].Offset+1)
				} else {
					for _, v := range recs {
						fmt.Printf("done commit for data : %s\n", string(v.Value))
					}
				}
			}
		}
	}
}
