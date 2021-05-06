# README

Last changes: 05.05.2021
Authors: Laback Jakob, Matschinegg Thomas, Arseniy Shukhov


# Citymesh - Guide

## Prerequirements

Before running the citymesh/complete_mesh_sql, the certificates have to be generated in the `cert` - folder with:

```bash
sh cert/gencertnew.sh
```

If you don't have docker installed on your system, install it with your favourite packet-manager. After the installation of `docker` and `docker-compose` you can start the Kafka-Broker with the command: (the files can be found in the [kafka-docker](https://github.com/Open-Twin/citymesh/tree/master/kafka-docker) directory)

```bash
docker-compose -f docker-compose-expose.yml up -d
```

It can be stopped with

```bash
docker-compose -f docker-compose-expose.yml down
```

## How to run the citymesh

To run the whole citymesh:

Go to the [complete_mesh_sql](https://github.com/Open-Twin/citymesh/tree/master/complete_mesh) folder and run the following commands:

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

## Deploying the Service Mesh

Instead of starting all services by hand the mesh can also be deployed using a docker-compose.yml file or Kubernetes Cluster-orchestrator

The Compose file can be found inside the `complete_mesh_sql` folder and can be deployed using:
```bash
docker-compose up -d -f <filename>
```
To deploy the mesh inside a Kubernetes node all important service depoyments have to be applied.

These deployments can be found inside the `Kubernetes` directory.

Before and building and deploying the services make sure your Kubernetes ist running. Minikube for example can be started with `minikube start`.

Before the services can be deployed the Dockerfiles have to be build.
```bash
sudo docker build -t digitaltwin/smesh -f DockerClient
sudo docker build -t digitaltwin/smesh2 -f DockerSidecar
sudo docker build -t digitaltwin/smesh3 -f DockerMaster
```
Afterwards the deployments can be applied.
```bash
sudo kubectl apply -f deploymentClient.yml
sudo kubectl apply -f deploymentSidecar.yml
sudo kubectl apply -f deploymentMaster.yml
```
Also the loadbalancing services have to be applied.
```bash
sudo kubectl apply -f serviceSidecar.yml
sudo kubectl apply -f serviceMaster.yml
```


## grpc Encryption

partially no working

## Big-Data Interface Guide
This section shows the last step of the Citymesh. It processes data from the Kafka broker and persists it in the file system.

### Beschreibung
The big-data interface mainly consists of a Spark project, which streams and persists data. The interface is in the directory ``spark_gradle``. Spark is written with Scala and the project uses Gradle as build management tool. 

### Versioning
* Scala: 2.12
* JDK: 11
* Spark: 3.0.1
* Gradle: 6.8.1
* Development with an Ubuntu VM: Ubuntu 18.04 Server

### Execution

#### Configuration
The big-data interface has multiple configurations that have to be set, as they can change for every system. The file ``application.properties`` in the directory ``src/main/resources`` can be changed to specify the master of the Spark session, the location of the warehouse and connection to the Kafka broker. By default it looks like this:

```properties
warehouseLocation=spark-warehouse
master=local[*]
kafka_ip=localhost:9092
kafka_topic=topic_test
```

#### Streaming
Parallel to the project, the docker container for Kafka has to run, otherwise the interface throws an exception.
To run the big-data interface you have to install Gradle. If you have done so, you can execute the interface followingly.
```bash=
gradle stream
```
The task stream executes the object ``StreamingCloudevent``. A streaming job is defined to read Cloudevent data from Kafka. It processes it into a dataframe and stores it. 
To send data to Kafka, please go to ``citymesh/complete_mesh/kafka`` and run the Producer.
```bash=
go run Producer.go
```
The data can be sent to the broker by another script too. As long as the data has the same structure as the Producer script, the interface will continue to work.


#### Verify the streaming
The task getAll looks into the warehouse and reads the persisted data. It shows the content in the console. This task is merely for debug purposes and has no additional features than to output persisted data to the console.
```bash=
gradle getAll
```
### Sockets

If you want to let the big-data interface read from a socket, the following command launches a Spark Job, which reads socket data from the port 3000.
```bash=
gradle socketData
```

### Tests

To execute the tests of the interface, you type in following command:

```bash=
gradle test
```

Make sure, that data has already been persisted and send data to the Kafka broker. Some tests depend on a persisted system and other tests require data from the Kafka broker.README

Last changes: 05.05.2021
Authors: Laback Jakob, Matschinegg Thomas, Arseniy Shukhov


# Citymesh - Guide

## Prerequirements

Before running the citymesh/complete_mesh_sql, the certificates have to be generated in the `cert` - folder with:

```bash
sh cert/gencertnew.sh
```



If you don't have docker installed on your system, install it with your favourite packet-manager. After the installation of `docker` and `docker-compose` you can start the Kafka-Broker with the command: (the files can be found in the [kafka-docker](https://github.com/Open-Twin/citymesh/tree/master/kafka-docker) directory)

```bash
docker-compose -f docker-compose-expose.yml up -d
```

It can be stopped with

```bash
docker-compose -f docker-compose-expose.yml down
```

## How to run the citymesh

To run the whole citymesh:

Go to the [complete_mesh_sql](https://github.com/Open-Twin/citymesh/tree/master/complete_mesh) folder and run the following commands:

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

Now the terminal outputs the transported data.

-------

## Deploying the Service Mesh

Instead of starting all services by hand the mesh can also be deployed using a docker-compose.yml file or Kubernetes Cluster-orchestrator

The Compose file can be found inside the `complete_mesh_sql` folder and can be deployed using:
```bash
docker-compose up -d -f <filename>
```
To deploy the mesh inside a Kubernetes node all important service deployments have to be applied.

These deployments can be found inside the `Kubernetes` directory.

Before and building and deploying the services make sure your Kubernetes is running. Minikube for example can be started with `minikube start`.

Before the services can be deployed the Dockerfiles have to be build.
```bash
sudo docker build -t digitaltwin/smesh -f DockerClient
sudo docker build -t digitaltwin/smesh2 -f DockerSidecar
sudo docker build -t digitaltwin/smesh3 -f DockerMaster
```
Afterwards the deployments can be applied.
```bash
sudo kubectl apply -f deploymentClient.yml
sudo kubectl apply -f deploymentSidecar.yml
sudo kubectl apply -f deploymentMaster.yml
```
Also the loadbalancing services have to be applied.
```bash
sudo kubectl apply -f serviceSidecar.yml
sudo kubectl apply -f serviceMaster.yml
```


## grpc Encryption

partly no working

## Running the Big-Data interface

### Beschreibung
Spark ist mit Scala geschrieben und das Projekt wird mit dem Build-Management Tool Gradle ausgeführt.

### Versionierung
* Scala: 2.12
* JDK: 11
* Spark: 3.0.1
* Gradle: 6.8.1
* Entwicklung über eine Ubuntu VM: Ubuntu 18.04 Server

### Ausführung

#### Konfiguration
application.properties

#### Streaming
Nebenbei sollte ein Docker Container für Kafka laufen. In diesem Repository gibt es einen kafka-docker Ordner dafür. Wenn man Gradle installiert hat, führt man das Spark Projekt folgend aus.
```bash=
gradle stream
```
Der Task stream führt die Klasse KafkaDemo aus. Dort ist ein streaming-Job definiert, um Cloudevents (Corona-Ampel) von Kafka zu lesen. Diese werden entsprechend in ein DataFrame gewandelt und im Filesystem persistiert. Um Messages in Kafka reinzubekommen, geht man zum citymesh/complete_mesh/kafka Ordner. Hier führt man den Producer aus, während das Gradle Projekt läuft.
```bash=
go run Producer.go
```
Die Daten können auch von einem anderen Script zum Kafka Broker gesendet werden. Solange die Daten dieselbe Struktur wie beim Producer-Skript haben kann das Big-Date Interface damit arbeiten.


#### Überprüfen
Der Task getAll liest das Persistierte im Warehouse und zeigt den Inhalt in der Konsole. Dieser Task ist lediglich Verwendung für Debugging und hat keine weiteren Funktionen.
```bash=
gradle getAll
```
### Sockets lesen

Wenn das Interface einen Socket lesen soll, dann führt man folgenden Befehl aus.
```bash=
gradle socketData
```