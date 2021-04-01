package ashukhov

import scala.io.Source

object ReadProperties {
    def readProperties(): Map[String, String] ={
        val configMap=Source.fromFile("src/main/resources/application.properties").getLines().filter(line => line.contains("=")).map{ line => val tkns=line.split("=")
          if(tkns.size==1){
            (tkns(0) -> "" )
          }else{
            (tkns(0) ->  tkns(1))
          }
        }.toMap
        configMap
    }
}
