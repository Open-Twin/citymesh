package test.ashukhov

import java.io.File

import org.apache.spark.sql._
import org.apache.log4j.{Level, Logger}


trait SparkSessionTestWrapper {
    val warehouseLocation = new File("spark-warehouse").getAbsolutePath
    lazy val spark = SparkSession
      .builder()
      .appName("Spark Structured Streaming Testing")
      .master("local[*]")
      .config("spark.sql.warehouse.dir", warehouseLocation)
      .enableHiveSupport()
      .getOrCreate()
}
