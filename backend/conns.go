package backend

import (
	"encoding/json"
	"fmt"
	"kafka_desktop/backend/conf"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/zhwei820/gconv"
)

type Conn struct {
	c                *conf.KafkaConfig
	InstanceAdmin    sarama.ClusterAdmin
	InstanceConsumer sarama.Consumer

	Name string
}

var ConnsOpened = map[string]*Conn{}

func OpenConn(c *conf.Config, key string) (res *Conn, err error) {
	log.Debug().Str("==>c", gconv.Export(c)).Send()
	log.Debug().Str("==>ConnsOpened", gconv.Export(ConnsOpened)).Send()

	if _, ok := ConnsOpened[c.Instance[key].Name]; ok {
		return ConnsOpened[c.Instance[key].Name], nil
	}
	res = &Conn{
		c:    c.Instance[key],
		Name: key,
	}
	if c.Instance[key].Type == "kafka" {
		config := sarama.NewConfig()
		config.Net.MaxOpenRequests = 255

		config.Consumer.Offsets.Initial = sarama.OffsetOldest

		config.Consumer.MaxWaitTime = 500 * time.Millisecond
		config.Consumer.MaxProcessingTime = 2000 * time.Millisecond
		config.Consumer.Group.Session.Timeout = 15 * time.Second
		config.Consumer.Group.Heartbeat.Interval = 5 * time.Second

		config.Version, _ = sarama.ParseKafkaVersion("2.2.0")

		res.InstanceAdmin, err = sarama.NewClusterAdmin(strings.Split(c.Instance[key].Host, ";"), config)
		res.InstanceConsumer, err = sarama.NewConsumer(strings.Split(c.Instance[key].Host, ";"), config)
	}
	if err == nil {
		ConnsOpened[c.Instance[key].Name] = res
	} else {
		log.Error().Err(err).Str("key", key).Send()
	}
	return res, err
}

type EventResultStatus uint8

const (
	EventProcessing EventResultStatus = iota
	EventSuccess
	EventFailed
)

const (
	FailedCodeOk = 0
)

// MessageEvent common event
type MessageEvent struct {
	Data    interface{}             `json:"data" gorm:"-"`
	Message *sarama.ConsumerMessage `json:"message" gorm:"-"`
}

type MessageLag struct {
	LatestOffset        int64
	ConsumerGroupOffset int64
	Partition           int32
	Topic               string
}

func (conn *Conn) GetConsumerGroupMessageLag(topic, group string) ([]*MessageLag, error) {
	consumer := conn.InstanceConsumer
	admin := conn.InstanceAdmin

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error fetching partitions: ", err)
	}

	response, err := admin.ListConsumerGroupOffsets(group, map[string][]int32{
		topic: partitions,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error ListConsumerGroupOffsets")
		return nil, err
	}

	res := []*MessageLag{}
	for _, partition := range partitions {
		block := response.GetBlock(topic, partition)
		if block == nil {
			continue
		}

		latestOffset, err := consumer.Client().GetOffset(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Error().Err(err).Msg("Error GetOffset")
		}
		offsetInGroupConsumer := block.Offset
		if offsetInGroupConsumer > 0 {
			res = append(res, &MessageLag{
				LatestOffset:        latestOffset,
				ConsumerGroupOffset: offsetInGroupConsumer,
				Partition:           partition,
				Topic:               topic,
			})
		}
	}
	return res, nil
}

func (conn *Conn) TailMessage(topic string, lookBackNum int64) (map[int32][]*MessageEvent, error) {
	consumer := conn.InstanceConsumer

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error fetching partitions: ", err)
		log.Error().Err(err).Msg("Error get Partitions")
		return nil, err
	}

	offsets := make(map[int32]int64)
	for _, partition := range partitions {
		latestOffset, err := consumer.Client().GetOffset(topic, partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Error fetching offset for partition : ", partition, err)
			log.Error().Err(err).Msg("Error GetOffset")
			continue
		}
		offsets[partition] = latestOffset
	}

	messages := make(map[int32][]*MessageEvent)

	for partition, offset := range offsets {
		partition, offset := partition, offset
		go func() {
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset-lookBackNum)
			if err != nil {
				fmt.Println("Error creating partition consumer: ", err)
				log.Error().Err(err).Msg("Error ConsumePartition")
				return
			}
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					fmt.Println("Error closing partition consumer: ", err)
					log.Error().Err(err).Msg("Error closing partition")
				}
			}()
			for message := range partitionConsumer.Messages() {
				messageEvent := MessageEvent{}
				messageEvent.Data = map[string]interface{}{}
				if err := json.Unmarshal(message.Value, &messageEvent); err != nil {
					log.Error().Err(err).Msg("Error Unmarshal")
				}
				if gconv.String(messageEvent.Data) == "{}" { //empty, try without data wrap
					if err := json.Unmarshal(message.Value, &messageEvent.Data); err != nil {
						log.Error().Err(err).Msg("Error Unmarshal")
					}
				}
				messageEvent.Message = message
				messages[partition] = append(messages[partition], &messageEvent)
			}
		}()
	}
	time.Sleep(2 * time.Second)
	return messages, err
}
