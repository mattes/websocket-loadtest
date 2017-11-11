# Websockets Loadtest Tool [![Build Status](https://travis-ci.org/templarbit/websocket-loadtest.svg?branch=master)](https://travis-ci.org/templarbit/websocket-loadtest)

A small helper tool to load test Websockets.

## Usage

```
$ go get github.com/templarbit/websocket-loadtest
websocket-loadtest -help
websocket-loadtest -c 100 -url wss://example.com -verbose -h cookie=user=123 -h origin=https://app.example.com
websocket-loadtest -file example.txt -c 2 -verbose
```

## TODO

  - [ ] Write to Websocket (it's read-only right now) 

