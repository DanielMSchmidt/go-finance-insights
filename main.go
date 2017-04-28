package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type transaction struct {
	date  string
	iban  string
	value float32
}

func loadFile(filePath string) string {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(in)
}

func getIBAN(content string) string {
	return "TODO"
}

func parseCSV(content string) []transaction {
	startOffset := 17
	endOffset := 4
	start := strings.Index(content, "Umsatz in EUR")
	end := strings.Index(content, "Alter Kontostand")
	rawEntries := content[start+startOffset : end-endOffset]
	entries := strings.Split(rawEntries, "\n")

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		entryParts := strings.Split(entry, ";")

		date := entryParts[0]
		iban := getIBAN(entryParts[3])
		amount := entryParts[4]
		fmt.Println("\n{")
		fmt.Println(date)
		fmt.Println(iban)
		fmt.Println(amount)
		fmt.Println("}\n")
	}

	return nil
}

func main() {
	raw := loadFile("data/foo.csv")
	parseCSV(raw)
	// fmt.Println(raw)

}
