<div align="center">
  
# 🍋 Lemonade Stand

> A centralized hub for interacting with [Lemonade](https://github.com/Jaezmien/lemonade).
 
</div>

# Usage
```bash
$ sudo ./stand-linux-arch64
```
```cmd
> stand-windows-arch64.exe
```

# Building
```bash
# Requires Go 1.25.3+ 

$ go mod tidy
$ make
```

# Flags

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| `deep` | No | `false` | Scans for NotITG by forcefully reading every program's memory |
| `pid` | No | `0` | Use a specific process |
| `verbose` | No | `false` | Enable debug messages |
| `port` | No | `8000` | The default port to run the server on |
| `version` | No | `false` | Shows the version of the binary, and exits |

# Routes

## GET `/`

Returns either an empty array, or an integer array of the entire current values of the external memory region.

## GET `/ws`

| Query | Type | Required | Description |
| --- | --- | --- | --- |
| `appid` | `int` | yes | The application id the client will be listening on |

### Receiving Messages

The websocket client will be receiving a `binary` message, with the first byte of data indicating the type of message.

| Value | Description | Notes |
| --- | --- | --- |
| `0x1` | NotITG has initialized | |
| `0x2` | NotITG has exited | |
| `0x3` | NotITG has sent a buffer to the client | The buffer content is the rest of the data after `0x3`. |

### Sending Messages

The websocket client can also send messages to the server, and it will be directly sent (and split if needed) as a buffer.

If the message type is a `text`, it will be automatically converted into a buffer.
