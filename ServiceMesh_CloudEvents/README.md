# ServiceMesh - How To

## Kafka - Producer and Consumer

Kafka Producer and Kafka Consumer are in `/Kafka` Director.
The scripts can be executed with the following commands.
```
go run Producer.go
```
```
go run Consumer.go
```

The Producer sends Data in corona-ampel or CoudData Format to the Kafka 
Broker. Data is saved in binary-Format via using the shopify sarama Library for Go.

The Consumer "pulls" the data from Kafka. Kafka is running in Docker and is listening
on Port `:9092`

## API - Grabber

The client is in `GETAPI` package and can gets Data from OpenData - API.
This example uses the corona-ampel. When starting the script with:
```
go run apiclient2.go
```
`http.Get` is used to get Data from the API. 

## CloudEvent Structure

