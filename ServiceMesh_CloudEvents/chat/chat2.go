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

func (s *Server) SayHello(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	//log.Printf("Recived message body from client %s , %s , %s , %s , %s ", message.Timestamp, message.Location, message.Sensortyp, message.SensorID, message.SensorData)
	Timestamp = message.Timestamp
	/*Location = message.Location
	Sensortyp = message.Sensortyp
	SensorID = message.SensorID
	SensorData = message.SensorData*/

	localmessage = Timestamp + "|" + Location + "|" + Sensortyp + "|" + string(SensorID) + "|" + SensorData
	go SafeToFile(Sensortyp, SensorData, Timestamp)
	//go initialize()

	return &MessageReply{Message: "Hello From the Server!"}, nil
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
	//Check ob Verbindung schonmal aufgebaut?
	b, err := ioutil.ReadFile("files/master2Config.txt")
	if err != nil {
		panic(err)
	}
	suchRegex := sid + ";"
	fmt.Println(suchRegex)
	isExist, err := regexp.Match(suchRegex, b)
	if err != nil {
		panic(err)
	}

	if isExist == false {
		fmt.Println("Speichern von Verbindungsdaten um für Health Checks")
		fmt.Println("Der Datensatz wird in ein File gespeichert")

		f, err := os.OpenFile("files/master2Config.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Wo ist meine ID!?")
		fmt.Println(sid)
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
		fmt.Println("Der Eintrag ist schon enthalten")

		file, err := os.Open("files/master2Config.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("Was wird gelesen")
			fmt.Println(line)
			lines = append(lines, line)
		}

		fmt.Println("Laenge von lines:", len(lines))
		fmt.Println("Vor der for")
		for i := 0; i < len(lines); i++ {
			fmt.Println("Lines:" + lines[i])
			if strings.Contains(lines[i], suchRegex) {
				//lines[i]= sid + ";" + ip + ";" + "Timestamp" + ";"

				lines[i] = sid + ";" + ip + ";" + tst + ";"
				break
			}
		}
		fmt.Println("Nach der for")

		f, err := os.OpenFile("files/master2Config.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, value := range lines {
			fmt.Println(value)
			//fmt.Println("Wir schreiben den Value:"+neuerValue)
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
