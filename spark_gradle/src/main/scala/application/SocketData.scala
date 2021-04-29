package application

import java.io.File
import scala.io.Source

import org.apache.log4j.{Level, Logger}


import org.apache.spark.sql.SparkSession

object SocketData {

    def main(args: Array[String]): Unit = {
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        
        querySocket()
    }

    def querySocket(): Unit = {
      val configMap = readProperties()
      val warehouseLocation = new File(configMap("warehouseLocation")).getAbsolutePath
      val spark = SparkSession
        .builder()
        .appName("Spark Structured Streaming Example")
        .master(configMap("master"))
        .config("spark.sql.warehouse.dir", warehouseLocation)
        .enableHiveSupport()
        .getOrCreate()

      import spark.sql

      val df = spark.readStream
        .format("socket")
        .option("host","localhost")
        .option("port","3000")
        .load()

      df.writeStream
        .format("console")
        .outputMode("append")
        .start()
        .awaitTermination()
    }

    def readProperties(): Map[String, String] ={
        val configMap = Source.fromFile("src/main/resources/application.properties").getLines().filter(line => line.contains("=")).map{ line => val config=line.split("=")
          if(config.size==1) {
            (config(0) -> "" )
          } else {
            (config(0) ->  config(1))
          }
        }.toMap
        configMap
    }

}
