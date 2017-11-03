package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

type ErrMalformed struct {
	Text string
}

func (e ErrMalformed) Error() string {
	return fmt.Sprintf("malformed: %v", e.Text)
}

func parseHeaders(in []string) (http.Header, error) {
	h := make(http.Header)
	for _, fh := range in {
		s := strings.SplitN(fh, "=", 2)
		if len(s) != 2 {
			return nil, ErrMalformed{fmt.Sprintf("must be key=value format: %s", fh)}
		}
		h.Add(s[0], s[1])
	}
	return h, nil
}

type ErrFileMalformed struct {
	Line int
	Text string
}

func (e ErrFileMalformed) Error() string {
	return fmt.Sprintf("%v in line %v", e.Text, e.Line)
}

// wss://example.com cookie=user=123 origin=https://app.example.com
func readFromFile(path string) ([]LoadTest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	l := make([]LoadTest, 0)
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines
		// skip lines that start with #
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		p := strings.SplitN(line, " ", 2)

		if len(p) == 2 {
			h, err := parseHeaders(strings.Split(p[1], " "))
			if err != nil {
				return nil, ErrFileMalformed{i, err.Error()}
			}

			l = append(l, LoadTest{
				Url:     p[0],
				Headers: h,
			})

		} else if len(p) == 1 {
			l = append(l, LoadTest{
				Url: p[0],
			})

		} else {
			return nil, ErrFileMalformed{i, line}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return l, nil
}
