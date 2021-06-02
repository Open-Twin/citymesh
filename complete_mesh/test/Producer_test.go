package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/broker"
	"github.com/Open-Twin/citymesh/complete_mesh/dataFormat"
	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

/*func TestApiclient(t *testing.T) {
	tests := []struct {
		name                  string
		wantCloudeventmessage broker.CloudEvent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCloudeventmessage := main.Apiclient(); !reflect.DeepEqual(gotCloudeventmessage, tt.wantCloudeventmessage) {
				t.Errorf("Apiclient() = %v, want %v", gotCloudeventmessage, tt.wantCloudeventmessage)
			}
		})
	}
}

func TestApiclient2(t *testing.T) {
	tests := []struct {
		name                  string
		wantCloudeventmessage broker.CloudEvent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCloudeventmessage := main.Apiclient2(); !reflect.DeepEqual(gotCloudeventmessage, tt.wantCloudeventmessage) {
				t.Errorf("Apiclient2() = %v, want %v", gotCloudeventmessage, tt.wantCloudeventmessage)
			}
		})
	}
}

func Test_initProducer(t *testing.T) {
	tests := []struct {
		name    string
		want    sarama.SyncProducer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := main.initProducer()
			if (err != nil) != tt.wantErr {
				t.Errorf("initProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initProducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_publish(t *testing.T) {
	type args struct {
		message  broker.CloudEvent
		producer sarama.SyncProducer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}*/

type Warning struct {
	Gkz       string `json:"GKZ"`
	Name      string `json:"Name"`
	Region    string `json:"Region"`
	Warnstufe string `json:"Warnstufe"`
}

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

func Test_Producer(t *testing.T) {
	/*tests := []struct {
		name    string
		want    sarama.SyncProducer
		wantErr bool
	}{
		// TODO: Add test cases.
	}*/

	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	cloudmessage := Apiclient2()

	//fmt.Println(cloudmessage.Data)

	publish(cloudmessage, producer)

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

var (
	addrs = []string{"localhost:9092"}
	topic = "my-topic"
)

func TestSendMessage(t *testing.T) {
	//given
	kafka := testcontainers.NewLocalDockerCompose(
		[]string{configs.RootDir() + "kafka/test-docker-compose.yml"},
		strings.ToLower(uuid.New().String()),
	)
	kafka.WithCommand([]string{"up", "-d"}).Invoke()
	time.Sleep(5 * time.Second)
	defer destroyKafka(kafka)
	//when
	SendMessage("test-topic", "test-key", "test-message")
	//then
	result := read()
	require.Equal(t, result["test-key"], "test-message")
}
func read() map[string]string {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	})
	m, _ := r.ReadMessage(context.Background())
	defer r.Close()

	result := make(map[string]string)
	result[string(m.Key)] = string(m.Value)
	return result
}
func destroyKafka(compose *testcontainers.LocalDockerCompose) {
	compose.Down()
	time.Sleep(1 * time.Second)
}

// Test mit meinem Producer und Consumer
