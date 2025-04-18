package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
)

type Consumer struct {
	MsgChan   chan *sarama.ConsumerMessage
	Wg        *sync.WaitGroup
	msgRouter router.MessageRouter
	cfg       *ConsumerConfig
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer setup...")
	for i := 0; i < c.cfg.WorkerPoolSize; i++ {
		c.Wg.Add(1)
		go c.worker()
	}
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer cleanup...")
	c.Wg.Wait()
	return nil
}

// ConsumeClaim processes messages from Kafka
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Println("Incomming message")
		c.MsgChan <- msg
		sess.MarkMessage(msg, "gmail-service")
	}
	return nil
}
func (c *Consumer) worker() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		c.Wg.Done()
	}()
	for msg := range c.MsgChan {
		c.msgRouter.Route(context.Background(), msg)
	}
}
func NewConsumer(msgRouter router.MessageRouter, cfg *ConsumerConfig) *Consumer {
	return &Consumer{
		MsgChan:   make(chan *sarama.ConsumerMessage, cfg.WorkerPoolSize),
		Wg:        &sync.WaitGroup{},
		msgRouter: msgRouter,
		cfg:       cfg,
	}
}
