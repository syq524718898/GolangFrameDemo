package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var kafkaclusaddr = []string{"192.168.0.102:9095","192.168.0.102:9093","192.168.0.102:9094"}

func SyncProduce() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{}
	msg.Topic = "test_topic"
	// 元数据
	msg.Metadata = ""
	content := "this is a sync message"
	msg.Value = sarama.StringEncoder(content)

	client, err := sarama.NewSyncProducer(kafkaclusaddr, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}

	fmt.Printf("send msg partition/offset: %d/%d, value is: %s",partition,offset,content)

}

func AsyncProduce() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewAsyncProducer(kafkaclusaddr, config)
	if err != nil {
		fmt.Println("error is:", err.Error())
		return
	}
	defer client.AsyncClose()

	go func(p sarama.AsyncProducer) {
		for {
			select {
			case msg := <-p.Successes():
				value,_ := msg.Value.Encode()
				fmt.Printf("send msg partition/offset: %d/%d, value is: %s",msg.Partition,msg.Offset,string(value))
				return
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
				return
			}
		}
	}(client)

	msg := &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.ByteEncoder("this is a async message hello5"),
	}
	client.Input() <- msg
	time.Sleep(time.Second * 1)
}

func main() {
	AsyncProduce()
	// fmt.Println()
	// SyncProduce()
}

