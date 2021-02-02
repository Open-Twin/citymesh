# Read Me

A Service Mesh is a structured network of communicating services.

This implementation here consists of three major layers.

* Master-Layer
* Sidecar-Layer
* Service-Layer

![](img\archi.png)

This project is written in GoLang and uses Protocol Buffers in combination with gRPC. [^1][^2][^3]

Go does not include all the needed libraries by default.

#### Needed installations:

```
go get google.golang.org/grpc
```

This Project consists of three main parts:

## 1) Client

The Client receives data from many different data sources: sensors, open data interfaces etc...

It can save messages which could not be send because of connection errors and also connect to other sidecars if the previous one has a failure.

The Client can be started by using: (```/cmd```)

```
go run clientMain.go
```

All the executable .go-files lie in the main package in the directory ```/cmd``` (goLang guidelines)

To operate the client needs multiple files, including: (```/files```)

* ```safeData.csv```

  This file stores all messages which could not be send because of network errors.

* ```clientCon.csv```

  ClientCon is the file which stores the IP-addresses of all the possible sidecars the client can connect to. Entries can be added manually.

```
Sidecar3;192.168.14.123:9000;
Sidecar4;192.168.33.156:9000;
```

```
<sidecar ID>;IP-address;
```

##### Proto

To establish a connection the client also needs access to the compiled Protobuf Files ```sidecar.pb.go``` and ```sidecar_grpc.pb.go```

```sidecar.proto``` can be compiled with:

```
 	//old Proto-Version
    /path/to/binary/protoc.exe (protoc) --		go_out=plugins=grpc:<foldername> ./<filename>.proto
    
    //new Proto-Version
    protoc --go_out=<foldername> <filename>.proto
```

Before compiling, the following Libraries should be installed:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

## 2) Sidecar

The sidecar is a proxy which receives the message from and hands it off to the master. The sidecars are also able to find a new master if the previous one breaks down. Each sidecar also saves the metadata of all clients which have previously created a connection. So that it can initiate health checks where the stats of the services is tested.

A sidecar can be started with: (```/cmd```)

```
go run sidecarMain.go
```

To operate the sidecar needs the following files: (```/sidecar```)

* ```sidecarServer.go```

  The server receives the messages from the clients and starts a new client which pushes the data to the next service in line.

* ```sidecarCom.go```

  The sidecarCom file defines the different gRPC functions which clients can invocate

* ```sidecarClient.go```

  The clients receives the data from the server and sends it to the masters

* ```sidecar.pb.go```

  Generated file for GRPC communication

* ```sidecar_grpc.pb.go```

  Generated file for GRPC communication

#### 3) Server (Master)

The Master receives the messages from the sidecar proxies and publishes them to to Kafka. The master also saves meta data of all sidecars which established a connection. The master can also initiate health-checks through the saved metadata.

The master can be started with: (```cmd```)

```
go run masterMain.go
```

To operate, the server needs the following files: (```/chat```)

* ```chat2.go```

  Chat2 holds all the gRPC communication functions.

* ```chat2.pb.go```

  Generated file for GRPC communication

* ```chat2.proto```

  Generated file for GRPC communication

* ```server.go```

  Starts the gRPC Server

### Tests

The test are in the directory ```/tests```. Each class has its own test file.

All tests can be started with:

```
go test
```

Passed tests look like this:

```
PASS
ok      github.com/Open-Twin/CityMesh-ProtoTypes/Jakob_GRPC/test        1.052s
```

### Healthchecks

Work in Progress

```/cmd```



```
go run masterHeartBeat.go
```



```
healthChechReport.csv
```





[^1]: https://golang.org/
[^2]: https://grpc.io/
[^3]: https://developers.google.com/protocol-buffers

