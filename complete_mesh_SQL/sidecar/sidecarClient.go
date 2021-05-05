package sidecar

import (
	"bufio"
	_ "errors"
	"github.com/Open-Twin/citymesh/complete_mesh/ddns"
	"github.com/Open-Twin/citymesh/complete_mesh/master"
	_ "github.com/gogo/protobuf/proto"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	_ "os"
	_ "reflect"
	"strings"
)

const (
	grpcport = "9010"

)

func client(cloudmessage *CloudEvent) {
	// The client establishes a connection with the master service and tries to send cloudevent messages
	var masterip net.IP
	masterip = dnscon()
	log.Println(masterip.String())
	if masterip == nil {
		masterip = net.ParseIP("127.0.0.1")
		log.Println(masterip.String())
	}

	log.Println("Debug: Initialised Client")

	// create client for GRPC Server
	var conn *grpc.ClientConn

	// Gathering the stored IPs

	var ips []string

	ips = GetIp()


	//creds, _ := credentials.NewClientTLSFromFile("cert/server.crt", "")
	//conn, error := grpc.Dial(":"+grpcport, grpc.WithTransportCredentials(creds))
	//conn, error := grpc.Dial(masterip.String()+":"+grpcport, grpc.WithTransportCredentials(creds))
	conn, error := grpc.Dial(masterip.String()+":"+grpcport, grpc.WithInsecure())

	if error != nil {
		log.Println("no server connection could be established cause: %v", error)
	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := master.NewChatServiceClient(conn)

	message := master.CloudEvent{
		IdService:   cloudmessage.IdService,
		Source:      cloudmessage.Source,
		SpecVersion: cloudmessage.SpecVersion,
		Type:        cloudmessage.Type,
		Attributes:  nil,
		Data:        nil,
		IdSidecar:   cloudmessage.IdSidecar,
		IpService:   cloudmessage.IpService,
		IpSidecar:   cloudmessage.IpSidecar,
		Timestamp:   cloudmessage.Timestamp,
	}

	response, err := c.DataFromSidecar(context.Background(), &message)

	if err != nil {
		log.Println("No con")
	}

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Message)

	} else {
		log.Println("Debug: Could not establish a connection")

		for _, newIP := range ips {
			log.Println("Debug: New Master Ip:")
			log.Println(newIP)
			conn, erro := grpc.Dial(newIP+":"+grpcport, grpc.WithInsecure())
			defer conn.Close()
			if erro != nil {
				log.Println(erro)
			}
			c := master.NewChatServiceClient(conn)
			response, err := c.DataFromSidecar(context.Background(), &message)

			if err != nil {
				log.Println(err)
			}

			if response != nil {
				log.Println("Connected to a new master")

				log.Printf("Response from Master: %s ", response.Message)
				break
			}
		}
	}
}
func GetIp() []string{

	var ips []string
	file, err := os.Open("files/sidecarCon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ";")

		ips = append(ips, res[1])

	}
	return ips
}

func dnscon() net.IP{
	// Querrying for master addresses
	var masterip net.IP
	ip, err := ddns.Query("Master")
	log.Println("IP API:")
	log.Println(ip)
	if err != nil || len(ip) == 0 {
		log.Println("The Query could not find any entries")
	}else {
		masterip = ip[0]
	}
	return masterip
}