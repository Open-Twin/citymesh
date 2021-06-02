package client

import (
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	. "github.com/onsi/gomega"
	"testing"
)

func TestApiclient(t *testing.T) {
	tests := []struct {
		name                  string
		wantCloudeventmessage sidecar.CloudEvent
		ampel                 string
		expectedErr           bool
	}{
		{
			name:        "getStructure",
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//testCase := tt
			gotCloudeventmessage := Apiclient()
			g := NewGomegaWithT(t)

			if tt.expectedErr {
				g.Expect(gotCloudeventmessage).ToNot(BeNil(), "Result should be nil")
			}

		})
	}
}

func TestWrongURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name          string
		args          args
		wantAmpelJson string
		url           string
		expectedErr   bool
	}{
		{
			name:        "getWrongURL",
			url:         "http",
			expectedErr: true,
		},

		/*{
			name: "messagesCreated",
		},
		{
			name: "messageSent",
		},*/
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			g := NewGomegaWithT(t)
			gotAmpelJson := URL(testCase.url)
			if tt.expectedErr {
				g.Expect(gotAmpelJson).ToNot(BeNil(), "Result should be nil")
			}

		})
	}
}

func TestURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name          string
		args          args
		wantAmpelJson string
		url           string
		expectedErr   bool
	}{
		{
			name:        "getURL",
			url:         "https://corona-ampel.gv.at/sites/corona-ampel.gv.at/files/assets/Warnstufen_Corona_Ampel_Gemeinden_aktuell.json",
			expectedErr: true,
		},

		/*{
			name: "messagesCreated",
		},
		{
			name: "messageSent",
		},*/
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			g := NewGomegaWithT(t)
			gotAmpelJson := URL(testCase.url)
			if tt.expectedErr {
				g.Expect(gotAmpelJson).ToNot(BeNil(), "Result should be nil")
			}

		})
	}
}
