package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/config"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
)

// Consumer implements sarama.ConsumerGroupHandler
type Consumer struct {
	MsgChan   chan *sarama.ConsumerMessage
	Wg        *sync.WaitGroup
	msgRouter router.MessageRouter
	cfg       *config.ConsumerConfig
}

// Setup runs when the consumer group session starts
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer setup...")
	for i := 0; i < c.cfg.WorkerPoolSize; i++ {
		c.Wg.Add(1)
		go c.worker()
	}
	return nil
}

// Cleanup runs when the consumer group session ends
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer cleanup...")
	c.Wg.Wait()
	return nil
}

// ConsumeClaim processes messages from Kafka
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
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
		ctx := context.Background()
		err := c.msgRouter.Route(ctx, msg)
		fmt.Println(err)
	}
}
