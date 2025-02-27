package initialize

import (
	"github.com/IBM/sarama"
	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
)

func InitializeKafkaProducer(config *config.KafkaConfig) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = config.MaxRetry
	kafkaConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
