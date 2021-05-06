package debug

import application.{PropertiesReader, SparkSessionWrapper}

object SelectAll extends Object with SparkSessionWrapper with PropertiesReader{

  def main(args: Array[String]): Unit = {
      val df = spark.read.json("spark-warehouse/"+configMap("kafka_topic")+"/*.json")
      df.show()
  }
}
