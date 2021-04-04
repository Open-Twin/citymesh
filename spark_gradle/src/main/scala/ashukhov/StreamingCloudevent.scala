package ashukhov

import ReadProperties.readProperties

import java.io.{File, FileNotFoundException}
import org.apache.log4j.{Level, Logger}
import broker.{Any, CloudEvent}
import org.apache.spark.sql._

import java.util.Properties
import scala.io.Source


object StreamingCloudevent {

    private val logger = Logger.getRootLogger
    /**
    * main-function calling to push to Kafka and stream from it
    */
    def main(args: Array[String]): Unit = {

        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        streamKafkaCloudEvent()
        //makeCloudDF()
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

        // ScalaPB
        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._ 

        
        import spark.sql

        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", configMap("kafka_ip"))
          .option("subscribe", configMap("kafka_topic"))
          .option("failOnDataLoss", "false")
          .load()
          .select("value")

        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }

        // val parseBytes = ProtoSQL.udf {bytes: Array[Byte] => bytes.toString()}

        val df_new = df.withColumn("cloud", parseCloud($"value"))

        df_new.select("cloud.*", "*").drop("value")//.drop("cloud")
          .select("proto_data.*", "*")
          //.withColumn("corona", unpackAny($"proto_data"))
          .selectExpr("CAST(value AS STRING)")
          .writeStream
          .format("json")
          .option("path", "spark-warehouse/corona")
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()
         
          
          /*.writeStream.format("console")
          .outputMode("append")
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()*/


        /*df_new.select("cloud.*", "*")
              .withColumn("warnstufe", explode($"warnstufen"))
              .select("warnstufe.*", "*")
              .drop("warnstufen", "value", "cloud", "warnstufe")
              .writeStream.format("json")
              //.option("parquet.block.size", 1024)
              .option("path", "spark-warehouse/corona")
              .option("checkpointLocation", "checkpoint")
              .start()
              .awaitTermination()*/
        
    }
}

