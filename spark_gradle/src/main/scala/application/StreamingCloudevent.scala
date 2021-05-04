package application

// import ReadProperties.readProperties

import java.io.{File, FileNotFoundException}
import org.apache.log4j.{Level, Logger}
import broker.{Any, CloudEvent}
import org.apache.spark.sql._

import java.util.Properties
import scala.io.Source

import org.apache.spark.sql.types.{StructType, StructField, StringType, ArrayType};


object StreamingCloudevent {

    private val logger = Logger.getRootLogger
    /**
    * main-function calling to push to Kafka and stream from it
    */
    def main(args: Array[String]): Unit = {
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        streamKafkaCloudEvent()
    }


    def streamKafkaCloudEvent(): Unit = {
        val configMap = readProperties()

        val warehousePath = new File(configMap("warehouseLocation")).getAbsolutePath
        val spark = SparkSession
          .builder()
          .appName("Spark Structured Streaming CloudEvent")
          .master(configMap("master"))
          .config("spark.sql.warehouse.dir", warehousePath)
          .config("spark.speculation","false")
          .enableHiveSupport()
          .getOrCreate()

        logger.info("Spark Session has been built.")
        logger.info("Spark master set to: " + configMap("master"))
        logger.info("warehouse location set to: " + warehousePath)

        // ScalaPB
        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._ 

        
        import spark.sql
        import org.apache.spark.sql.functions._
        import org.apache.spark.sql.Column

        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", configMap("kafka_ip"))
          .option("subscribe", configMap("kafka_topic"))
          .option("failOnDataLoss", "false")
          .load()
          .select("value")

        logger.info("Reading data from Kafka topic: " + configMap("kafka_topic"))  

        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }

        val df_new = df.withColumn("cloud", parseCloud($"value"))

        df_new.select("cloud.*", "*").drop("value")
          .writeStream
          .format("json")
          .option("path", "spark-warehouse/" + configMap("kafka_topic"))
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()
    }

    def readProperties(): Map[String, String] ={
        val configMap=Source.fromFile("src/main/resources/application.properties").getLines().filter(line => line.contains("=")).map{ line => val config=line.split("=")
          if(config.size==1){
            (config(0) -> "" )
          }else{
            (config(0) ->  config(1))
          }
        }.toMap
        configMap
    }

}

