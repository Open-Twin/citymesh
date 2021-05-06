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


Make sure, that data has already been persisted and send data to the Kafka broker. Some tests depend on a persisted system and other tests require data from the Kafka broker.