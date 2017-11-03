package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	// create temp test file
	f, err := ioutil.TempFile("", "websocket-loadtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`

		wss://example.com cookie=user=123 origin=https://app.example.com
		wss://example2.com 
		# comment 

	`)

	// test with contents from test file
	l, err := readFromFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 2 {
		t.Fatalf("expected 2, got %v", len(l))
	}

	if l[0].Url != "wss://example.com" {
		t.Errorf("expected wss://example.com, got %v", l[0].Url)
	}

	if len(l[0].Headers) != 2 {
		t.Errorf("expected 2, got %v", len(l[0].Headers))
	}

	if l[0].Headers.Get("cookie") != "user=123" {
		t.Errorf("expected user=123, got %v", l[0].Headers.Get("cookie"))
	}

	if l[0].Headers.Get("origin") != "https://app.example.com" {
		t.Errorf("expected https://app.example.com, got %v", l[0].Headers.Get("origin"))
	}
}
