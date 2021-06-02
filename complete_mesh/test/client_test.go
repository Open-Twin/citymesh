package test

import (
	"log"
	"os"
)

func setup() {
	e := os.Remove("../files/saveData.csv")
	if e != nil {
		log.Fatal(e)
	}
}

/*func TestDataSave(t *testing.T) {
	setup()
	var message sidecar.Message

	message = sidecar.Message{
		Timestamp:  "16:00:00",
		Location:   "Alte Börse",
		Sensortyp:  "192.168.41.140",
		SensorID:   123,
		SensorData: "S1234",
	}

	client.DataSave(message)

	// Established a connection
	fmt.Println("Sending old data")
	file, err := os.Open("files/saveData.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		// Trying to push the old data
		res := strings.Split(line, ";")
		fmt.Println(res)
		if res[0] != "Alte Börse" {
			t.Errorf("Value 0 correct got: %s, want: %s.", res[0], "16:00:00")
		}
		if res[2] != "192.168.41.140" {
			t.Errorf("Value 0 correct got: %s, want: %s.", res[0], "16:00:00")
		}
		if res[4] != "S1234" {
			t.Errorf("Value 0 correct got: %s, want: %s.", res[0], "16:00:00")
		}
		if res[6] != "16:00:00" {
			t.Errorf("Value 0 correct got: %s, want: %s.", res[0], "16:00:00")
		}
	}

}

func TestNewIP(t *testing.T) {

}
*/
