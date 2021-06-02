package sidecarCom

import (
	"fmt"
	"github.com/Open-Twin/citymesh/complete_mesh/sidecar"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	_ "gotest.tools/assert"
	"reflect"
	"testing"
)

func TestServer_DataFromService(t *testing.T) {
	type args struct {
		ctx     context.Context
		message *sidecar.CloudEvent
	}
	tests := []struct {
		name        string
		args        args
		want        *sidecar.MessageReply
		wantErr     bool
		req         *sidecar.CloudEvent
		message     string
		expectedErr bool
	}{
		{
			name: "req ok",
			req: &sidecar.CloudEvent{
				IdService:   "123",
				Source:      "123",
				SpecVersion: "123",
				Type:        "Test",
				Attributes:  nil,
				Data:        nil,
				IdSidecar:   "",
				IpService:   "",
				IpSidecar:   "23",
				Timestamp:   "2021",
			},
			message:     "hello me",
			expectedErr: false,
		},
		{
			name:        "req with empty message",
			req:         &sidecar.CloudEvent{},
			expectedErr: true,
		},
		{
			name:        "nil request",
			req:         nil,
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &sidecar.Server{}
			//greeter := sidecar.New()
			g := NewGomegaWithT(t)
			ctx := context.Background()
			got, err := s.DataFromService(ctx, testCase.req)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataFromService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataFromService() got = %v, want %v", got, tt.want)
			}

			if tt.expectedErr {
				g.Expect(got).ToNot(BeNil(), "Result should be nil")
				g.Expect(err).ToNot(BeNil(), "Result should be nil")
			} else {
				g.Expect(got.Message).To(Equal(testCase.message))
			}
		})
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
