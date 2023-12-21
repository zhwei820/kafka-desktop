package main

import (
	"kafka_desktop/backend"
	"kafka_desktop/backend/conf"
	"log"

	"github.com/IBM/sarama"
)

func (app *App) AppConf(key string) *conf.Config {
	return app.conf
}

func (app *App) DBConf(key string) map[string]*conf.KafkaConfig {
	return app.conf.Instance
}

func (app *App) SetAPPConf(c *conf.Config) {
	app.conf = c
}
func (app *App) SetDBConf(c map[string]*conf.KafkaConfig) {
	app.conf.Instance = c
	app.conf.SaveAll()
}

func (app *App) OpenConn(key string) error {
	_, err := backend.OpenConn(app.conf, key)
	return err
}

func (app *App) ConnsOpened() map[string]*backend.Conn {
	return backend.ConnsOpened
}

func (app *App) ListConsumerGroups(connName string) (map[string]string, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.ListConsumerGroups()
}
func (app *App) DescribeConsumerGroups(connName string, groups []string) ([]*sarama.GroupDescription, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.DescribeConsumerGroups(groups)
}

func (app *App) ListTopics(connName string) (map[string]sarama.TopicDetail, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.ListTopics()
}
func (app *App) DescribeTopics(connName string, groups []string) ([]*sarama.TopicMetadata, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.DescribeTopics(groups)
}
func (app *App) TailMessage(connName string, topic string) (map[int32][]*sarama.ConsumerMessage, error) {
	consumer := backend.ConnsOpened[connName].InstanceConsumer

	// Fetch the latest offset for each partition
	// "高水位标记"（high water marks）是指 Kafka 中每个分区最新的已提交消息的偏移量（offset）。这个偏移量表示了消费者（consumer）在分区中已经读取的消息位置。
	offsets := consumer.HighWaterMarks()

	// Define a map to store fetched messages
	messages := make(map[int32][]*sarama.ConsumerMessage)

	// Fetch the latest 500 messages for each partition
	for partition, offset := range offsets[topic] {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset-500)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %v", err)
		}
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalf("Error closing partition consumer: %v", err)
			}
		}()

		// Iterate through fetched messages
		for message := range partitionConsumer.Messages() {
			messages[partition] = append(messages[partition], message)
		}
	}

	return messages, nil
}
