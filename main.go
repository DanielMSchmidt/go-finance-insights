package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type transaction struct {
	date  string
	iban  string
	value float64
}

func loadFile(filePath string) string {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(in)
}

func getIBAN(content string) string {
	startOffset := 6
	start := strings.Index(content, "IBAN:")
	end := strings.Index(content, " BLZ")
	if start < 0 || end < 0 {
		return ""
	}

	return content[start+startOffset : end]
}

func transformEntry(line string) (transaction, error) {
	parts := strings.Split(line, ";")

	if len(parts) < 5 {
		fmt.Println("This entry looks wrong")
		fmt.Println(line)
		return transaction{}, fmt.Errorf("Entry is not well formed")
	}

	if len(parts[0]) < 3 {
		return transaction{}, fmt.Errorf("Date is too short")
	}

	if len(parts[4]) < 3 {
		return transaction{}, fmt.Errorf("Amount is too short")
	}

	date := parts[0][1 : len(parts[0])-1]
	iban := getIBAN(parts[3])
	rawAmount := parts[4][1 : len(parts[4])-1]
	amount, err := strconv.ParseFloat(strings.Replace(rawAmount, ",", ".", -1), 32)

	if err != nil {
		return transaction{}, fmt.Errorf("Amount is not a float")
	}

	return transaction{date, iban, amount}, nil
}

func parseCSV(content string) []transaction {
	startOffset := 16
	endOffset := 4
	start := strings.Index(content, "Umsatz in EUR")
	end := strings.Index(content, "Alter Kontostand")
	if start < 0 || end < 0 {
		panic("Invalid CSV, no start or end delimiter found")
	}
	rawEntries := content[start+startOffset : end-endOffset]
	entries := strings.Split(rawEntries, "\n")
	var result []transaction

	for i := 0; i < len(entries); i++ {
		entry, err := transformEntry(entries[i])
		if err == nil {
			result = append(result, entry)
		}
	}

	return result
}

func main() {
	raw := loadFile("data/foo.csv")
	transactions := parseCSV(raw)
	for i := 0; i < len(transactions); i++ {
		fmt.Println(transactions[i])
	}
}
