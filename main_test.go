package main

import "testing"

func TestGetIBAN(t *testing.T) {
	var ibanTests = []struct {
		in  string // input
		out string // expected result
	}{
		{`Empf�nger: Daniel Schmidt, Carina Heitmann Kto/IBAN: DE60200411550892963000 BLZ/BIC: COBADEHD055 Buchungstext: End-to-End-Ref.: nicht angegeben Ref. JF217118L1327255/2";`, "DE60200411550892963000"},
		{`Empf�nger: Landeshauptstadt Kiel Kto/IBAN: DE76210501701001803152 BLZ/BIC: NOLADE21KIE Buchungstext: 1057005241701 End-to-End-Ref.: nicht angegeben Ref. JF217118L1327288/2;`, "DE76210501701001803152"},
	}

	for _, tt := range ibanTests {
		actual := getIBAN(tt.in)
		if actual != tt.out {
			t.Errorf("TestGetIBAN(%s): expected %s, actual %s", tt.in, tt.out, actual)
		}
	}
}

func TestParseCSV(t *testing.T) {
	assertet := []transaction{
		{"02.05.2017", "DE60200411550892963000", -100.00},
		{"12.05.2017", "DE76210501701001803152", -48.50},
	}
	demoCSV := loadFile("fixtures/demo1.csv")
	result := parseCSV(demoCSV)

	if len(result) != 2 {
		t.Error("Not enough items found")
		return
	}

	for i := 0; i < 2; i++ {
		r := result[i]
		a := assertet[i]
		if r.date != a.date {
			t.Error("Names are not the same")
		}

		if r.iban != a.iban {
			t.Error("IBANs are not the same")
		}

		if r.value != a.value {
			t.Error("values are not the same")
		}
	}
}
