package test.ashukhov

import org.apache.spark.sql.SaveMode
import org.apache.spark.sql.functions.explode
import org.scalatest.funsuite.AnyFunSuite
import broker.CloudEvent

class CloudEventTest extends AnyFunSuite with SparkSessionTestWrapper {
    test("Json and Kafka") {
        import spark.sql

        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._

        val cloudDF = spark.read.json("attribute.json")

        sql("CREATE TABLE IF NOT EXISTS test_corona (stand STRING, region STRING, gkz STRING, name STRING, warnstufe STRING) STORED AS PARQUET")

        cloudDF.write.mode(SaveMode.Overwrite).saveAsTable("test_corona")

        val sql_df_cloud = sql("SELECT * FROM test_corona")

        sql("DELETE FROM test_corona WHERE stand=='2020-...'")

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

        df_new.select("cloud.*", "*")
          .withColumn("warnstufe", explode($"warnstufen"))
          .select("warnstufe.*", "*")
          .drop("warnstufen")
          .drop("value")
          .drop("cloud")
          .drop("warnstufe")
          .writeStream.format("parquet")
          .option("parquet.block.size", 1024)
          .option("path", "spark-warehouse/test_corona")
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()

        val sql_df_kafka = sql("SELECT * FROM test_corona")

        assert(sql_df_kafka.select("stand").first().get(0) === sql_df_cloud.select("stand").first().get(0))

    }
}
