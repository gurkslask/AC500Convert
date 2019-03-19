package AC500Convert

import (
	"fmt"
	"log"
	"testing"
)

func TestExtract(t *testing.T) {
	istr := []string{"var1 AT %RX0.1.0:BOOL;(*kommentar*)",
		"var2 AT %RX0.2.0:BOOL;(*kommentar2*)",
		"var4 AT %RX0.42.0:BOOL;(*kommentar2*)",
		"var3 AT %RW1.3:UINT;           (*uint*)",
	}

	want := []VARS{
		VARS{tag: "var1", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00001", comment: "kommentar"},
		VARS{tag: "var2", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00002", comment: "kommentar2"},
		VARS{tag: "var4", datatype: "BOOL", globaldatatype: "DEFAULT", adress: "00052", comment: "kommentar2"},
		VARS{tag: "var3", datatype: "UINT16", globaldatatype: "FLOAT", adress: "R00003", comment: "uint"},
	}

	got, err := ExtractData(istr)
	if err != nil {
		fmt.Println(err)
	}

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant:%v\n", got[key], want[key])
		}

	}
}
func TestGenerateAccess(t *testing.T) {
	istr := []string{"var1 :BOOL;(*kommentar*)",
		"(* BJUMP 50 *)",
		"var2:bool;(*kommentar2*)",
		"var3:uint;           (*uint*)",
		"(* RJUMP 10 *)",
		"var4:uint;           (*uint*)",
		"var5 : uint ;           (*uint*)",
	}
	got, err := GenerateAccess(istr)

	if err != nil {
		log.Fatal(err)
	}
	want := []string{
		"var1  AT %RX0.1.0:BOOL;(*KOMMENTAR*)",
		"var2 AT %RX0.52.0:BOOL;(*KOMMENTAR2*)",
		"var3 AT %RW1.1:UINT;           (*UINT*)",
		"var4 AT %RW1.12:UINT;           (*UINT*)",
		"var5  AT %RW1.13: UINT ;           (*UINT*)",
	}
	//want := "var1  AT %RX0.1.0:BOOL;(*KOMMENTAR*)\nvar2 AT %RX0.52.0:BOOL;(*KOMMENTAR2*)\nvar3 AT %RW1.1:UINT;           (*UINT*)\nvar4 AT %RW1.12:UINT;           (*UINT*)\n"

	for key, _ := range want {
		if got[key] != want[key] {
			t.Fatalf("Got: %v\nWant: %v\n", got[key], want[key])
		}

	}
	/*if got != want {
		t.Fatalf("Got:\n%v\nWant:\n%v\n", got, want)
	}*/
}
