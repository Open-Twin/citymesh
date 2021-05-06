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