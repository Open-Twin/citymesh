package ddns

import (
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"strconv"
	"time"
)

type answerFormat struct {
	Domain string
	Error  string
	Value  string
}

const (
	lbip = "127.0.0.1"
	metaport = 20000
	dnsapiport = 10000
	dnsport = 52
)

func Register(hostname string,ip string, rtype string){
	msg := bson.M{
		"Hostname": hostname,
		"Ip" : ip,
		"RequestType" : rtype,
	}

	ans := SendBsonMessage(lbip+":"+strconv.Itoa(dnsapiport),msg)
	answerVals := answerFormat{}
	bson.Unmarshal(ans, &answerVals)
	if answerVals.Error != "ok" {
		log.Println("Registering has failed!")
	}
}

func SendBsonMessage(address string, msg bson.M) []byte {
	d := net.Dialer{Timeout: 10*time.Millisecond}
	conn, err := d.Dial("udp", address)
	defer conn.Close()

	if err != nil {
		fmt.Printf("Error on establishing connection: %s\n", err)
	}
	sendMsg, _ := bson.Marshal(msg)

	conn.Write(sendMsg)
	fmt.Printf("Message sent: %s\n", sendMsg)

	timeoutDuration := 5 * time.Second

	answer := make([]byte, 2048)
	erroro := conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	if err != nil {
		fmt.Printf("Error on receiving answer: %v", erroro.Error())
	}
	_, err = conn.Read(answer)

	if err != nil {
		fmt.Printf("Error on receiving answer: %v", err.Error())
	} else {
		fmt.Printf("Answer:\n%s\n", answer)
	}

	return answer
}

func Query(hostname string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, lbip+":"+strconv.Itoa(dnsport))
		},
	}
	return r.LookupIP(context.Background(), "ip", hostname)
}

func SendMetadataPacket(service string, rtype string, key string, value string) {
	address := lbip + ":" + strconv.Itoa(metaport)
	msg := bson.M{
		"service": service,
		"type": rtype,
		"key": key,
		"value": value,
	}

	SendBsonMessage(address, msg)
}