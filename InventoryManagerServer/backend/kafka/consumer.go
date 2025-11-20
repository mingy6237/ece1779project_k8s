package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

var Consumer sarama.Consumer

func InitConsumer(brokers []string) error {
	var err error
	Consumer, err = sarama.NewConsumer(brokers, nil)
	if err != nil {
		return err
	}
	log.Println("Kafka broadcast consumer initialized.")
	return nil
}

func StartConsumer(ctx context.Context, topic string, handler func(*sarama.ConsumerMessage)) error {
	partitions, err := Consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, p := range partitions {
		pc, err := Consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for {
				select {
				case msg := <-pc.Messages():
					if msg != nil {
						handler(msg)
					}
				case err := <-pc.Errors():
					log.Printf("Kafka partition %d error: %v", p, err)
				case <-ctx.Done():
					return
				}
			}
		}(pc)
	}

	log.Printf("Broadcast consumer started on topic %s (all partitions)", topic)
	return nil
}

func CloseConsumer() error {
	if Consumer != nil {
		return Consumer.Close()
	}
	return nil
}
