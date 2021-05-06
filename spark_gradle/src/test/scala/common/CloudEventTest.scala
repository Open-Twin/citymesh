package common

import application.SparkSessionWrapper
import application.StreamingCloudevent.{configMap, logger, persistData, processData, readData, streamKafkaCloudEvent}
import broker.CloudEvent
import org.apache.log4j.{Level, Logger}
import org.junit.runner.RunWith
import org.scalatest.flatspec.AnyFlatSpec
import org.scalatest._
import org.scalatestplus.junit.JUnitRunner
import org.apache.spark.sql.{AnalysisException, SaveMode, SparkSession}

import java.io.File

/** Test class for the Structured Streaming of the Cloudevents */
@RunWith(classOf[JUnitRunner])
class CloudEventTest extends AnyFlatSpec with SparkSessionWrapper {
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
     * reads a stream of data from Kafka
     * asserts whether it is streaming
     */
    it should "read a non-empty dataframe from Kafka" in {
        val df = readData()
        assert(df.isStreaming)
        df.writeStream
          .format("console")
          .outputMode("append")
          .option("checkpointLocation", "test_checkpoint")
          .start()
          .awaitTermination(1000)
    }

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
        // no dataframe should be empty
        assert(!persist_df.isEmpty)
        assert(!cloudDF.isEmpty)
        // assert to see if persisted data is correct
        assert(persist_df.select("ipService").first().get(0) === cloudDF.selectExpr("ipService").first().get(0))
        assert(persist_df.select("source").first().get(0) === cloudDF.selectExpr("source").first().get(0))
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
     * processes read data from Kafka and asserts whether processing works
     */
    it should "process into multiple columns" in {
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
        assert(persist_df.columns.length > 0)
    }

    /**
     * assert, whether correct configurations have been set to the SparkSession
     */
    it should "start with correct configurations" in {
        assert(spark.conf.get("spark.sql.warehouse.dir") == new File(configMap("warehouseLocation")).getAbsolutePath)
        assert(spark.conf.get("spark.speculation") == "false")
    }
}
