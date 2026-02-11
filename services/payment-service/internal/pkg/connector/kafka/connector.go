package kafka

import (
	"log"
	"time"

	"payment-service/config"

	"github.com/IBM/sarama"
)

func MustConsumerGroup() sarama.ConsumerGroup {
	cfg := sarama.NewConfig()

	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	cfg.Consumer.Group.Rebalance.Retry.Max = 4                          // количество повторов при ребалансе
	cfg.Consumer.Group.Rebalance.Retry.Backoff = 500 * time.Millisecond // интервал между попытками

	cfg.Consumer.Group.Session.Timeout = 30 * time.Second   // таймаут сессии (heartbeat)
	cfg.Consumer.Group.Heartbeat.Interval = 5 * time.Second // heartbeat для координации

	cfg.Consumer.Offsets.AutoCommit.Enable = false
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	cfg.Consumer.Fetch.Min = 1
	cfg.Consumer.Fetch.Default = 1024 * 1024  // размер батча для чтения
	cfg.Consumer.Fetch.Max = 10 * 1024 * 1024 // максимум за один fetch
	cfg.Consumer.MaxProcessingTime = 500 * time.Millisecond

	group, err := sarama.NewConsumerGroup(config.Instance().Kafka.Brokers, config.Instance().Kafka.ConsumerGroup, cfg)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}

	return group
}

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
