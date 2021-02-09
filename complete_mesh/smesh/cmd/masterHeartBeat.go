package main

import (
	"fmt"
	"github.com/Open-Twin/citymesh/service_mesh/smesh/sidecar"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
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

	// Placeholder until Gateway
	conn, erro := grpc.Dial(":9000", grpc.WithInsecure())
	// real code:
	//conn, err := grpc.Dial(target, grpc.WithInsecure())

	if erro != nil {
		log.Fatalf("no server connection could be established cause: %v", erro)

	}

	// defer runs after the functions finishes
	defer conn.Close()

	c := sidecar.NewChatServiceClient(conn)

	var message sidecar.Health


	message = sidecar.Health{
		Reply:  "abc",
	}

	response, err := c.HealthCheck(context.Background(), &message)
	var data string
	if response != nil {
		fmt.Sprintf("Healthcheck Report: %s ", response.Reply)
		fmt.Println(data)
		saveHealthCheck(response.Reply)

	}

	if err != nil {
		log.Printf("Response from Server: %s , ", err)
	}


}

func saveHealthCheck(data string){

	f, err := os.OpenFile("../files/healthCheckReport.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString( data)
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