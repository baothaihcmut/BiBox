package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	producer sarama.SyncProducer
}

func (k *KafkaService) PublishMessage(ctx context.Context, topic string, value any, headers map[string]string) (int32, int64, error) {
	//value
	data, err := json.Marshal(value)
	if err != nil {
		return 0, 0, err
	}
	//headers
	header := make([]sarama.RecordHeader, 0)
	for key, val := range headers {
		header = append(header,
			sarama.RecordHeader{
				Key:   []byte(key),
				Value: []byte(val),
			},
		)
	}
	message := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(data),
		Headers:   header,
		Timestamp: time.Now(),
	}
	partition, offset, err := k.producer.SendMessage(message)
	if err != nil {
		return 0, 0, err
	}
	return partition, offset, nil
}

func NewKafkaService(
	producer sarama.SyncProducer,
) QueueService {
	return &KafkaService{
		producer: producer,
	}
}
