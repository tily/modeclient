# modeclient

## Usage

### Devices

Create new device:

```
endpoint := "api.tinkermode.com"
deviceAPIKey := "Your API key of device's"
deviceId := 123 // Your device ID
device := modeclient.NewDevice(endpoint, deviceAPIKey, deviceId)
```

Trigger events:

```
device.TriggerEvent(modeclient.Event{EventType: "Hello", EventData: map[string]string{"Hello": "World"}})
```

Listen to commands:

```
device.ListenToCommands(func(command modeclient.Command) {
	fmt.Printf("Device received command: %+v\n", command)
})
```

### Users


Create new user:

```
endpoint := "api.tinkermode.com"
userAPIKey := "Your API key of user's"
user := modeclient.NewUser(endpoint, userAPIKey)
```

Send commands:

```
user.SendCommand(deviceId, modeclient.Command{Action: "Hello", Parameters: map[string]string{"Hello": "World"}})
```

Listen to events:

```
user.ListenToEvents(func(event modeclient.Event) {
	fmt.Printf("User received event: %+v\n", event)
})
```
