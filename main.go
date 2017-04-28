package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	in, err := ioutil.ReadFile("data/foo.csv")
	if err != nil {
		log.Fatal(err)
	}

	raw = string(in)

	fmt.Println(raw)

}
