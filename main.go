package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"github.com/Jaezmien/notitg-external-go"
	"github.com/gorilla/websocket"
)

var DeepScan = false
var ProcessID = 0
var Verbose = false
var Port = 8080

func init() {
	flag.BoolVar(&DeepScan, "deep", false, "Scan deeply by checking each process' memory")
	flag.IntVar(&ProcessID, "pid", 0, "Use a specific process")
	flag.BoolVar(&Verbose, "verbose", false, "Enable debug messages")
	flag.IntVar(&Port, "port", 8000, "Sets the server port")

	flag.Parse()
}

func main() {
	done := make(chan bool, 1)

	s := NewLemonadeStand(
		WithDeepScan(DeepScan),
		WithProcessID(ProcessID),
		WithLogger(Verbose),
		WithTickRate(10),
	)

	s.OnRead = func(l *LemonadeStand, appid int32, buffer []int32) {
		// TODO:
	}
	s.OnExit = func(l *LemonadeStand) {
		// TODO:
	}
	s.OnClose = func(l *LemonadeStand) {
		done <- true
	}

	go s.Run()

	upgrader := websocket.Upgrader{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(400)
			fmt.Fprintf(w, "unknown method")
			return
		}

		q, _ := url.ParseQuery(r.URL.RawQuery)
		i, err := strconv.ParseInt(q.Get("appid"), 10, 32)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "unknown method")
			return
		}

		_ = int32(i)

		_, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Error("error in upgrading connection", slog.Any("err", err))

			w.WriteHeader(400)
			fmt.Fprintf(w, "internal error")
			return
		}

		// TODO:
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(400)
			fmt.Fprintf(w, "unknown method")
			return
		}

		if !s.HasNotITG() {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "[]")
			return
		}

		data := make([]int32, 0)

		m := 64
		if s.NotITG.Version >= notitg.V4_2 {
			m = 256
		}

		for i := range m {
			data = append(data, s.NotITG.GetExternal(int(i)))
		}

		j, err := json.Marshal(data)
		if err != nil {
			s.logger.Error("marshal error:", slog.Any("error", err))

			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error")
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	})
	go func(l *LemonadeStand) {
		err := http.ListenAndServe(fmt.Sprintf("localhost:%d", Port), nil)
		if err != nil {
			l.logger.Error("http error", slog.Any("error", err))
			l.Close()
			done <- true
		}
	}(s)

	s.logger.Info("lemonade stand is ready")

	termChannel := make(chan os.Signal, 2)
	signal.Notify(termChannel, os.Interrupt)
	go func() {
		<-termChannel
		s.Close()
		done <- true
	}()
	<-done

	s.logger.Info("exiting")
}
