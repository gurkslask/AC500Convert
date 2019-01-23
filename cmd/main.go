package main

import (
	"flag"
	"fmt"
	"github.com/gurkslask/AC500Convert"
	"log"
	"strings"
)

func main() {
	var path = flag.String("path", "./data.txt", "path to file")
	flag.Parse()
	filename := *path

	text, err := AC500Convert.Openfile(filename)
	stext := strings.Split(text, "\n")
	if err != nil {
		log.Fatal(err)
	}

	vars, err := AC500Convert.ExtractData(stext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(AC500Convert.OutputToText(vars))

}
