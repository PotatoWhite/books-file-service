package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/potatowhite/books/file-service/handler"
	"log"
	"sync"
)

type worker struct {
	id       int
	handler  handler.Handler
	messages chan *kafka.Message
	done     chan bool
}

func (w *worker) start(wg *sync.WaitGroup) {
	go func() {
		for msg := range w.messages {
			wg.Wait()
			if err := w.handler.HandleMessage(msg); err != nil {
				log.Printf("failed to handle message: %v", err)
			}
		}
		w.done <- true
	}()

	logger.Printf("worker %d started", w.id)
}

func (w *worker) stop() {
	// close if not closed
	select {
	case <-w.done:
		return
	default:
		close(w.messages)
	}

	logger.Printf("worker %d stopped", w.id)
}

func newWorker(id int, handler handler.Handler) *worker {

	return &worker{
		id:       id,
		handler:  handler,
		messages: make(chan *kafka.Message),
		done:     make(chan bool),
	}
}
