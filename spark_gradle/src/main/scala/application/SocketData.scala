package application

import java.io.File
import scala.io.Source

import org.apache.log4j.{Level, Logger}


import org.apache.spark.sql.SparkSession

/** object for reading data from a socket */
object SocketData extends Object with SparkSessionWrapper with PropertiesReader {
  /**
   * main function to call the socket reading job
   *
   * @param args command line arguments
   */
  def main(args: Array[String]): Unit = {
        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        
        querySocket()
    }

  /**
   * read data from a local socket with the port 3000
   */
  def querySocket(): Unit = {
      val df = spark.readStream
        .format("socket")
        .option("host","localhost")
        .option("port","3000")
        .load()
      // output the socket data to the console
      df.writeStream
        .format("console")
        .outputMode("append")
        .start()
        .awaitTermination()
    }
}
