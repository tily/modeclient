# modeclient

[![GoDoc](https://godoc.org/github.com/tily/modeclient?status.svg)](http://godoc.org/github.com/tily/modeclient)

## Usage

### Devices

Create new device:

```go
endpoint := "api.tinkermode.com"
deviceAPIKey := "Your API key of device's"
deviceId := 123 // Your ID of device's
device := modeclient.NewDevice(endpoint, deviceAPIKey, deviceId)
```

Trigger events:

```go
device.TriggerEvent(modeclient.Event{EventType: "Hello", EventData: map[string]string{"Hello": "World"}})
```

Listen to commands:

```go
device.ListenToCommands(func(command modeclient.Command) {
	fmt.Printf("Device received command: %+v\n", command)
})
```

### Users


Create new user:

```go
endpoint := "api.tinkermode.com"
userAPIKey := "Your API key of user's"
user := modeclient.NewUser(endpoint, userAPIKey)
```

Send commands:

```go
user.SendCommand(deviceId, modeclient.Command{Action: "Hello", Parameters: map[string]string{"Hello": "World"}})
```

Listen to events:

```go
user.ListenToEvents(func(event modeclient.Event) {
	fmt.Printf("User received event: %+v\n", event)
})
```
