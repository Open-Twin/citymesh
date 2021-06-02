package Consumer

import (
	_ "context"
	"encoding/json"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/broker"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"gotest.tools/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"testing"
	"time"
)

type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

var (
	addrs   = []string{"localhost:9092"}
	brokers = []string{"0.0.0.0:9092"}
	topic   = "topic_test"
)

const (
	kafkaConn = "localhost:9092"
)

/*func Test_consume(t *testing.T) {
	type args struct {
		topics []string
		master sarama.Consumer
	}
	tests := []struct {
		name  string
		args  args
		want  chan *sarama.ConsumerMessage
		want1 chan *sarama.ConsumerError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := main.consume(tt.args.topics, tt.args.master)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("consume() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("consume() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}*/

func TestConsume(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, config)
	require.NoError(t, err)

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder([]byte("foobar")),
	})
	require.NoError(t, err)
	t.Logf("Sent message to partition %d with offset %d", partition, offset)

	//ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	master, err := sarama.NewConsumer(brokers, config)

	topics, _ := master.Topics()

	consumer, errors := consume(topics, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				receivedMessage := &sidecar.CloudEvent{}
				err := proto.Unmarshal(msg.Value, receivedMessage)
				if err != nil {
					log.Fatalln("Failed to unmarshal message:", err)
				}

				//log.Printf("Message received: %s", receivedMessage.Data)
				//fmt.Println("Message received: ", receivedMessage.Data)
				fmt.Println("Message received: ", receivedMessage)
				//fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case consumerError := <-errors:
				msgCount++
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	log.Println("Sarama consumer up and running!")

	time.Sleep(1 * time.Second)

	wg.Wait()
	//close()
}

func TestConsumeTwice(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, config)
	require.NoError(t, err)

	data1, data2 := "foobar1", "foobar2"

	for _, data := range []string{data1, data2} {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("foobar"),
			Value: sarama.StringEncoder(data),
		})
		require.NoError(t, err)
		t.Logf("Sent message to partition %d with offset %d", partition, offset)
	}

	producertwo, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	cloudmessage := Apiclient2()

	//fmt.Println(cloudmessage.Data)

	publish(cloudmessage, producertwo)

	//ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topics, _ := master.Topics()

	consumer, errors := consume(topics, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				receivedMessage := &sidecar.CloudEvent{}
				err := proto.Unmarshal(msg.Value, receivedMessage)
				if err != nil {
					log.Fatalln("Failed to unmarshal message:", err)
				}

				//log.Printf("Message received: %s", receivedMessage.Data)
				//fmt.Println("Message received: ", receivedMessage.Data)
				fmt.Println("Message received: ", receivedMessage)
				//fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case consumerError := <-errors:
				msgCount++
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	messageReceived := make(chan []byte)
	/*cons := &
	consumer := &Consumer{
		ready: make(chan bool),
		handle: func(data []byte) error {
			messageReceived <- data
			fmt.Printf("Received message: %s\n", data)
			return nil
		},
	}*/

	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				receivedMessage := &sidecar.CloudEvent{}
				err := proto.Unmarshal(msg.Value, receivedMessage)
				if err != nil {
					log.Fatalln("Failed to unmarshal message:", err)
				}

				//log.Printf("Message received: %s", receivedMessage.Data)
				//fmt.Println("Message received: ", receivedMessage.Data)
				fmt.Println("Message received: ", receivedMessage)
				//fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case consumerError := <-errors:
				msgCount++
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	//close := main.consume(ctx, &wg, consumer)

	//<-consumer.ready
	log.Println("Sarama consumer up and running!")

	for i := 0; i < 2; i++ {
		data := <-messageReceived
		switch i {
		case 0:
			assert.Equal(t, data1, string(data))
		case 1:
			assert.Equal(t, data2, string(data))
		}
	}

	//cancel()
	wg.Wait()
	//close()
}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = 10000000

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}
		partitions, _ := master.Partitions(topic)
		// this only consumes partition no 1, you would probably want to consume all partitions
		consumer, err := master.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			fmt.Printf("Topic %v Partitions: %v", topic, partitions)
			panic(err)
		}
		fmt.Println(" Start consuming topic ", topic)
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					consumers <- msg
					fmt.Println("Got message on topic ", topic, msg.Value)
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}

func publish(message broker.CloudEvent, producer sarama.SyncProducer) {
	// publish sync

	messageToSend := &message
	messageToSendBytes, err := proto.Marshal(messageToSend)

	//fmt.Println(messageToSendBytes)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageToSendBytes),
	}

	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}

func Apiclient() (cloudeventmessage broker.CloudEvent) {

	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	ampelJson := string(body)

	value := gjson.Get(ampelJson, "#.Stand")
	println("Length: ", len(value.Array()))
	println(value.String())

	messages := make([]*dataFormat.Message, 0, 0)
	//0,0 checken

	datenres := gjson.Get(ampelJson, "#.Stand")
	datenres.ForEach(func(key, value gjson.Result) bool {
		stand := value.String()
		warnstufen := make([]*dataFormat.Warnstufen, 0, 0)

		path := "#(Stand==" + stand + ").Warnstufen"
		result := gjson.Get(ampelJson, path)
		result.ForEach(func(key, value gjson.Result) bool {
			var warning Warning
			if err := json.Unmarshal([]byte(value.Raw), &warning); err != nil {
				panic(err)
			}
			warnstufen = append(warnstufen,
				&dataFormat.Warnstufen{
					Region:    warning.Region,
					GKZ:       warning.Gkz,
					Name:      warning.Name,
					Warnstufe: warning.Warnstufe,
				})
			return true
		})
		messages = append(messages,
			&dataFormat.Message{Stand: stand,
				Warnstufen: warnstufen,
			})

		return true
	})

	marshalMessages, err := ptypes.MarshalAny(&dataFormat.Messages{Message: messages})
	if err != nil {
		panic(err)
	}

	data := broker.CloudEvent_ProtoData{ProtoData: marshalMessages}

	cloudeventmessage = broker.CloudEvent{
		IdService:   "",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &data,
		IpService:   "192.168.0.10",
	}

	//fmt.Print(cloudeventmessage.Data)

	//fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}

func Apiclient2() (cloudeventmessage broker.CloudEvent) {

	ampel, erro := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	if erro != nil {
		fmt.Print(erro.Error())
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(ampel.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	ampelJson := string(body)

	value := gjson.Get(ampelJson, "#.Stand")
	println("Length: ", len(value.Array()))
	println(value.String())

	messages := make([]*dataFormat.Message, 0, 0)
	//0,0 checken

	datenres := gjson.Get(ampelJson, "#.Stand")
	datenres.ForEach(func(key, value gjson.Result) bool {
		stand := value.String()
		warnstufen := make([]*dataFormat.Warnstufen, 0, 0)

		path := "#(Stand==" + stand + ").Warnstufen"
		result := gjson.Get(ampelJson, path)
		result.ForEach(func(key, value gjson.Result) bool {
			var warning Warning
			if err := json.Unmarshal([]byte(value.Raw), &warning); err != nil {
				panic(err)
			}
			warnstufen = append(warnstufen,
				&dataFormat.Warnstufen{
					Region:    warning.Region,
					GKZ:       warning.Gkz,
					Name:      warning.Name,
					Warnstufe: warning.Warnstufe,
				})
			return true
		})
		messages = append(messages,
			&dataFormat.Message{Stand: stand,
				Warnstufen: warnstufen,
			})
		//jsonmessage,err := ProtobufToJSON(messages)

		return true
	})

	print("Message", messages)

	/*marshalMessages, err := ptypes.MarshalAny(&dataFormat.Messages{Message: messages})
	if err != nil {
		panic(err)
	}*/

	//data := sidecar.CloudEvent_ProtoData{ProtoData: marshalMessages}

	//jsondata,err := ProtobufToJSON(messages)

	text := ""
	for _, element := range messages {
		text = text + string(proto.MarshalTextString(element))
		//result, _ := ProtobufToJSON(element)
		//text = text + result

	}

	finishedMessage := broker.CloudEvent_TextData{TextData: text}

	cloudeventmessage = broker.CloudEvent{
		IdService:   "",
		Source:      "corona-ampel",
		SpecVersion: "1.0",
		Type:        "json",
		Attributes:  nil,
		Data:        &finishedMessage,
		IdSidecar:   "01",
		IpService:   "192.168.0.10",
		IpSidecar:   "192.168.0.11",
	}

	//fmt.Println(cloudeventmessage.Data)

	/*response, err := c.DataFromService(context.Background(), &cloudeventmessage)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())*/

	return cloudeventmessage

}
