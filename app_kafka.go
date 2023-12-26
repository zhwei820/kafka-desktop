package main

import (
	"fmt"
	"kafka_desktop/backend"
	"kafka_desktop/backend/conf"

	"github.com/IBM/sarama"
	"github.com/zhwei820/gconv"
	"gitlab.matrixport.com/common/generic_tool/gmap"
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

func (app *App) ListConsumerGroups(connName string) (map[string]string, error) { // "asset_liability_group": "consumer", // groups
	return backend.ConnsOpened[connName].InstanceAdmin.ListConsumerGroups()
}

func (app *App) DescribeConsumerGroups(connName string, groups []string) (map[string][]int32, error) {
	result, err := backend.ConnsOpened[connName].InstanceAdmin.DescribeConsumerGroups(groups)
	if err != nil {
		return nil, err
	}
	res, _ := result[0].Members[gmap.Keys(result[0].Members)[0]].GetMemberAssignment()
	fmt.Println("GetMemberAssignment", gconv.Export(res))

	return res.Topics, err
}

func (app *App) ListTopics(connName string) (map[string]sarama.TopicDetail, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.ListTopics()
}
func (app *App) DescribeTopics(connName string, topics []string) ([]*sarama.TopicMetadata, error) {
	return backend.ConnsOpened[connName].InstanceAdmin.DescribeTopics(topics)
}

func (app *App) TailMessage(connName, topic string, num int64) (map[int32][]*backend.MessageEvent, error) {
	return backend.ConnsOpened[connName].TailMessage(topic, num)
}

func (app *App) GetConsumerGroupMessageLag(connName string, topic, group string) ([]*backend.MessageLag, error) {
	return backend.ConnsOpened[connName].GetConsumerGroupMessageLag(topic, group)
}
