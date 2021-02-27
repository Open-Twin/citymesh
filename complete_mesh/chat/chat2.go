package chat

import (
	"bufio"
	"strings"
	//"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type Server struct {
}

var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string

var localmessage string

const (
	kafkaConn = "localhost:9092"
	topic     = "topic_test"
)

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	panic("implement me")
}

func (s *Server) DataFromSidecar(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	//log.Printf("Received message body from client %s , %s , %s , %s , %s ", message.IdService, message.Source, message.SpecVersion, message.Type, message.IdService, message.IpSidecar, message.IpSidecar, message.Timestamp, message.Data)
	log.Printf("Received message body from client: %s, %s", message.IpSidecar, message.Source)
	fmt.Println(message.Data)
	//log.Printf( "Received message body from client")
	message.IdSidecar = "sd123"
	message.IpSidecar = "123.123.123.123"
	//msg := message
	SafeToFile(message.IdSidecar, message.IpSidecar, message.Timestamp)
	return &MessageReply{Message: "Sidecar Proxy: Received the message"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, health2 *Health) (*MessageReply, error) {
	/*
		fmt.Println("Sending pings to all registered entries")

		file, err := os.Open("files/master2Config.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string
		var report string
		var health sidecar.Health



		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		for _, value := range lines {
			res := strings.Split(value, ";")
			report = report + res[0] + "; \n"
			fmt.Println(res)
			// get()
			// create client for GRPC Server
			var conn *grpc.ClientConn
			// Placeholder until Gateway
			conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
			//conn, err := grpc.Dial(target+":9000", grpc.WithInsecure())

			if erro != nil {
				log.Fatalf("no server connection could be established cause: %v", erro)

			}

			// defer runs after the functions finishes
			defer conn.Close()

			c := sidecar.NewChatServiceClient(conn)

			response, err := c.HealthCheck(context.Background(), &health)

			if response != nil {
				log.Printf("Response from Sidecar: %s ", response.Reply)
			}

			if err != nil {
				fmt.Println(err)
			}

			report = report + response.Reply

		}
		return &MessageReply{Reply: report}, nil
	*/

	return &MessageReply{Message: "Hello"}, nil
}

func (s *Server) GetMessage() string {
	return localmessage
}

func initialize() {
	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	msg := localmessage

	// publish without goroutene
	publish(msg, producer)

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

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
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

func SafeToFile(ip string, sid string, tst string) {
	//Checks if Sidecar ID already has an entry
	b, err := ioutil.ReadFile("files/master2Config.csv")
	if err != nil {
		panic(err)
	}
	suchRegex := sid + ";"
	isExist, err := regexp.Match(suchRegex, b)
	if err != nil {
		panic(err)
	}

	if isExist == false {

		f, err := os.OpenFile("files/master2Config.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}

		l, err := f.WriteString(sid + ";" + ip + ";" + "Timestamp" + ";" + "\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "bytes written successfully")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// entry already exists
		fmt.Println("Debug: ID already stored")
		file, err := os.Open("files/master2Config.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string

		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], suchRegex) {

				lines[i] = sid + ";" + ip + ";" + tst + ";"
				break
			}
		}

		f, err := os.OpenFile("files/master2Config.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, value := range lines {
			fmt.Fprintln(f, value) // print values to f, one per line
		}
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
