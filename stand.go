package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"git.jaezmien.com/Jaezmien/lemonade-stand/buffer"
	"github.com/Jaezmien/notitg-external-go"
)

var ErrNotDetected = errors.New("notitg not detected")
var ErrUnsupported = errors.New("unsupported version")

var TickerScanning = time.Second * 2

type LemonadeStand struct {
	NotITG    *notitg.NotITG
	OnConnect func(l *LemonadeStand)
	OnRead    func(l *LemonadeStand, appid int32, buffer []int32)
	OnExit    func(l *LemonadeStand)
	OnClose   func(l *LemonadeStand)

	tickRateMs time.Duration
	processID  int
	deepScan   bool

	heartbeatStatus HeartbeatStatus
	initialized     bool

	exit      chan struct{}
	isRunning bool

	readManager  *buffer.LemonadeBufferManager
	writeManager *buffer.LemonadeBufferManager

	logger        *slog.Logger
	loggerEnabled bool
}

type LemonadeStandOption func(s *LemonadeStand)

func WithProcessID(pid int) LemonadeStandOption {
	return func(s *LemonadeStand) {
		s.processID = pid
	}
}
func WithDeepScan(deep bool) LemonadeStandOption {
	return func(s *LemonadeStand) {
		s.deepScan = deep
	}
}
func WithTickRate(ms int) LemonadeStandOption {
	return func(s *LemonadeStand) {
		s.tickRateMs = time.Duration(ms)
	}
}
func WithLogger(enabled bool) LemonadeStandOption {
	return func(s *LemonadeStand) {
		s.loggerEnabled = enabled
	}
}

func NewLemonadeStand(options ...LemonadeStandOption) *LemonadeStand {
	l := &LemonadeStand{
		heartbeatStatus: HEARTBEAT_EXITED,
		initialized:     false,

		tickRateMs: 10,
		deepScan:   false,
		processID:  0,
		exit:       make(chan struct{}),
		isRunning:  false,

		readManager:  buffer.NewManager(),
		writeManager: buffer.NewManager(),
	}

	for _, opt := range options {
		opt(l)
	}

	slogOptions := &slog.HandlerOptions{}
	slogLevel := new(slog.LevelVar)
	if l.loggerEnabled {
		slogLevel.Set(slog.LevelDebug)
	} else {
		slogLevel.Set(slog.LevelInfo)
	}
	slogOptions.Level = slogLevel

	l.logger = slog.New(slog.NewTextHandler(os.Stdout, slogOptions))

	return l
}
func (l *LemonadeStand) Close() {
	if l.exit == nil {
		return
	}

	close(l.exit)
	l.exit = nil

	if l.OnClose != nil {
		l.OnClose(l)
	}
}

func (l *LemonadeStand) HasNotITG() bool {
	if l.NotITG == nil {
		return false
	}
	return l.NotITG.Heartbeat()
}
func (l *LemonadeStand) HasInitialized() bool {
	return l.initialized
}

func (l *LemonadeStand) isInitializing() bool {
	return l.NotITG.GetExternal(INIT_STATE) != INIT_EXPECTED_VALUE
}

func (l *LemonadeStand) isScanningByPID() bool {
	return l.processID != 0
}
func (l *LemonadeStand) scan() error {
	var n *notitg.NotITG
	var err error

	if l.isScanningByPID() {
		n, err = notitg.ScanProcessID(l.processID)
	} else {
		n, err = notitg.Scan(l.deepScan)
	}

	if err != nil {
		return err
	}

	if n == nil {
		l.heartbeatStatus = HEARTBEAT_NOTFOUND
		return ErrNotDetected
	}
	if n.Version <= notitg.V2 {
		l.heartbeatStatus = HEARTBEAT_NOTFOUND
		return ErrUnsupported
	}

	l.NotITG = n
	l.heartbeatStatus = HEARTBEAT_FOUND

	return nil
}

func (l *LemonadeStand) isOutgoingAvailable() bool {
	return l.NotITG.GetExternal(OUTGOING_STATE) == STATE_OUTGOING_AVAILABLE
}

// NotITG -> Lemonade Stand
func (l *LemonadeStand) read() {
	if !l.isOutgoingAvailable() {
		return
	}

	appid := l.NotITG.GetExternal(OUTGOING_ID)
	readLength := int(l.NotITG.GetExternal(OUTGOING_LENGTH))
	readBuffer := make([]int32, readLength)

	for idx := range readLength {
		flagIdx := OUTGOING_DATA_START + idx
		readBuffer[idx] = l.NotITG.GetExternal(flagIdx)
		l.NotITG.SetExternal(flagIdx, 0)
	}

	isBufferEnd := l.NotITG.GetExternal(OUTGOING_TYPE) == int32(buffer.BUFFER_END)

	l.NotITG.SetExternal(OUTGOING_LENGTH, 0)
	l.NotITG.SetExternal(OUTGOING_TYPE, 0)
	l.NotITG.SetExternal(OUTGOING_ID, 0)
	l.NotITG.SetExternal(OUTGOING_STATE, STATE_OUTGOING_IDLE)

	buff := l.readManager.NewBuffer(appid)

	if isBufferEnd {
		data := buff.AppendBuffer(readBuffer)

		if l.OnRead != nil {
			l.OnRead(l, appid, readBuffer)
		}

		l.logger.Debug(
			"read buffer from notitg",
			slog.String("data", fmt.Sprintf("%+v", data)),
			slog.Int("length", len(data)),
			slog.Int("appid", int(appid)),
		)

		l.readManager.CloseBuffer(appid)
	} else {
		buff.AppendBuffer(readBuffer)
	}
}

func (l *LemonadeStand) hasWriteBuffer() bool {
	return l.writeManager.Count() > 0
}
func (l *LemonadeStand) isIncomingAvailable() bool {
	return l.NotITG.GetExternal(INCOMING_STATE) == STATE_INCOMING_IDLE
}

// Lemonade Stand -> NotITG
func (l *LemonadeStand) write() {
	if !l.hasWriteBuffer() {
		return
	}
	if !l.isIncomingAvailable() {
		return
	}

	l.NotITG.SetExternal(INCOMING_STATE, STATE_INCOMING_BUSY)

	appid, err := l.writeManager.GetFirstID()
	if err != nil {
		panic(err)
	}

	buffer, err := l.writeManager.TryGetBuffer(appid)
	if err != nil {
		panic(err)
	}

	for idx, value := range buffer.Buffer {
		flagIdx := INCOMING_DATA_START + idx
		l.NotITG.SetExternal(flagIdx, value)
	}
	l.NotITG.SetExternal(INCOMING_LENGTH, int32(len(buffer.Buffer)))

	l.NotITG.SetExternal(INCOMING_TYPE, int32(buffer.Set))
	l.NotITG.SetExternal(INCOMING_ID, appid)
	l.NotITG.SetExternal(INCOMING_STATE, STATE_INCOMING_AVAILABLE)

	l.writeManager.CloseBuffer(appid)

	l.logger.Debug(
		"written buffer to notitg",
		slog.String("data", fmt.Sprintf("%+v", buffer.Buffer)),
		slog.Int("length", len(buffer.Buffer)),
		slog.Int("appid", int(appid)),
	)
}

func (l *LemonadeStand) Run() {
	if l.isRunning {
		return
	}

	l.isRunning = true
	defer func(s *LemonadeStand) {
		s.isRunning = false
	}(l)

	ticker := time.NewTicker(TickerScanning)

	for {
		select {
		case <-ticker.C:
			heartbeat := l.HasNotITG()
			if !heartbeat {
				if l.heartbeatStatus != HEARTBEAT_FOUND {
					err := l.scan()
					if err != nil {
						if errors.Is(err, ErrNotDetected) {
							l.logger.Info("notitg not detected, retrying in 2 seconds")
							continue
						}
						if errors.Is(err, ErrUnsupported) {
							l.logger.Info("unsupported version of notitg (found v2, or below)")
							continue
						}

						l.logger.Error("found error while scanning for notitg", slog.Any("error", err))
					}

					l.logger.Info("notitg found")

					if l.OnConnect != nil {
						l.OnConnect(l)
					}
				} else if l.heartbeatStatus == HEARTBEAT_EXITED {
					l.heartbeatStatus = HEARTBEAT_NOTFOUND
					l.initialized = false

					continue
				} else if l.heartbeatStatus == HEARTBEAT_FOUND {
					l.heartbeatStatus = HEARTBEAT_EXITED
					l.initialized = false
					l.NotITG = nil
					ticker.Reset(TickerScanning)

					if l.isScanningByPID() {
						l.logger.Warn("notitg scanning by pid, but process has exited.")
						l.logger.Warn("closing.")
						l.Close()
					} else {
						l.logger.Warn("notitg has closed.")

						if l.OnExit != nil {
							l.OnExit(l)
						}
					}

					continue
				}

			}

			if !l.initialized {
				if l.isInitializing() {
					l.logger.Info("notitg currently initializing")
					continue
				}
				l.logger.Info("notitg has initialized")

				l.initialized = true
				ticker.Reset(time.Millisecond * l.tickRateMs)
			}

			l.read()
			l.write()
		case <-l.exit:
			ticker.Stop()
			return
		}
	}
}
