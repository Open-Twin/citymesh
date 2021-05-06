package application

import application.StreamingCloudevent.logger
import org.apache.log4j.{Level, Logger}
import broker.{Any, CloudEvent}
import org.apache.spark.sql._
import org.apache.spark.sql.streaming.StreamingQuery
import scalapb.spark.ProtoSQL


object StreamingCloudevent extends Object with PropertiesReader with SparkSessionWrapper {

    private val logger = Logger.getRootLogger
    /**
    * main-function calling to read Kafka data and stream from it to warehouse
    */
    def main(args: Array[String]): Unit = {
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        streamKafkaCloudEvent()
    }


    def streamKafkaCloudEvent(): Unit = {
        logger.info("Spark Session has been built.")
        logger.info("Spark master set to: " + configMap("master"))
        logger.info("warehouse location set to: " + configMap("warehouseLocation"))

        val df_read = readData()
        logger.info("Reading data from Kafka topic: " + configMap("kafka_topic"))

        val proc_data = processData(df_read)
        proc_data.printSchema()
        persistData(proc_data)
    }

    def readData(): DataFrame = {
        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", configMap("kafka_ip"))
          .option("subscribe", configMap("kafka_topic"))
          .option("failOnDataLoss", "false")
          .load()
          .select("value")
        df
    }

    def processData(df_read: DataFrame): DataFrame = {
        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL
        import scalapb.spark.Implicits._

        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }
        val df_new = df_read.withColumn("cloud", parseCloud($"value"))
        logger.info("Processing data into original format")
        df_new
    }

    def persistData(proc_data: DataFrame): Unit = {
        /*import spark.sql
        sql("CREATE TABLE IF NOT EXISTS " + configMap("kafka_topic") + "(idService STRING, source STRING, " +
          "spec_version STRING, type STRING, attributes STRUCT<key: STRING, value: STRING>, " +
          "binary_data BINARY, text_data STRING, proto_data STRUCT<type_url: STRING, value: BINARY>, " +
          "idSidecar STRING, ipService STRING, ipSidecar STRING, timestamp STRING)")
        logger.info("Created table in warehouse")*/

        logger.info("Persisting processed data in warehouse")
        proc_data.select("cloud.*", "*")
          .drop("value", "cloud")
          .writeStream
          .format("json")
          .option("path", "spark-warehouse/" + configMap("kafka_topic"))
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()
    }
}

