package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/potatowhite/books/file-service/handler"
	"log"
	"sync"
)

type Consumer interface {
	Run() error
	Close()
}

type kafkaConsumer struct {
	consumer   *kafka.Consumer
	handler    handler.Handler
	workerPool []*worker
	wg         sync.WaitGroup
}

func (c *kafkaConsumer) Close() {
	c.consumer.Close()
}

func NewConsumer(bootstrapServers string, groupId string, handler handler.Handler) (Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err := consumer.Subscribe(handler.Topic(), nil); err != nil {
		return nil, err
	}

	// wait for assignment
	for {

		ev := consumer.Poll(10000)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			consumer.Assign(e.Partitions)
			log.Printf("Assigned partitions: %v", e.Partitions)
			break
		case kafka.RevokedPartitions:
			consumer.Unassign()
			log.Printf("Revoked partitions: %v", e.Partitions)
			break
		case kafka.Error:
			log.Printf("Error: %v", e)
			break
		default:
			log.Printf("Ignored event: %s", e)
			break
		}

		break
	}

	partitions, err := consumer.Assignment()
	if err != nil {
		return nil, err
	}

	workerPool := make([]*worker, len(partitions))
	for i := range workerPool {
		workerPool[i] = newWorker(i, handler)
	}

	return &kafkaConsumer{
		consumer:   consumer,
		handler:    handler,
		workerPool: workerPool,
	}, nil
}

func (c *kafkaConsumer) Run() error {
	// Initialize the worker pool for each partition
	partitions, err := c.consumer.Assignment()
	if err != nil {
		return err
	}

	c.workerPool = make([]*worker, len(partitions))
	for i := range c.workerPool {
		c.workerPool[i] = newWorker(i, c.handler)
	}

	// Start the worker goroutines
	for _, w := range c.workerPool {
		w.start(&c.wg)
	}

	for {
		ev := c.consumer.Poll(0)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			if c.handler != nil {
				// Find the worker for this message's partition
				partition := e.TopicPartition.Partition
				worker := c.workerPool[partition]

				// Send the message to the worker's message queue
				worker.messages <- e
			}
		case kafka.PartitionEOF:
			continue
		case kafka.Error:
			log.Printf("Consumer error: %v (%v)", e, e.Code())
			if e.Code() == kafka.ErrAllBrokersDown {
				c.shutdownWorkers()
			}
		case kafka.AssignedPartitions:
			// Update the worker pool for the new set of partitions
			c.handleRebalance(e.Partitions)
		case kafka.RevokedPartitions:
			// Stop all worker goroutines for the old set of partitions
			c.stopWorkers()
		}
	}
}

func (c *kafkaConsumer) handleRebalance(partitions []kafka.TopicPartition) {
	// Stop all worker goroutines for the old set of partitions
	c.stopWorkers()

	// Initialize the worker pool for the new set of partitions
	c.workerPool = make([]*worker, len(partitions))
	for i := range c.workerPool {
		c.workerPool[i] = newWorker(i, c.handler)
	}

	// Start the worker goroutines
	for _, w := range c.workerPool {
		w.start(&c.wg)
	}
}

func (c *kafkaConsumer) stopWorkers() {
	for _, w := range c.workerPool {
		w.stop()
	}
	for _, w := range c.workerPool {
		<-w.done
	}
	c.wg.Wait()
}

type worker struct {
	id       int
	handler  handler.Handler
	messages chan *kafka.Message
	done     chan bool
}

func (w *worker) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range w.messages {
			if err := w.handler.HandleMessage(msg); err != nil {
				log.Printf("failed to handle message: %v", err)
			}
		}
		w.done <- true
	}()
}

func (w *worker) stop() {
	close(w.messages)
}

func newWorker(id int, handler handler.Handler) *worker {
	return &worker{
		id:       id,
		handler:  handler,
		messages: make(chan *kafka.Message),
		done:     make(chan bool),
	}
}

func (c *kafkaConsumer) shutdownWorkers() {
	log.Println("Shutting down workers...")
	for _, worker := range c.workerPool {
		worker.stop()
	}
	log.Println("Workers shut down.")
}
