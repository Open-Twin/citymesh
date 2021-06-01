package application

import application.StreamingCloudevent.logger
import org.apache.log4j.{Level, Logger}
import broker.{Any, CloudEvent}
import org.apache.spark.sql._
import org.apache.spark.sql.streaming.StreamingQuery

/** main object of the big-data interface */
object StreamingCloudevent extends Object with PropertiesReader with SparkSessionWrapper {

  private val logger = Logger.getRootLogger
  /**
  * main-function calling to read Kafka data and stream to warehouse
  */
  def main(args: Array[String]): Unit = {
      // turn off debug logging for better overview
      // is optional
      //Logger.getLogger("org").setLevel(Level.OFF)
      //Logger.getLogger("akka").setLevel(Level.OFF)
      streamKafkaCloudEvent()
  }

  /**
   * function to call every process of the big-data interface
   */
  def streamKafkaCloudEvent(): Unit = {
        logger.info("Spark Session has been built.")
        logger.info("Spark master set to: " + configMap("master"))
        logger.info("warehouse location set to: " + configMap("warehouseLocation"))
        // read data from Kafka
        val df_read = readData()
        logger.info("Reading data from Kafka topic: " + configMap("kafka_topic"))

        // process data from proto format
        val proc_data = processData(df_read)

        // persist data in file system
        persistData(proc_data)
    }

  /**
   * reads a stream of data from Kafka with given properties
   *
   * @return the read dataframe from Kafka
   */
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

  /**
   * function to process binary protobuf data with ScalaPB into original format
   *
   * @param df_read the streaming dataframe with binary protobuf data
   * @return the processed original data as a DataFrame
   */
  def processData(df_read: DataFrame): DataFrame = {
      import spark.implicits.StringToColumn
      import scalapb.spark.ProtoSQL
      import scalapb.spark.Implicits._

      // parse the value of the dataframe as a byte array and expand the struct
      val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }
      val df_new = df_read.withColumn("cloud", parseCloud($"value"))
      logger.info("Processing data into original format")
      df_new
  }

  /**
   * function which takes the processed dataframe and stores it in the file system
   *
   * @param proc_data the processed original dataframe for persisting
   */
  def persistData(proc_data: DataFrame): Unit = {
        // sql statement to create table, yet not working fully
        import spark.sql
        sql("CREATE TABLE IF NOT EXISTS " + configMap("kafka_topic") + "(idService STRING, idSidecar STRING, " +
          "ipService STRING, ipSidecar STRING, source STRING, spec_version STRING, text_data STRING, timestamp STRING, " +
          "type STRING) STORED AS ORC")
        logger.info("Created table in warehouse")

        // store the processed original data in the warehouse in json format
        logger.info("Persisting processed data in warehouse")
        proc_data.select("cloud.*", "*")
          .drop("value", "cloud")
          .writeStream
          .format("orc")
          .option("path", "spark-warehouse/" + configMap("kafka_topic"))
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()
    }
}

