package healthcheck

import (
	"bufio"
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/chat"

"strings"

	//"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func Healthcheck(masterIDs []string, masterIps []string) {

	var data string
	for index, ip  := range masterIps {
		fmt.Println("Starting Healthcheck for Master : "+masterIDs[index] +"|"+ ip)
		now := time.Now()
		// Placeholder until Ips
		conn, erro := grpc.Dial(":9001", grpc.WithInsecure())
		// real code:
		//conn, err := grpc.Dial(ip+":9001", grpc.WithInsecure())

		if erro != nil {
			log.Fatalf("no server connection could be established cause: %v", erro)

		}

		// defer runs after the functions finishes
		defer conn.Close()

		c := chat.NewChatServiceClient(conn)

		var message chat.Health
		var sidecarIDs []string
		var sidecarIPs []string

		message = chat.Health{Health: "abc"}
		response, err := c.HealthCheck(context.Background(), &message)
		fmt.Println("Message gelesen:"+response.Message)
		if response != nil {
			//fmt.Sprintf("Healthcheck Report: %s ", response.Message)
			scanner := bufio.NewScanner(strings.NewReader(response.Message))
			for scanner.Scan() {
				res := strings.Split(scanner.Text(),";")
				sidecarIDs = append(sidecarIDs, res[0])
				sidecarIPs = append(sidecarIPs, res[1])
			}
			fmt.Println(sidecarIDs)
			fmt.Println(sidecarIPs)
			data+= "Master: "+masterIDs[index] +" | "+ip+" Timestamp: " +now.String()+"\n"
			//fmt.Println(response.Message)
			fmt.Println("Sidecars Hello?")
			sideData := SidecarHealth(sidecarIDs,sidecarIPs)
			fmt.Println(sideData)
			data += sideData
			fmt.Println("Wo seids?")

		}

		if err != nil {
			log.Printf("Response from Server: %s , ", err)
		}

	}
	fmt.Println(data)
	saveHealthCheck(data)

}

func saveHealthCheck(data string) {

	f, err := os.OpenFile("files/healthCheckReport.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(data)
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
}
