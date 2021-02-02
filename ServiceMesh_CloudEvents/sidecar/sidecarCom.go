package sidecar

import (
	"bufio"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Server struct {
}

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	panic("implement me")
}

var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string

var localmessage string

var msg *CloudEvent

func (s *Server) DataFromService(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	log.Printf("Recived message body from client %s , %s , %s , %s , %s ", message.Type, message.Timestamp, message.IdService, message.IdSidecar, message.IpSidecar, message.IpService, message.Source)
	Timestamp = message.Timestamp
	Type := message.Type
	IDService := message.IdService
	IDSidecar := message.IdSidecar
	IPSidecar := message.IpSidecar
	IPService := message.IpService
	Source := message.Source

	localmessage = Type + "|" +  Timestamp + "|" + IDService + "|" + IDSidecar + "|" + IPService + "|" + IPSidecar + "|" + Source
	msg = message
	SafeToFile(Sensortyp, SensorData, Timestamp)
	go client(msg)
	return &MessageReply{
		Message: "Hello From the Sidecar Proxy!",
	}, nil
}

func (s *Server) HealthCheck(ctx context.Context, schlecht *ProtoSchlecht) (*MessageReply, error) {
	/*fmt.Println("Starte Health Check")
	log.Printf("Starte Health Check")

	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run() // blocks until finished
	if err != nil {
		panic(err)
	}
	//stats := pinger.Statistics() // get send/receive/rtt stats*/

	return &MessageReply{Message: "HealthCheck started"}, nil
}

func (s *Server) GetMessage() string {

	return localmessage
}

func SafeToFile(ip string, sid string, tst string) {
	//Check ob Verbindung schonmal aufgebaut?
	b, err := ioutil.ReadFile("files/sidecarConfig.txt")
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

		f, err := os.OpenFile("files/sidecarConfig.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

		file, err := os.Open("files/sidecarConfig.txt")
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

		f, err := os.OpenFile("files/sidecarConfig.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
