package ashukhov

import java.io.File
import java.util.Properties

import org.apache.kafka.clients.producer.{KafkaProducer, ProducerRecord}
import org.apache.spark.sql._
import serialize.Serializer.{deserialize, serialize}
// import org.apache.spark.sql.DataFrame
import chat.cloudevent.{CloudEventAttribute}
import org.apache.log4j.{Level, Logger}


object PushToKafka {

    /**
    * main-function calling to push to Kafka and stream from it
    */
    def main(args: Array[String]): Unit = {
        // TODO: Metals Extension
        // TODO: bloopInstall
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        // pushToKafka("corona_3", 1)
    }

    /**
    * Pushes to Kafka using the KafkaProducer
    */
    /*def pushToKafka(topic: String, amount: Int) : Unit ={
        val props = new Properties()
        props.put("bootstrap.servers", "localhost:9092")
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer")
        props.put("value.serializer", "org.apache.kafka.common.serialization.ByteArraySerializer")
        val producer = new KafkaProducer[String, Array[Byte]](props)

        val corona = new Warnstufen("Alser", "12042", "Wien", "3")
        val cloud = CloudEventAttribute("2021-02-02", corona)
        val event = CloudEvent("001", "Jakobs Hintereingang", "002", "Corona", cloud)

        val byteEvent = serialize(event)

        for(a <- 1 to amount) {
            val record = new ProducerRecord[String, Array[Byte]](topic, "key" + a, byteEvent)
            producer.send(record)
        }
        producer.close()
    }*/

}

