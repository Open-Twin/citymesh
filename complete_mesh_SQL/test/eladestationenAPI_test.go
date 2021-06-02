package test

import (
	"github.com/Open-Twin/citymesh/complete_mesh/GETAPI"
	"net/http"
	"net/http/httptest"

	"testing"
)

func getDataFromEladestationen(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "https://data.wien.gv.at/daten/geo?service=WFS&request=GetFeature&version=1.1.0&typeName=ogdwien:ELADESTELLEOGD&srsName=EPSG:4326&outputFormat=json")
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	//api := API{server.Client(), server.URL}
	body, err := GETAPI.APIData()

	ok(t, err)
	equals(t, []byte("OK"), body)
}




