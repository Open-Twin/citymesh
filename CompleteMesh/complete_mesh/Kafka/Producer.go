package main

import (
	"github.com/Open-Twin/CityMesh-ProtoTypes/Jakob_GRPC/chat"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
)

type Datastructtwo []struct {
	stand      string `json:"Stand"`
	Warnstufen []struct {
		Gkz       string `json:"GKZ"`
		Name      string `json:"Name"`
		Region    string `json:"Region"`
		Warnstufe string `json:"Warnstufe"`
	} `json:"Warnstufen"`
}

type dataMessagetwo struct {
	message []Datastructtwo
}

const (
	kafkaConn = "localhost:9092"
	topic     = "topic_test"
)

func main() {

}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func publish(message chat.Message, producer sarama.SyncProducer) {
	// publish sync

	messageToSend := &message
	messageToSendBytes, err := proto.Marshal(messageToSend)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageToSendBytes),
	}

	producer.SendMessage(msg)
	log.Printf("Message sent")

	if err != nil {
		log.Fatalln("Failed", err)
	}

	/*msg := &sarama.ProducerMessage {
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}*/
	/*p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}*/

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	//fmt.Println("Partition: ", p)
	//fmt.Println("Offset: ", o)
}
