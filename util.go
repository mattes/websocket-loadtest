package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func fatal(format string, a ...interface{}) {
	format = format + "\n"
	fmt.Fprintf(os.Stderr, format, a)
	os.Exit(1)
}

func parseHeaders(in []string) http.Header {
	h := make(http.Header)
	for _, fh := range in {
		s := strings.SplitN(fh, "=", 2)
		if len(s) != 2 {
			fatal("must be key=value format: %s", fh)
		}
		h.Add(s[0], s[1])
	}
	return h
}
