package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/potatowhite/books/file-service/handler"
	"log"
	"os"
	"sync"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

type Consumer interface {
	Run() error
	Close()
}

type kafkaConsumer struct {
	consumer   *kafka.Consumer
	handler    handler.Handler
	workerPool map[int32]*worker
	wg         sync.WaitGroup
}

func (c *kafkaConsumer) Close() {
	logger.Println("Closing Kafka consumer")

	// unasign partitions
	c.consumer.Unassign()
	c.consumer.Close()
}

func NewConsumer(bootstrapServers string, groupId string, handler handler.Handler) (Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               bootstrapServers,
		"group.id":                        groupId,
		"auto.offset.reset":               "earliest",
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		return nil, err
	}

	_kafkaConsumer := &kafkaConsumer{
		consumer:   consumer,
		handler:    handler,
		workerPool: make(map[int32]*worker),
	}

	if err := consumer.Subscribe(handler.Topic(), rebalanceCb(_kafkaConsumer)); err != nil {
		return nil, err
	}

	return _kafkaConsumer, nil
}

func rebalanceCb(svcConsumer *kafkaConsumer) kafka.RebalanceCb {
	logger.Println("Kafka consumer rebalance callback")

	return func(c *kafka.Consumer, e kafka.Event) error {
		switch ev := e.(type) {
		case kafka.AssignedPartitions:
			svcConsumer.wg.Add(1)
			// make workers for each partition
			for _, partition := range ev.Partitions {
				worker := newWorker(int(partition.Partition), svcConsumer.handler)
				worker.start(&svcConsumer.wg)
				svcConsumer.workerPool[partition.Partition] = worker
			}

			c.Assign(ev.Partitions)
			logger.Printf("Kafka consumer assigned partitions: %v", ev.Partitions)
			svcConsumer.wg.Done()
		case kafka.RevokedPartitions:
			c.Unassign()
			log.Printf("Kafka consumer revoked partitions: %v", ev.Partitions)

			svcConsumer.wg.Add(1)
			// stop and remove worker for each partition
			for _, partition := range ev.Partitions {
				worker := svcConsumer.workerPool[partition.Partition]
				if worker != nil {
					worker.stop()
					delete(svcConsumer.workerPool, partition.Partition)
				}
			}

			svcConsumer.wg.Done()
		}

		return nil
	}
}

func (c *kafkaConsumer) Run() error {

	for {
		ev := c.consumer.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			if c.handler != nil {
				// Find the worker for this message's partition
				partition := e.TopicPartition.Partition
				worker := c.workerPool[partition]

				if worker != nil {
					worker.messages <- e
					continue
				}

				logger.Printf("no worker found for partition %d", partition)

				// nack message
				c.consumer.CommitMessage(e)

			}
		case kafka.PartitionEOF:
			continue
		case kafka.Error:
			log.Printf("Consumer error: %v (%v)", e, e.Code())
			if e.Code() == kafka.ErrAllBrokersDown {
				c.shutdownWorkers()
				return e
			}
		default:
			log.Printf("Ignored event: %s", e)
		}
	}
}

func (c *kafkaConsumer) stopWorkers() {
	for _, w := range c.workerPool {
		w.stop()
		<-w.done
	}
}

func (c *kafkaConsumer) shutdownWorkers() {
	log.Println("Shutting down workers...")

	// close all workers which are still running
	c.stopWorkers()

	// unassign partitions
	c.consumer.Unassign()

	// wait for all workers to finish
	c.wg.Wait()

	log.Println("Workers shut down.")
}
