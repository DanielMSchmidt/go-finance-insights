package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type transaction struct {
	date   string
	issuer string
	value  float64
}

func loadFile(filePath string) string {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(in)
}

func getContractee(content string) string {
	startOffset := 14
	start := strings.Index(content, "Auftraggeber: ")
	end := strings.Index(content, " Buchungstext")

	if start < 0 || end < 0 {
		return ""
	}

	return content[start+startOffset : end]
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
	var issuer string
	if iban == "" {
		issuer = getContractee(parts[3]) // lastschrift
	} else {
		issuer = iban
	}

	rawAmount := parts[4][1 : len(parts[4])-1]
	amount, err := strconv.ParseFloat(strings.Replace(rawAmount, ",", ".", -1), 32)

	if err != nil {
		return transaction{}, fmt.Errorf("Amount is not a float")
	}

	return transaction{date, issuer, amount}, nil
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

func floatEquals(a float64, b float64) bool {
	return math.Abs(a-b) < 0.001
}

func findEqualAmounts(entries []transaction) [][]transaction {
	buckets := [][]transaction{}
	itemAdded := false
	for _, entry := range entries {
		for bucketIndex, bucket := range buckets {
			if floatEquals(entry.value, bucket[0].value) {
				buckets[bucketIndex] = append(bucket, entry)
				itemAdded = true
			}
		}
		if !itemAdded {
			list := []transaction{}
			list = append(list, entry)
			buckets = append(buckets, list)
		}

	}

	return buckets
}

func filterMatches(entries [][]transaction) (ret [][]transaction) {
	for _, entry := range entries {
		if len(entry) > 1 {
			ret = append(ret, entry)
		}
	}
	return

}

func main() {
	raw := loadFile("data/foo.csv")
	buckets := filterMatches(findEqualAmounts(parseCSV(raw)))
	for _, bucket := range buckets {
		fmt.Println("\n{")
		for _, trans := range bucket {
			fmt.Println(trans.date, trans.issuer, trans.value)
		}
		fmt.Println("}")
	}

}
