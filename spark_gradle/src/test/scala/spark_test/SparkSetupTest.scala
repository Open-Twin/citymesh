package spark_test

import application.SparkSessionWrapper
import application.StreamingCloudevent.{processData, readData}
import org.apache.log4j.{Level, Logger}
import org.apache.spark.sql.AnalysisException
import org.junit.runner.RunWith
import org.scalatest.flatspec.AnyFlatSpec
import org.scalatestplus.junit.JUnitRunner

import java.io.File

/** Test class for the Setup of Spark and Structured Streaming */
@RunWith(classOf[JUnitRunner])
class SparkSetupTest extends AnyFlatSpec with SparkSessionWrapper {
  // Turn off the debug logging for easier overview
  Logger.getLogger("org").setLevel(Level.OFF)
  Logger.getLogger("akka").setLevel(Level.OFF)

  behavior of "spark"

  /**
   * uses a simple function by the SparkSession
   * to see if it works properly without Exception
   */
  it should "create a session" in {
    spark.emptyDataFrame.show()
    // kein assert notwendig
  }

  /**
   * asserts, whether Spark throws an Exception when using streaming dataframes
   */
  it should "not allow direct access to data" in {
    val df = readData()
    val proc_data = processData(df)
    assertThrows[AnalysisException](df.select("value").collect())
    assertThrows[AnalysisException](proc_data.select("cloud.source").collect())
  }

  /**
   * assert, whether correct configurations have been set to the SparkSession
   */
  it should "start with correct configurations" in {
    assert(spark.conf.get("spark.sql.warehouse.dir") == new File(configMap("warehouseLocation")).getAbsolutePath)
    assert(spark.conf.get("spark.speculation") == "false")
  }

  /**
   * checks, if version of Spark is 3.0.1
   */
  it should "run with version of at least 3.0.1" in {
    assert(spark.version.equals("3.0.1"))
  }
}
