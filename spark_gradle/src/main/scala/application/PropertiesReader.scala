package application

import scala.io.Source

trait PropertiesReader {
  lazy val configMap: Map[String, String] = {
    val txtSource = Source.fromFile("src/main/resources/application.properties")
    val all = txtSource.getLines()
     all.filter(line => line.contains("="))
        .map{ line => val config=line.split("=")
        if(config.size==1){
          (config(0) -> "" )
        }else{
          (config(0) ->  config(1))
        }
    }.toMap
  }
}
