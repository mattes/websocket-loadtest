package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type LoadTest struct {
	Url     string
	Headers http.Header
}

func main() {
	var err error

	// parse flags
	var file string
	flag.StringVar(&file, "file", "", "Read test matrix from file")

	var headers sliceFlag
	flag.Var(&headers, "h", "Add header to request (key=value format)")

	var url string
	flag.StringVar(&url, "url", "", "Target URL")

	var connections uint
	flag.UintVar(&connections, "c", 1, "Number of concurrent websocket connections")

	var throttleConnections time.Duration
	flag.DurationVar(&throttleConnections, "throttle", 200*time.Millisecond, "Allow new websocket connection every n milliseconds")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	// initialize load test
	var loadtest []LoadTest

	// if -file is present, read that file
	if file != "" {
		loadtest, err = readFromFile(file)
		if err != nil {
			fatal(err)
		}

	} else if url != "" {
		// no -file, but -url is present
		loadtest = []LoadTest{
			{Url: url},
		}

		if len(headers) > 0 {
			loadtest[0].Headers, err = parseHeaders(headers)
			if err != nil {
				fatal(err)
			}
		}
	} else {
		fatal("either provide -url or -file")
	}

	var connectionsCount int64
	var connectionsErrCount uint64
	var readMessagesCount uint64

	// start logging stats
	go func() {
		t := time.Tick(5 * time.Second)
		for range t {
			log.Printf("connections=%v, read_messages=%v errors_since_start=%v",
				atomic.LoadInt64(&connectionsCount),
				atomic.LoadUint64(&readMessagesCount),
				atomic.LoadUint64(&connectionsErrCount))
		}
	}()

	// create tickets over time to throttle connection creation
	// (can only create connection, if ticket is available)
	ticket := make(chan struct{}, 10)
	go func() {
		tick := time.Tick(throttleConnections)
		for range tick {
			ticket <- struct{}{}
		}
	}()

	// create goroutine for every load test
	// and multiple that with num of connections
	var wg sync.WaitGroup
	wg.Add(int(uint(len(loadtest)) * connections))
	for _, l := range loadtest {
		for i := uint(0); i < connections; i++ {
			go func(l LoadTest) {
				defer wg.Done()

				// try forever if connection fails
				for {
					// wait for ticket to become available
					<-ticket

					// create new connection
					c, err := newWebsocketConnection(l.Url, l.Headers)
					if err != nil {
						atomic.AddUint64(&connectionsErrCount, 1)
						if verbose {
							log.Printf("err: %v", err)
						}

					} else {
						atomic.AddInt64(&connectionsCount, 1)

						// start reading from connection
						if err := handleConnection(c, &readMessagesCount); err != nil {
							atomic.AddUint64(&connectionsErrCount, 1)
							if verbose {
								log.Printf("err: %v", err)
							}
						}

						// close connection
						c.Close()
						atomic.AddInt64(&connectionsCount, -1)
					}
				}
			}(l)
		}
	}

	wg.Wait()
}

func newWebsocketConnection(url string, h http.Header) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, h)
	return c, err
}

func handleConnection(c *websocket.Conn, messageCount *uint64) error {
	// read and discard messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			return err
		}

		atomic.AddUint64(messageCount, 1)
	}
}
