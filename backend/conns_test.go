package backend

import (
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	"github.com/zhwei820/gconv"
)

var conn *Conn

func init() {

	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0

	admin, err := sarama.NewClusterAdmin([]string{"10.171.22.153:30200"}, config)
	if err != nil {
		panic(err)
	}
	consumer, err := sarama.NewConsumer([]string{"10.171.22.153:30200"}, config)
	if err != nil {
		panic(err)
	}
	conn = &Conn{
		InstanceAdmin:    admin,
		InstanceConsumer: consumer,
	}

	gconv.SetExportExpand(true)
}

var (
	group = "inner_loan_group_test"
	topic = "inner_loan_order_topic"
)

func TestConn_GetConsumerGroupMessageLag(t *testing.T) {
	messageLagResult, err := conn.GetConsumerGroupMessageLag(topic, group)
	fmt.Println("messageLagResult, err", gconv.Export(messageLagResult), err)
}

// messageLagResult, err [
// 	{
// 		"LatestOffset": 111685804,
// 		"ConsumerGroupOffset": 111685623,
// 		"Partition": 0,
// 		"Topic": "inner_loan_order_topic"
// 	}
// ]

// messageLagResult, err [
// 	{
// 		"LatestOffset": 111771570,
// 		"ConsumerGroupOffset": 111771463,
// 		"Partition": 0,
// 		"Topic": "inner_loan_order_topic"
// 	}
// ]

func TestConn_TailMessage(t *testing.T) {
	messageResult, err := conn.TailMessage(topic, 500)
	fmt.Println("messageResult, err", gconv.Export(messageResult), err)
}
