package serialize

import java.io._

object Serializer {

  def serialize(obj: Any): Array[Byte] = {
    val byteOut = new ByteArrayOutputStream()
    val objOut = new ObjectOutputStream(byteOut)
    objOut.writeObject(obj)
    objOut.close()
    byteOut.close()
    byteOut.toByteArray
  }

  /*def deserialize[T <: Serializable](bytes: Array[Byte]): T = {
    val byteIn = new ByteArrayInputStream(bytes)
    val objIn = new ObjectInputStream(byteIn)
    val obj = objIn.readObject().asInstanceOf[T]
    byteIn.close()
    objIn.close()
    obj
  }*/

  def deserialize(bytes: Array[Byte]): Any = {
    val ois = new ObjectInputStream(new ByteArrayInputStream(bytes))
    val value = ois.readObject
    ois.close()
    value
  }
}