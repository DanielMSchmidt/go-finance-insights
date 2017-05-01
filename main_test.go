package main

import (
	"math"
	"testing"
)

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

func TestGetContractee(t *testing.T) {
	var auftraggeberTests = []struct {
		in  string
		out string
	}{
		{`Auftraggeber: mobilcom-debitel Buchungstext: Kd 0032219297 Wir sagen danke. RG-N r. M17029539092 12,98 EUR End-to-End-Ref.: 0032219297/938113306921 CORE / Mandatsref.: MC-32219297-00000001 Gl�ubiger-ID: DE43ZZZ00000074855 Ref. JK217111A2804618/69196`, "mobilcom-debitel"},
		{`Auftraggeber: Deutsche Post AG Buchungstext: KD0020920000 2809206926 REFSAP31427 45412/EFI RE7455307009 DAT24.04.201 7 VIELEN DANK, IHRE EFILIALE 314274 5412 End-to-End-Ref.: LS-FP2-840-1000-2809206926 CORE / Mandatsref.: POSTAG0053000004218857 Gl�ubiger-ID: DE65ZZZ00000210259 Ref. HN217115G1934031/697`, "Deutsche Post AG"},
	}

	for _, tt := range auftraggeberTests {
		actual := getContractee(tt.in)
		if actual != tt.out {
			t.Errorf("TestgetContractee(%s): expected %s, actual %s", tt.in, tt.out, actual)
		}
	}
}

func TestParseCSV(t *testing.T) {
	assertet := []transaction{
		{"02.05.2017", "DE60200411550892963000", -100.00},
		{"12.05.2017", "DE76210501701001803152", -48.50},
		{"25.04.2017", "Penny Mundsburg Ce", -27.9},
		{"24.04.2017", "2708 GAMESTOP", -57.98},
	}
	demoCSV := loadFile("fixtures/demo1.csv")
	result := parseCSV(demoCSV)

	if len(result) != 4 {
		t.Error("Not enough items found")
		return
	}

	for i := 0; i < 4; i++ {
		r := result[i]
		a := assertet[i]

		if r.date != a.date {
			t.Error("Dates are not the same")
			break
		}

		if r.issuer != a.issuer {
			t.Error("IBANs are not the same")
			break
		}

		if math.Abs(r.value-a.value) > 0.001 {
			t.Error("Values are not the same")
			t.Error(r.value)
			t.Error(a.value)
			break
		}
	}
}
