package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	KafkaCluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"time"
)

var kafkaclusaddr = []string{"192.168.0.102:9095","192.168.0.102:9093","192.168.0.102:9094"}

func main(){
	KafkaConsumerCluster("")
}

type Canal struct {
	Data []Data		`json:"data"`
	Database string `json:"database"`
}

type Data struct {
	Id string	`json:"id"`
	Name string	`json:"name"`
	Age string	`json:"age"`
}

func KafkaConsumerCluster(consumerId string)  {

	topics := []string{"test_topic"}
	//topics := []string{"canal_topic"}
	config := KafkaCluster.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Consumer.Offsets.AutoCommit.Enable=true
	//config.Consumer.Offsets.AutoCommit.Interval=1*time.Second
	config.Consumer.Offsets.CommitInterval=1*time.Second
	// 关闭自动提交
	//config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial=sarama.OffsetOldest
	config.Group.Return.Notifications = true

	//第二个参数是groupId
	consumer, err := KafkaCluster.NewConsumer(kafkaclusaddr, "consumer-group1", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// 接收错误
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// 打印一些rebalance的信息
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// 消费消息
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s : %s/%d/%d\t%s\t%s\n", consumerId,msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				data := &Canal{}
				json.Unmarshal(msg.Value,data)
				fmt.Println(data)

				consumer.MarkOffset(msg, "")   // 提交offset
			}
		case <-signals:
			return
		}
	}
}
