package common

import application.StreamingCloudevent.{logger, streamKafkaCloudEvent}
import broker.CloudEvent
import org.apache.log4j.{Level, Logger}
import org.junit.runner.RunWith
import org.scalatest.flatspec.AnyFlatSpec
import org.scalatest._
import org.scalatestplus.junit.JUnitRunner
import org.apache.spark.sql.{SaveMode, SparkSession}

@RunWith(classOf[JUnitRunner])
class CloudEventTest extends AnyFlatSpec with SparkSessionTestWrapper {

    Logger.getLogger("org").setLevel(Level.OFF)
    Logger.getLogger("akka").setLevel(Level.OFF)

    behavior of "spark"

    it should "create a session" in {
        spark.emptyDataFrame.show()
        // kein assert notwendig
    }

    it should "read a non-empty dataframe from Kafka" in {
        val df = spark
            .readStream
            .format("kafka")
            .option("kafka.bootstrap.servers", "localhost:9092")
            .option("subscribe", "topic-test")
            .option("failOnDataLoss", "false")
            .load()
            .select("value")

        assert(df.isStreaming)
        df.writeStream.format("console").outputMode("append").option("checkpointLocation", "test_checkpoint").start().awaitTermination(1000)
    }

    it should "store Kafka data in file system" in {
        import spark.sql

        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._

        val cloudDF = spark.read.json("attribute.json")

        cloudDF.printSchema()

        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._


        import spark.sql
        import org.apache.spark.sql.functions._
        import org.apache.spark.sql.Column

        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", "localhost:9092")
          .option("subscribe", "topic_test")
          .option("failOnDataLoss", "false")
          .load()
          .select("value")


        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }

        val df_new = df.withColumn("cloud", parseCloud($"value"))

        df_new.select("cloud.*", "*").drop("value")
          .writeStream
          .format("json")
          .option("path", "spark-warehouse/test_corona")
          .option("checkpointLocation", "test_checkpoint")
          .start()
          .awaitTermination(20000)
        //val sql_df_kafka = sql("SELECT * FROM test_corona")
        val persist_df = spark.read.json("spark-warehouse/test_corona/*.json")
        assert(persist_df.select("cloud.ipService").first().get(0) === cloudDF.selectExpr("ipService").first().get(0))
    }
}
