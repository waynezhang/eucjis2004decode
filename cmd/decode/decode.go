package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/waynezhang/eucjis2004decode/decode"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage:\t%s file_name\n", filepath.Base(args[0]))
		os.Exit(0)
	}

	f, err := os.Open(args[1])
	if err != nil {
		log.Fatal("cannot open file", err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)

	s := bufio.NewScanner(f)
	for s.Scan() {
		buf.Reset()

		err := decode.Convert(s.Bytes(), buf)
		if err != nil {
			fmt.Println("invalid line", s.Text())
		}
		fmt.Println(buf.String())
	}
}
