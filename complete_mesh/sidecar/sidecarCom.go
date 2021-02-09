package sidecar

import (
	"bufio"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type Server struct {
}

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	//panic("implement me")
}

var Timestamp string
var Location string
var Sensortyp string
var SensorID int32
var SensorData string

var localmessage string

var msg *CloudEvent

func (s *Server) DataFromService(ctx context.Context, message *CloudEvent) (*MessageReply, error) {
	log.Printf("Received message body from client %s , %s , %s , %s , %s ", message.IdService, message.Source, message.SpecVersion, message.Type, message.IdService, message.IpSidecar, message.IpSidecar, message.Timestamp)
	//log.Printf( "Received message body from client")
	message.IdSidecar = "sd123"
	message.IpSidecar = "123.123.123.123"
	msg = message
	SafeToFile(message.IdService, message.IpService, message.Timestamp)
	go client(msg)
	return &MessageReply{Message: "Sidecar Proxy: Received the message"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, health *Health) (*MessageReply, error) {
	fmt.Println("Sending pings to all registered entries")

	file, err := os.Open("../files/sidecarConfig.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	var response string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for _, value := range lines {
		res := strings.Split(value, ";")
		fmt.Println(res)
		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", "www.google.com")
		//ra, err := net.ResolveIPAddr("ip4:icmp", res[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p.AddIPAddr(ra)
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			//response = fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)

			response = fmt.Sprintf("%s ID: %s IP Addr: %s receive, RTT: %v; \n ", response, res[0], addr.String(), rtt)

		}
		p.OnIdle = func() {

		}
		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	return &MessageReply{Message: response}, nil
}

func SafeToFile(ip string, sid string, tst string) {

	b, err := ioutil.ReadFile("../files/sidecarConfig.csv")
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
		fmt.Println("Debug: Saving the connection data")

		f, err := os.OpenFile("../files/sidecarConfig.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
		fmt.Println(l, "Debug: Bytes written successfully")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Debug: ID already stored")

		file, err := os.Open("../files/sidecarConfig.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		var lines []string

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			lines = append(lines, line)
		}

		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], suchRegex) {
				//lines[i]= sid + ";" + ip + ";" + "Timestamp" + ";"

				lines[i] = sid + ";" + ip + ";" + tst + ";"
				break
			}
		}

		f, err := os.OpenFile("../files/sidecarConfig.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
