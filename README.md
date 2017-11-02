# Websockets Loadtest Tool

A small helper tool to load test Websockets.

## Usage

```
$ go get github.com/templarbit/websocket-loadtest
websocket-loadtest -c 100 -url wss://example.com -verbose -h cookie=user=123 -h origin=https://app.example.com
```

## TODO

  - [ ] Write to Websocket (it's read-only right now) 
  - [ ] Read URL & HTTP headers test matrix from file

