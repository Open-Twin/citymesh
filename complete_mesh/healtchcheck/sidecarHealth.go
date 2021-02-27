package healthcheck

import (
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)
func SidecarHealth(sidecarID []string ,sidecarIps []string) (string) {

	/*
		pinger, err := ping.NewPinger("www.google.com")
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run() // blocks until finished
		if err != nil {
			panic(err)
		}
		stats := pinger.Statistics() // get send/receive/rtt stats
		fmt.Println(stats)
	*/
	/*
		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", "www.google.com")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p.AddIPAddr(ra)
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		}
		p.OnIdle = func() {
			fmt.Println("finish")
		}
		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}

		/*
		host := "www.google.com"
		port := ""
		timeout := time.Duration(1 * time.Second)
		_, err := net.DialTimeout("tcp", host+":"+port, timeout)
		if err != nil {
			fmt.Printf("%s %s %s\n", host, "not responding", err.Error())
		} else {
			fmt.Printf("%s %s %s\n", host, "responding on port:", port)
		}
	*/
	fmt.Println("Infos")
	var data string
	for index, ip  := range sidecarIps {
		fmt.Println("Starting Healthcheck for Sidecar : "+sidecarID[index] +"|"+ ip)
		now := time.Now()
		// Placeholder until Ips
		conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
		// real code:
		//conn, err := grpc.Dial(ip+":9000", grpc.WithInsecure())

		if erro != nil {
			log.Fatalf("no server connection could be established cause: %v", erro)

		}

		// defer runs after the functions finishes
		defer conn.Close()

		c := sidecar.NewChatServiceClient(conn)

		var message sidecar.Health

		message = sidecar.Health{Health: "abc"}
		response, err := c.HealthCheck(context.Background(), &message)

		if response != nil {
			fmt.Sprintf("Healthcheck Report: %s ", response.Message)
			data+= "Sidecar: "+sidecarID[index] +" | "+ip+" Timestamp: " +now.String()+"\n"
			//fmt.Println(response.Message)
			data += response.Message


		}

		if err != nil {
			log.Printf("Response from Server: %s , ", err)
		}

	}
	fmt.Println(data)
	return data
}