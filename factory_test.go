package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGewusstUpdate(t *testing.T) {
	var cgw = updateGewusst()
	fmt.Println(cgw)
	var buf = bytes.NewBufferString("")
	testOutput(buf)
	for _, v := range messages {
		fmt.Fprintf(buf, "%s", v.String())
	}
	io.Copy(os.Stdout, buf)
}
