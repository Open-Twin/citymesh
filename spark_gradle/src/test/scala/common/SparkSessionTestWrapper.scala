package common

import org.apache.spark.sql.SparkSession

import java.io.File

trait SparkSessionTestWrapper {
  lazy val spark: SparkSession = {
    SparkSession.builder()
      .appName("Spark Structured Streaming Testing")
      .master("local[*]")
      .config("spark.sql.warehouse.dir", "spark-warehouse")
      .enableHiveSupport()
      .getOrCreate()
  }
}
