package main

import (
	"github.com/Jaezmien/notitg-external-go"
)

type HeartbeatStatus uint

const (
	HEARTBEAT_EXITED HeartbeatStatus = iota
	HEARTBEAT_NOTFOUND
	HEARTBEAT_FOUND
)

const (
	INCOMING_ID         = 0
	INCOMING_TYPE       = 1
	INCOMING_LENGTH     = 2
	INCOMING_DATA_START = 3
	INCOMING_DATA_END   = 29
	INCOMING_STATE      = 30

	OUTGOING_ID         = 31
	OUTGOING_TYPE       = 32
	OUTGOING_LENGTH     = 33
	OUTGOING_DATA_START = 34
	OUTGOING_DATA_END   = 60
	OUTGOING_STATE      = 61

	INIT_STATE = 63
)

const (
	STATE_INCOMING_IDLE      = 0
	STATE_INCOMING_BUSY      = 1
	STATE_INCOMING_AVAILABLE = 2

	STATE_OUTGOING_IDLE      = 0
	STATE_OUTGOING_AVAILABLE = 1
)

const (
	INIT_EXPECTED_VALUE = 56
)

var VersionsAsBytes = map[notitg.NotITGVersion][]byte{
	notitg.V1:     []byte{0x1, 0x0, 0x0, 0x0},
	notitg.V2:     []byte{0x2, 0x0, 0x0, 0x0},
	notitg.V3:     []byte{0x3, 0x0, 0x0, 0x0},
	notitg.V3_1:   []byte{0x3, 0x1, 0x0, 0x0},
	notitg.V4:     []byte{0x4, 0x0, 0x0, 0x0},
	notitg.V4_0_1: []byte{0x4, 0x0, 0x1, 0x0},
	notitg.V4_2:   []byte{0x4, 0x2, 0x0, 0x0},
	notitg.V4_3:   []byte{0x4, 0x3, 0x0, 0x0},
	notitg.V4_9:   []byte{0x4, 0x9, 0x0, 0x0},
	notitg.V4_9_1: []byte{0x4, 0x9, 0x1, 0x0},
}
