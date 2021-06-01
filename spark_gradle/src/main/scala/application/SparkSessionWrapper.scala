package application

import org.apache.spark.sql.SparkSession
import java.io.File

/** creates a SparkSession with the given properties */
trait SparkSessionWrapper extends Object with PropertiesReader{
    // the variable with the SparkSession
    lazy val spark: SparkSession = {
      SparkSession
        .builder()
        .appName("Spark Structured Streaming CloudEvent")
        .master(configMap("master")) // set master to by default local[*]
        .config("spark.sql.warehouse.dir",
          new File(configMap("warehouseLocation")).getAbsolutePath) // set the location of the warehouse
        .config("spark.speculation","false")
        .enableHiveSupport() // enables support to use hive for sql queries
        .getOrCreate()
    }
}
