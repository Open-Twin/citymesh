package debug

import application.{PropertiesReader, SparkSessionWrapper}

/** object to read data from file system */
object SelectAll extends Object with SparkSessionWrapper with PropertiesReader{
  /**
   * main function to read data from file system and output it to console
   * @param args command line arguments
   */
  def main(args: Array[String]): Unit = {
      val df = spark.read.json("spark-warehouse/"+configMap("kafka_topic")+"/*.json")
      df.show()
  }
}
