package main

import (
	"fmt"
)

type sliceFlag []string

func (f *sliceFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *sliceFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}
