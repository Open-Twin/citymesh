package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var responseData data

func main() {
	fmt.Println("Hello, Data!")
	get()
}

type data[] struct {
	string `json:"Stand"`
	Warnstufen[] struct {
		Gkz       string `json:"GKZ"`
		Name      string `json:"Name"`
		Region    string `json:"Region"`
		Warnstufe string `json:"Warnstufe"`
	} `json:"Warnstufen"`
}




func get() {
	response, err := http.Get("https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json")

	/*defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseData)*/

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &responseData)
	// fmt.Printf("Result: %v\n", responseData)
	fmt.Println("_______________________")


	fmt.Println(responseData.Warnstufen[0], responseData.Name, responseData.)
	os.Exit(0)

}


