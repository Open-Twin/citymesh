package application

import application.StreamingCloudevent.configMap
import org.apache.spark.sql.SparkSession

import java.io.File

trait SparkSessionWrapper extends Object with PropertiesReader{
    lazy val spark: SparkSession = {
      SparkSession
        .builder()
        .appName("Spark Structured Streaming CloudEvent")
        .master(configMap("master"))
        .config("spark.sql.warehouse.dir", new File(configMap("warehouseLocation")).getAbsolutePath)
        .config("spark.speculation","false")
        .enableHiveSupport()
        .getOrCreate()
    }
}
