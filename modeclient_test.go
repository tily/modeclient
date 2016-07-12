package modeclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func getTLSClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func Test_NewDevice(t *testing.T) {
	device := NewDevice("dummy-endpoint", "dummy-device-api-key", 1)
	expect := Device{
		Client: Client{
			Endpoint:   "dummy-endpoint",
			Token:      "dummy-device-api-key",
			HTTPClient: &http.Client{},
		},
		DeviceId: 1,
	}
	if !reflect.DeepEqual(device, expect) {
		t.Logf("Expected: %+v\n", expect)
		t.Logf("Actually: %+v\n", device)
		t.Fatal()
	}
}

func Test_Device_TriggerEvent(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "Hello")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	device := NewDevice(u.Host, "dummy-device-api-key", 1)
	device.HTTPClient = getTLSClient()
	device.TriggerEvent(Event{EventType: "Hello", EventData: map[string]string{"Hello": "World"}})
}

func Test_Device_ListenToCommands(t *testing.T) {
}

func Test_NewUser(t *testing.T) {
	user := NewUser("dummy-endpoint", "dummy-user-api-key")
	expect := User{
		Client: Client{
			Endpoint:   "dummy-endpoint",
			Token:      "dummy-user-api-key",
			HTTPClient: &http.Client{},
		},
	}
	if !reflect.DeepEqual(user, expect) {
		t.Logf("Expected: %+v\n", expect)
		t.Logf("Actually: %+v\n", user)
		t.Fatal()
	}
}

func Test_User_SendCommand(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "Hello")
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	user := NewUser(u.Host, "dummy-device-api-key")
	user.HTTPClient = getTLSClient()
	user.SendCommand(123, Command{Action: "Hello", Parameters: map[string]string{"Hello": "World"}})
}

func Test_User_ListenToEvents(t *testing.T) {
}
