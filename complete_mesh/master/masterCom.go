package master

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/broker"

	//"github.com/tatsushid/go-fastping"
	//"github.com/Open-Twin/citymesh/complete_mesh/broker"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	//"net"
	"os"
	"regexp"
	"strings"
	//"time"
)

/*
type Server struct {
}

func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	//panic("implement me")
}*/

var msg *CloudEvent
var data *broker.CloudEvent

func (s *Server) DataFromMaster(ctx context.Context, message *CloudEvent) {
	//log.Printf("Received message body from client %s , %s , %s , %s , %s ", message.IdService, message.Source, message.SpecVersion, message.Type, message.IdService, message.IpSidecar, message.IpSidecar, message.Timestamp, message.Data)
	//fmt.Printf( "Received message body from client")
	log.Printf("Received: %s", message.Data)
	msg = message
	//data = message.Data
	SafeToFileMaster(message.IdService, message.IpService, message.Timestamp)
	go client(msg)
}

func SafeToFileMaster(ip string, sid string, tst string) {

	b, err := ioutil.ReadFile("files/sidecarConfig.csv")
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

		f, err := os.OpenFile("files/sidecarConfig.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

		file, err := os.Open("files/sidecarConfig.csv")
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

		f, err := os.OpenFile("files/sidecarConfig.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
