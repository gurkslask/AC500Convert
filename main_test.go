package AC500Convert

import (
	"fmt"
	"log"
	"testing"
)

func TestExtractComli(t *testing.T) {
	istr := []string{"var1 AT %RX0.1.0:BOOL;(*kommentar*)",
		"var2 AT %RX0.2.0:BOOL;(*kommentar2*)",
		"var4 AT %RX0.42.0:BOOL ;(*kommentar2*)",
		"var3 AT %RW1.3:UINT;           (*uint*)",
	}

	want := []VARS{
		VARS{tag: "var1", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00001", comment: "kommentar"},
		VARS{tag: "var2", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00002", comment: "kommentar2"},
		VARS{tag: "var4", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00052", comment: "kommentar2"},
		VARS{tag: "var3", datatype: "UINT16", globaldatatype: "FLOAT", adress: "R00003", comment: "uint"},
	}

	got, err := ExtractDataComli(istr)
	if err != nil {
		fmt.Println(err)
	}

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant:%v\n", got[key], want[key])
		}

	}
}
func TestGenerateAccessComli(t *testing.T) {
	istr := []string{"var1 :BOOL;(*kommentar*)",
		"(* BJUMP 50 *)",
		"var2:bool;(*kommentar2*)",
		"var3:uint;           (*uint*)",
		"(* RJUMP 10 *)",
		"var4:uint;           (*uint*)",
		"var5 : uint ;           (*uint*)",
	}
	got, err := GenerateAccessComli(istr)

	if err != nil {
		log.Fatal(err)
	}
	want := []string{
		"var1 AT %RX0.1.0:BOOL;(*KOMMENTAR*)",
		"var2 AT %RX0.52.0:BOOL;(*KOMMENTAR2*)",
		"var3 AT %RW1.1:UINT;           (*UINT*)",
		"var4 AT %RW1.12:UINT;           (*UINT*)",
		"var5 AT %RW1.13:UINT ;           (*UINT*)",
	}

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant: %v\n", got[key], want[key])
		}

	}
}

func TestGenerateAccessModbus(t *testing.T) {
	istr := []string{
		"var1 :BOOL;(*kommentar*)",
		"(* BJUMP 50 *)",
		"var2:bool;(*kommentar2*)",
		"var3:uint;           (*uint*)",
		"(* RJUMP 10 *)",
		"var4:uint;           (*uint*)",
		"var5 : uint ;           (*uint*)",
		"var5 : uint ;           (*uint*)",
	}
	got, err := GenerateAccessModbus(istr)

	if err != nil {
		log.Fatal(err)
	}
	want := []string{
		"var1 AT %RX0.0.0:BOOL;(*KOMMENTAR*)",
		"var2 AT %RX0.6.3:BOOL;(*KOMMENTAR2*)",
		"var3 AT %RW1.1:UINT;           (*UINT*)",
		"var4 AT %RW1.12:UINT;           (*UINT*)",
	}

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant: %v\n", got[key], want[key])
		}

	}
}
func TestExtractModbus(t *testing.T) {
	istr := []string{"var1 AT %RX0.0.1:BOOL;(*kommentar*)",
		"var2 AT %RX0.2.0:BOOL;(*kommentar2*)",
		"var4 AT %RX0.42.4:BOOL ;(*kommentar2*)",
		"var5 AT %RX1.42.4:BOOL ;(*kommentar2*)",
		"var3 AT %RW1.3:UINT;           (*uint*)",
	}

	want := []VARS{
		VARS{tag: "var1", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00001", comment: "kommentar"},
		VARS{tag: "var2", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00016", comment: "kommentar2"},
		VARS{tag: "var4", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00340", comment: "kommentar2"},
		VARS{tag: "var5", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "33108", comment: "kommentar2"},
		VARS{tag: "var3", datatype: "UINT16", globaldatatype: "FLOAT", adress: "432771", comment: "uint"},
	}

	got, err := ExtractDataModbus(istr)
	if err != nil {
		fmt.Println(err)
	}

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant:%v\n", got[key], want[key])
		}

	}
}
