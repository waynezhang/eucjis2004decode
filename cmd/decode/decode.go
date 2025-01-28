package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/waynezhang/eucjis2004decode/eucjis2004"
	"golang.org/x/text/transform"
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

	dec := transform.NewReader(f, &eucjis2004.EUCJIS2004Decoder{})
	s := bufio.NewScanner(dec)
	for s.Scan() {
		text := s.Text()
		fmt.Println(text)
	}
}
