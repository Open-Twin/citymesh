package master

import (
	//"bufio"
	_ "errors"
	//"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	//"github.com/golang/protobuf/ptypes"

	//"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	//"github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/proto"

	_ "reflect"
)

func client(cloudmessage *CloudEvent) {

	//c := broker.NewChatServiceClient(conn)

	//broker

	//data := CloudEvent_ProtoData{cloudmessage.Data}

	newdata := cloudmessage.Data
	print(newdata)
	/*data := ptypes.UnmarshalAny(cloudmessage.Data)


	message := broker.CloudEvent{
		IdService:   cloudmessage.IdService,
		Source:      cloudmessage.Source,
		SpecVersion: cloudmessage.SpecVersion,
		Type:        cloudmessage.Type,
		Attributes:  nil,
		Data:        cloudmessage.Data,
		IdSidecar:   cloudmessage.IdSidecar,
		IpService:   cloudmessage.IpService,
		IpSidecar:   cloudmessage.IpSidecar,
		Timestamp:   cloudmessage.Timestamp,
	}

	response, err := c.DataFromSidecar(context.Background(), &message)

	if response != nil {
		log.Printf("Response from Server: %s , ", response.Message)

	} else {
		fmt.Println("Debug: Could not establish a connection")

		for _, newIP := range ips {
			fmt.Println("Debug: New Sidecar Ip:")
			fmt.Println(newIP)
			//conn, erro := grpc.Dial("target", grpc.WithInsecure())
			//defer conn.Close()
			//c := sidecar.NewChatServiceClient(conn)
			//response, err := c.DataFromSidecar(context.Background(), &message2)

			/*
				if response != nil {
					fmt.Println("New Target gesichtet")
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Response from Server: %s ", response.Reply)
					break
				}


		}
	}*/
}
