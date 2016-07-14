package modeclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tily/gofibwait"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type Event struct {
	HomeId            int         `json:"homeId"`
	Timestamp         time.Time   `json:"timestamp"`
	EventType         string      `json:"eventType"`
	EventData         interface{} `json:"eventData"`
	OriginDeviceId    int         `json:"originDeviceId"`
	OriginDeviceClass string      `json:"originDeviceClass"`
	OriginDeviceIp    string      `json:"originDeviceIp"`
}

type Command struct {
	Action     string      `json:"action"`
	Parameters interface{} `json:"parameters"`
}

type Client struct {
	Endpoint   string
	Token      string
	HTTPClient *http.Client
}

func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("ModeCloud %s", c.Token))
	req.Header.Add("Content-Type", "application/json")
	return c.HTTPClient.Do(req)
}

func (c *Client) DoListen(url string, origin string, callback func(*websocket.Conn) error) {
	w := gofibwait.NewWaiter(60)
	ws := c.ConnectWS(url, origin)
	for {
		err := callback(ws)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			err := ws.Close()
			if err != nil {
			}
			w.Wait()
			ws = c.ConnectWS(url, origin)
		} else {
			w.Reset()
		}
	}
}

func (c *Client) ConnectWS(url string, origin string) *websocket.Conn {
	w := gofibwait.NewWaiter(60)
	for true {
		w.Wait()
		ws, err := websocket.Dial(url, "", origin)
		if err == nil {
			return ws
		}
	}
	return nil
}

type User struct {
	Client
}

type Device struct {
	Client
	DeviceId int
}

func NewDevice(endpoint string, token string, deviceId int) Device {
	return Device{Client: Client{Endpoint: endpoint, Token: token, HTTPClient: &http.Client{}}, DeviceId: deviceId}
}

func NewUser(endpoint string, token string) User {
	return User{Client: Client{Endpoint: endpoint, Token: token, HTTPClient: &http.Client{}}}
}

func (d *Device) TriggerEvent(event Event) (*http.Response, error) {
	body, _ := json.Marshal(event)
	url := fmt.Sprintf("https://%s/devices/%d/event", d.Endpoint, d.DeviceId)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	return d.DoRequest(req)
}

func (u *User) SendCommand(deviceId int, command Command) (*http.Response, error) {
	body, _ := json.Marshal(command)
	url := fmt.Sprintf("https://%s/devices/%d/command", u.Endpoint, deviceId)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	return u.DoRequest(req)
}

func (d *Device) ListenToCommands(callback func(Command)) {
	url := fmt.Sprintf("ws://%s/devices/%d/command?authToken=%s", d.Endpoint, d.DeviceId, d.Token)
	origin := fmt.Sprintf("http://%s/", d.Endpoint)
	d.DoListen(url, origin, func(ws *websocket.Conn) error {
		command := Command{}
		if err := websocket.JSON.Receive(ws, &command); err != nil {
			return err
		} else {
			callback(command)
			return nil
		}
	})
}

func (u *User) ListenToEvents(callback func(Event)) {
	url := fmt.Sprintf("ws://%s/userSession/websocket?authToken=%s", u.Endpoint, u.Token)
	origin := fmt.Sprintf("http://%s/", u.Endpoint)
	u.DoListen(url, origin, func(ws *websocket.Conn) error {
		event := Event{}
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			return err
		} else {
			callback(event)
			return nil
		}
	})
}
