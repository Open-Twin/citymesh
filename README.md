# Citymesh - Guide

## Prerequirements

Before running the citymesh/complete_mesh, the certificates have to be generated with:

```bash
sh cert/gencertnew.sh
```



If you don't have docker installed on your system, install it with your favourite packet-manager. After the installation of `docker` and `docker-compose` you can start the Kafka-Broker with the command:

```bash
docker-compose -f docker-compose-expose.yml up -d
```

It can be stopped with

```bash
docker-compose -f docker-compose-expose.yml down
```



## How to run the citymesh

To run the whole citymesh:

go to the [complete_mesh](complete_mesh_sql) folder and run the following commands:

**Starting the sidecar-Server:**

```bash
go run cmd/sidecarMain.go
```

**Starting the client which grabs data from API:**

```bash
go run cmd/clientMain.go
```

**Starting the master which receives data from sidecar and sends it to the Kafka Broker**

````bash
go run cmd/masterMain.go 
````



Now the terminal outputs show the transported data





-------

**Todo Makefile**
