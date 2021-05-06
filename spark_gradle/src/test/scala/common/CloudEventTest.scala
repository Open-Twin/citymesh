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

@RunWith(classOf[JUnitRunner])
class CloudEventTest extends AnyFlatSpec with SparkSessionWrapper {

    Logger.getLogger("org").setLevel(Level.OFF)
    Logger.getLogger("akka").setLevel(Level.OFF)

    behavior of "spark"

    it should "create a session" in {
        spark.emptyDataFrame.show()
        // kein assert notwendig
    }

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
        assert(!persist_df.isEmpty)
        assert(!cloudDF.isEmpty)
        assert(persist_df.select("ipService").first().get(0) === cloudDF.selectExpr("ipService").first().get(0))
        assert(persist_df.select("source").first().get(0) === cloudDF.selectExpr("source").first().get(0))
    }

    it should "not allow direct access to data" in {
        val df = readData()
        val proc_data = processData(df)
        assertThrows[AnalysisException](df.select("value").collect())
        assertThrows[AnalysisException](proc_data.select("cloud.source").collect())
    }

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

    it should "start with correct configurations" in {
        assert(spark.conf.get("spark.sql.warehouse.dir") == new File(configMap("warehouseLocation")).getAbsolutePath)
        assert(spark.conf.get("spark.speculation") == "false")
    }
}
