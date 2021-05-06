# Big Data Interface mit Spark

**Autor**: Vasco Shukhov
**Datum**: 02-02-2021

## Beschreibung

Spark ist mit Scala geschrieben und das Projekt wird mit dem Build-Management Tool Gradle ausgeführt.

## Versionierung

* Scala: 2.12
* JDK: 8
* Spark: 3.0.1
* Gradle: 6.8.1
* Entwicklung über eine Ubuntu VM: Ubuntu 18.04 Server

## Ausführung

### Streaming

Nebenbei sollte ein Docker Container für Kafka laufen. In diesem Repository gibt es einen kafka-docker Ordner dafür. 
Wenn man Gradle installiert hat, führt man das Spark Projekt folgend aus.

```bash
gradle stream
```
Der Task stream führt die Klasse KafkaDemo aus. Dort ist ein streaming-Job definiert, um Cloudevents (Corona-Ampel) von Kafka zu lesen. 
Diese werden entsprechend in ein DataFrame gewandelt und im Filesystem persistiert. Um Messages in Kafka reinzubekommen, geht man zum citymesh/complete_mesh/kafka Ordner. Hier führt man den Producer aus, während das Gradle Projekt läuft.

```bash
go run Producer.go
```

### Überprüfen

Der Task getAll führt eine SQL-Query aus und schaut sich den Inhalt des Filesystems an. 

```bash
gradle getAll
```
