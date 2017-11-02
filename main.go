package main

import (
	"flag"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// parse flags
	var headers sliceFlag
	flag.Var(&headers, "h", "Add header to request")

	var url string
	flag.StringVar(&url, "url", "", "Target URL")

	var connections uint
	flag.UintVar(&connections, "c", 1, "Number of open Websocket Connections")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")

	flag.Parse()

	// build http.Headers
	h := parseHeaders(headers)

	var connectionsCount int64
	var connectionsErrCount int64

	// start logging stats
	go func() {
		t := time.Tick(5 * time.Second)
		for range t {
			log.Printf("connections=%v, errors_since_start=%v",
				atomic.LoadInt64(&connectionsCount),
				atomic.LoadInt64(&connectionsErrCount))
		}
	}()

	// initialize connections pool
	// failed connections will add a "ticket" to this pool
	// and then will restart over time.
	pool := make(chan struct{}, connections)
	for i := uint(0); i < connections; i++ {
		pool <- struct{}{}
	}

	// start connections from pool and maintain pool size
	for range pool {
		go func() {
			c, err := newWebsocketConnection(url, h)

			if err == nil {
				defer func() {
					c.Close()
					atomic.AddInt64(&connectionsCount, -1)
				}()

				atomic.AddInt64(&connectionsCount, 1)
				err = handleConnection(c)
			}

			if err != nil {
				atomic.AddInt64(&connectionsErrCount, 1)
				if verbose {
					log.Printf("err: %v", err)
				}
				pool <- struct{}{} // add ticket to pool, block if pool is full
			}
		}()

		time.Sleep(100 * time.Millisecond)
	}
}

func newWebsocketConnection(url string, h http.Header) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, h)
	return c, err
}

func handleConnection(c *websocket.Conn) error {
	// read and discard messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			return err
		}
	}

	panic("never")
}
