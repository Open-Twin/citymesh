package common

import application.SparkSessionWrapper
import application.StreamingCloudevent.{configMap, logger, persistData, processData, readData, streamKafkaCloudEvent}
import broker.CloudEvent
import org.apache.log4j.{Level, Logger}
import org.junit.runner.RunWith
import org.scalatest.flatspec.AnyFlatSpec
import org.scalatest._
import org.scalatestplus.junit.JUnitRunner
import org.apache.spark.sql.{AnalysisException, DataFrame, Row, SaveMode, SparkSession}

import java.io.File
import scala.util.Try

/** Test class for the Structured Streaming of the Cloudevents */
@RunWith(classOf[JUnitRunner])
class PersistDataTest extends AnyFlatSpec with SparkSessionWrapper {
  /**
   * stores processed data in file system and asserts whether it's correct
   */
  it should "store Kafka data in file system" in {
    val cloudDF = spark.read.json("attribute.json")

    cloudDF.printSchema()

    val df = readData()
    val proc_data = processData(df)

    assert(proc_data.isStreaming)

    proc_data.select("cloud.*", "*")
      .drop("value", "cloud")
      .writeStream
      .format("json")
      .option("path", "spark-warehouse/" + configMap("kafka_topic"))
      .option("checkpointLocation", "checkpoint")
      .start()
      .awaitTermination(10000)

    val persist_df = spark.read.json("spark-warehouse/" + configMap("kafka_topic")+ "/*.json")
    // assert to see if persisted data is correct
    assert(persist_df.select("ipService").first().get(0) === cloudDF.selectExpr("ipService").first().get(0))
    assert(persist_df.select("source").first().get(0) === cloudDF.selectExpr("source").first().get(0))
  }

  it should "throw an exception when not persisting Cloudevent data" in {
    import spark.implicits._
    import org.apache.spark.sql.types.{StructType, StructField, StringType, IntegerType};
    val schema = StructType(Array(StructField("cloud.one", StringType), StructField("cloud.two", StringType)))
    val simple_seq = Seq(Row("one data", "two data"), Row("first data", "second data"))
    val df = spark.createDataFrame(
      spark.sparkContext.parallelize(simple_seq),schema)
    df.printSchema()
    df.show()
    assertThrows[AnalysisException](persistData(df))
  }

  it should "should not have empty data after persisting" in {

    val df = readData()
    val proc_data = processData(df)
    proc_data.select("cloud.*", "*")
      .drop("value", "cloud")
      .writeStream
      .format("json")
      .option("path", "spark-warehouse/" + configMap("kafka_topic"))
      .option("checkpointLocation", "checkpoint")
      .start()
      .awaitTermination(10000)

    val persist_df = spark.read.json("spark-warehouse/" + configMap("kafka_topic")+ "/*.json")
    assert(!persist_df.isEmpty)
  }
}
