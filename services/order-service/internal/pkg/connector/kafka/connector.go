package kafka

import (
	"log"

	"order-service/config"

	"github.com/IBM/sarama"
)

func MustSyncProducer() sarama.SyncProducer {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer(config.Instance().Kafka.Brokers, saramaConfig)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}

	return producer
}
