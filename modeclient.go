package modeclient

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Endpoint string
	Token    string
}

func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("ModeCloud %s", c.Token))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func (c *Client) DoListen(url string, origin string, callback func(*websocket.Conn)) {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic("Dial: " + err.Error())
	}
	for {
		select {
		default:
			callback(ws)
		}
	}
}

type User struct {
	Client
}

type Device struct {
	Client
	DeviceId int
}

func NewDevice(endpoint string, token string, deviceId int) Device {
	return Device{Client: Client{Endpoint: endpoint, Token: token}, DeviceId: deviceId}
}

func NewUser(endpoint string, token string) User {
	return User{Client: Client{Endpoint: endpoint, Token: token}}
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
	d.DoListen(url, origin, func(ws *websocket.Conn) {
		command := Command{}
		websocket.JSON.Receive(ws, &command)
		callback(command)
	})
}

func (u *User) ListenToEvents(callback func(Event)) {
	url := fmt.Sprintf("ws://%s/userSession/websocket?authToken=%s", u.Endpoint, u.Token)
	origin := fmt.Sprintf("http://%s/", u.Endpoint)
	u.DoListen(url, origin, func(ws *websocket.Conn) {
		event := Event{}
		websocket.JSON.Receive(ws, &event)
		callback(event)
	})
}
