<div align="center">
  
# 🍋 Lemonade Stand

> A centralized hub for interacting with [Lemonade](https://github.com/Jaezmien/lemonade).
 
</div>

# Usage
```bash
# Requires Go 1.24.0+ 

$ go mod tidy
$ go run .
$ sudo go run . # If you're running on Linux
```

# Routes

## GET `/`

Returns either an empty array, or a JSON array of the entire current state of the external memory region.

## GET `/ws`

| Query | Type | Required | Description |
| --- | --- | --- | --- |
| `appid` | `int` | yes | The AppID the client will be listening on |

### Receiving Messages

The client will be receiving a `binary` message, with the first byte of data indicating the type of message.

| Value | Description | Notes |
| --- | --- | --- |
| `0x1` | NotITG has initialized | |
| `0x2` | NotITG has exited | |
| `0x3` | NotITG has send a buffer to the client | The buffer content is the rest of the data after `0x3`. |

### Sending Messages

The websocket client can also send messages to the server, and it will be directly sent (and split if needed) as a buffer.

If the message type is a `text`, it will be automatically converted into a buffer.
