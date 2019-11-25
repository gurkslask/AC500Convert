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
	var access = flag.Bool("access", false, "Access or not")
	var comli = flag.Bool("comli", false, "comli or modbus")
	flag.Parse()
	filename := *path

	text, err := AC500Convert.Openfile(filename)
	stext := strings.Split(text, "\n")
	if err != nil {
		log.Fatal(err)
	}
	if *comli {

		if *access {
			rvars, err := AC500Convert.GenerateAccessComli(stext)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(AC500Convert.OutputToText(rvars))
			for i := 0; i < len(rvars); i++ {
				fmt.Println(rvars[i])
			}
		} else {
			fmt.Println(stext)
			data, err := AC500Convert.ExtractDataComli(stext)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(data)
			fmt.Println(AC500Convert.OutputToText(data))
		}
	} else {

		if *access {
			rvars, err := AC500Convert.GenerateAccessModbus(stext)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(AC500Convert.OutputToText(rvars))
			for i := 0; i < len(rvars); i++ {
				fmt.Println(rvars[i])
			}
		} else {
			data, err := AC500Convert.ExtractDataModbus(stext)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range data {
				fmt.Println(i)
			}
		}

	}
}
