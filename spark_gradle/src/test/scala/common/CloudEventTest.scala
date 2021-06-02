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
class CloudEventTest extends AnyFlatSpec with SparkSessionWrapper {
    // Turn off the debug logging for easier overview
    Logger.getLogger("org").setLevel(Level.OFF)
    Logger.getLogger("akka").setLevel(Level.OFF)

    behavior of "spark"

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
     * check whether processed data is actually in correct structure
     */
    it should "have specified variables in processed and persisted data" in {
        def hasColumn(df: DataFrame, path: String) = Try(df(path)).isSuccess

        val persist_df = spark.read.json("spark-warehouse/" + configMap("kafka_topic")+ "/*.json")
        assert(hasColumn(persist_df, "ipService"))
        assert(hasColumn(persist_df, "idService"))
        assert(hasColumn(persist_df, "ipSidecar"))
        assert(hasColumn(persist_df, "idSidecar"))
        assert(hasColumn(persist_df, "spec_version"))
        assert(hasColumn(persist_df, "type"))
        assert(hasColumn(persist_df, "text_data"))
    }

    /**
     * Exception should be thrown when persisting gets the wrong binary data
     */
    it should "throw an exception when given wrong binary data" in {
        import spark.implicits._
        val df = Seq("default data").toDF
        val new_df = processData(df)
        assertThrows[AnalysisException](persistData(new_df))
    }

    /**
     * without a created table, an Exception should be thrown when accessed
     */
    it should "not allow selects without tables" in {
        import spark.sql
        assertThrows[AnalysisException](sql("select * from " + configMap("kafka_topic")))
    }

    /**
     * sequenced data can be processed without being correct
     */
    it should "be able to process sequenced data" in {
        import spark.implicits._
        val df = Seq(("default data", "data 1"), ("second default data", "data 2")).toDF
        assertThrows[AnalysisException](processData(df))
    }

    /**
     * non streaming dataframes can be accessed
     */
    it should "be able to access non-streaming data" in {
        import spark.implicits._
        val df = Seq("default data").toDF
        val proc_data = processData(df)
        proc_data.selectExpr("value")
    } 

}
