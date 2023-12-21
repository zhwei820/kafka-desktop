package backend

import (
	"kafka_desktop/backend/conf"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/zhwei820/gconv"
)

type Conn struct {
	c             *conf.KafkaConfig
	InstanceAdmin sarama.ClusterAdmin
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
