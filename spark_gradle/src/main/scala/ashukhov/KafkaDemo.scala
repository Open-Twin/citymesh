package ashukhov

import java.io.File

import serialize.Serializer.{serialize, deserialize}

import org.apache.spark.sql._

import org.apache.log4j.{Level, Logger}


import chat.cloudevent.{CloudEventAttribute, Warnstufen}

// makeCloudDF

import org.apache.spark.sql.types.{BinaryType, StructType}


import org.apache.commons.lang3.SerializationUtils

import org.apache.spark.sql.functions.col

import org.apache.spark.sql.functions.explode



object KafkaDemo {

    /**
    * main-function calling to stream Cloudevents from Kafka with Spark
    */
    def main(args: Array[String]): Unit = {
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        
        streamKafkaCloudEvent()
    }


    def streamKafkaCloudEvent(): Unit = {
        val warehouseLocation = new File("spark-warehouse").getAbsolutePath
        val spark = SparkSession
          .builder()
          .appName("Spark Structured Streaming CloudEvent")
          .master("local[*]")
          .config("spark.sql.warehouse.dir", warehouseLocation)
          .enableHiveSupport()
          .getOrCreate()

        // ScalaPB
        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._ 

        import spark.sql

        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", "localhost:9092")
          .option("subscribe", "topic_test")
          .option("failOnDataLoss", "false")
          .load()
          .select("value")

        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEventAttribute.parseFrom(bytes) }

        val df_new = df.withColumn("cloud", parseCloud($"value"))

        sql("CREATE TABLE IF NOT EXISTS corona (stand STRING, region STRING, gkz STRING, name STRING, warnstufe STRING) STORED AS PARQUET")

        df_new.select("cloud.*", "*")
              .withColumn("warnstufe", explode($"warnstufen"))
              .select("warnstufe.*", "*")
              .drop("warnstufen")
              .drop("value")
              .drop("cloud")
              .drop("warnstufe")
              .writeStream.format("parquet")
              .option("parquet.block.size", 1024)
              .option("path", "spark-warehouse/corona")
              .option("checkpointLocation", "checkpoint")
              .start()
              .awaitTermination()

    }

}

