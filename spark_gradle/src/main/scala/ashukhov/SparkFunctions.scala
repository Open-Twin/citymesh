package ashukhov

import java.io.File

import org.apache.log4j.{Level, Logger}


import org.apache.spark.sql.SparkSession

object SparkFunctions {

    def main(args: Array[String]): Unit = {
        // TODO: Metals Extension
        // TODO: bloopInstall
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        
        getAll()
        //makeCloudDF()
    }
    /**
     * Streams from a Kafka topic to another topic
     */
    def streamToKafka(): Unit ={
      val spark = SparkSession
        .builder()
        .appName("Spark Structured Streaming Example")
        .master("local[*]")
        .getOrCreate()


      val df = spark
        .readStream
        .format("kafka")
        .option("kafka.bootstrap.servers", "localhost:9092")
        .option("subscribe", "corona_spark")
        .option("failOnDataLoss", "false")
        .load()

      df.selectExpr("CAST(value as STRING)")
        .writeStream
        .format("console")
        .outputMode("append")
        .option("checkpointLocation", "/tmp/spark_checkpoint")
        .start()
        .awaitTermination()

      /*import spark.implicits._

      val words = df.selectExpr("CAST(value AS STRING)").as[String]

      val query = words.writeStream
        .format("kafka")
        .option("kafka.bootstrap.servers", "localhost:9092")
        .option("checkpointLocation", "/tmp/spark")
        .option("topic", "kafka-demo")
        .outputMode("append")
        .start()

      query.awaitTermination()*/
    }


    /*
    * Pushes to Hive or SQL Storage from Kafka
    */

    def streamKafkaToHive(topicFrom: String, hiveUrl: String):Unit= {
      val warehouseLocation = new File("spark-warehouse").getAbsolutePath
      val spark = SparkSession
        .builder()
        .appName("Spark Structured Streaming Example")
        .master("local[*]")
        .config("spark.sql.warehouse.dir", warehouseLocation)
        .enableHiveSupport()
        .getOrCreate()

      import spark.sql
      import spark.implicits._


      /*sql("CREATE TABLE IF NOT EXISTS src(value STRING) STORED AS PARQUET")

      sql("INSERT INTO src VALUES('GOOD')")

      val words = df.selectExpr("CAST(value AS STRING)").as[String]

      val words_df = words.toDF("value")

      words_df.writeStream
        .format("parquet")
        .option("parquet.block.size", 1024)
        .option("path", "spark-warehouse/src")
        .option("checkpointLocation", "checkpoint")
        .start()
        .awaitTermination()*/
    }

    def getAll(): Unit = {
      val warehouseLocation = new File("spark-warehouse").getAbsolutePath
      val spark = SparkSession
        .builder()
        .appName("Spark Structured Streaming Example")
        .master("local[*]")
        .config("spark.sql.warehouse.dir", warehouseLocation)
        .enableHiveSupport()
        .getOrCreate()

      import spark.sql

      sql("SELECT * FROM corona").show()
    }

}
