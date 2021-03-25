package ashukhov

import java.io.File

import serialize.Serializer.{serialize, deserialize}

import org.apache.spark.sql._

import org.apache.log4j.{Level, Logger}


// import chat.cloudevent.{CloudEventAttribute, Warnstufen}

import broker.{CloudEvent, Any}

import corona.{Corona, Warnstufe, Values}


// makeCloudDF

import org.apache.spark.sql.types.{BinaryType, StructType}


import org.apache.commons.lang3.SerializationUtils

import org.apache.spark.sql.functions.col

import org.apache.spark.sql.functions.explode


object KafkaDemo {

    /**
    * main-function calling to push to Kafka and stream from it
    */
    def main(args: Array[String]): Unit = {

        Logger.getLogger("org").setLevel(Level.OFF)
        Logger.getLogger("akka").setLevel(Level.OFF)
        
        streamKafkaCloudEvent()
        //makeCloudDF()
    }


    def streamKafkaCloudEvent(): Unit = {
        val warehouseLocation = new File("spark-warehouse").getAbsolutePath
        val spark = SparkSession
          .builder()
          .appName("Spark Structured Streaming CloudEvent")
          .master("local[*]")
          .config("spark.sql.warehouse.dir", warehouseLocation)
          .config("spark.speculation","false")
          .enableHiveSupport()
          .getOrCreate()

        // ScalaPB
        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._ 

        
        import spark.sql

        
        val df = spark
          .readStream
          .format("kafka")
          .option("kafka.bootstrap.servers", "localhost:9092")
          .option("subscribe", "topic_test")
          .option("failOnDataLoss", "false")
          .load()
          .select("value")

        // df.select("value").writeStream.format("console").outputMode("append").option("checkpointLocation", "checkpoint").start().awaitTermination()

        val parseCloud = ProtoSQL.udf { bytes: Array[Byte] => CloudEvent.parseFrom(bytes) }

        // val parseValue = ProtoSQL.udf {bytes: Array[Byte] => ByteString.parseFrom(bytes)}

        // val parseBytes = ProtoSQL.udf {bytes: Array[Byte] => bytes.toString()}

        val unpackAny = ProtoSQL.udf { any: com.google.protobuf.any.Any => any.unpack[Corona] }


        val df_new = df.withColumn("cloud", parseCloud($"value"))


       df_new.select("cloud.*", "*").drop("value")//.drop("cloud")
          //.select("proto_data.*", "*")
          .withColumn("corona", unpackAny($"proto_data"))
          //.selectExpr("CAST(value AS STRING)")
          .writeStream
          .format("json")
          .option("path", "spark-warehouse/corona")
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()
         
          
          /*.writeStream.format("console")
          .outputMode("append")
          .option("checkpointLocation", "checkpoint")
          .start()
          .awaitTermination()*/


        /*df_new.select("cloud.*", "*")
              .withColumn("warnstufe", explode($"warnstufen"))
              .select("warnstufe.*", "*")
              .drop("warnstufen", "value", "cloud", "warnstufe")
              .writeStream.format("json")
              //.option("parquet.block.size", 1024)
              .option("path", "spark-warehouse/corona")
              .option("checkpointLocation", "checkpoint")
              .start()
              .awaitTermination()*/
        

        /*df_new.writeStream.foreach(new ForeachWriter[Row] {
            override def open(partitionId: Long, epochId: Long): Boolean = {return true}

            override def process(value: Row): Unit = {
                

                //val cloud = deserialize[CloudEventAttribute[CoronaAmpel]](value.apply(1).asInstanceOf[Array[Byte]])

                //val cloud = deserialize[CloudEventAttribute[CoronaAmpel]](bytes)
                //val seq_cloud = Seq(cloud)

                /*val df_cloud = seq_cloud.toDF()
                df_cloud.show()
                df_cloud.writeStream.format("parquet")
                  .option("parquet.block.size", 1024)
                  .option("path", "spark-warehouse/corona")
                  .option("checkpointLocation", "checkpoint")
                  .start()*/
            }

            override def close(errorOrNull: Throwable): Unit = {}
        }).start().awaitTermination()*/

    }

    def makeCloudDF(): Unit = {
        val warehouseLocation = new File("spark-warehouse").getAbsolutePath
        val spark = SparkSession
          .builder()
          .appName("Spark Structured Streaming CloudEvent")
          .master("local[*]")
          .config("spark.sql.warehouse.dir", warehouseLocation)
          .enableHiveSupport()
          .getOrCreate()

        import spark.implicits.StringToColumn
        import scalapb.spark.ProtoSQL

        import scalapb.spark.Implicits._
      

        /*val cloud1 = CloudEvent[Warnstufen]("001", "Allermannen Hintereingang", "002", "Corona",
            CloudEventAttribute("2021-02-02", 
              new Warnstufen("Alser", "12042", "Wien", "3")))
        
        val cloud2 = CloudEvent[Warnstufen]("002", "Jakobs Vordereingang", "003", "Corona", 
            CloudEventAttribute("2021-02-03", 
              new Warnstufen("Döbling", "11190", "Wien", "5")))*/

        /*val corona1 = Warnstufen("Döbling", "2150", "Wien", "7")
        val corona2 = Warnstufen("Innere Stadt", "55254", "NÖ", "9")

        val seq_corona = Seq(corona1, corona2)
        val cloud1 = CloudEventAttribute("2020-04-01", seq_corona)

        val df = ProtoSQL.createDataFrame(spark, Seq(cloud1))
        df.printSchema()
        df.show()*/
        // val cloudS = deserialize(serialize(seq_event.select($"value").first().get(0)))
        // println(cloudS)

        // cloudS.asInstanceOf[CloudEvent[CoronaAmpel]]

        val cloudDF = spark.read.json("attribute.json")
        cloudDF.printSchema()
        cloudDF.show()

        // val better = spark.createDataFrame(spark.sparkContext.parallelize(Seq(cloud1, cloud2)), StructType(cloudDF.schema))

        /*sql("CREATE TABLE IF NOT EXISTS corona (id STRING, source STRING, spec_version STRING, data_type STRING, attribute STRUCT<stand: STRING, message: STRUCT<region: STRING, gkz: STRING, name: STRING, warnstufe: STRING>>) STORED AS PARQUET")

        cloudDF.write.mode(SaveMode.Overwrite).saveAsTable("corona")

        sql("SELECT * FROM corona").show()*/
    }

}

