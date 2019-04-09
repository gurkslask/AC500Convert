package AC500Convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractDataComli(input []string) ([]VARS, error) {
	// Extracts data from AC500 accessvariables

	//Regexes
	bitstr := regexp.MustCompile(`^\s*(.*) AT %RX0.(\d*)\.0:\s*(\w*) *;(.*)`)
	regstr := regexp.MustCompile(`^\s*(.*) AT %RW1.(\d*):\s*(\w*);(.*)`)

	//Translate datatypes
	regmap := map[string]string{"UINT": "UINT16", "WORD": "UINT16"}
	regglobmap := map[string]string{"UINT": "FLOAT", "WORD": "FLOAT"}
	bitmap := map[string]string{"BOOL": "BOOL"}
	bitglobmap := map[string]string{"BOOL": "DEFAULT"}

	var vars []VARS
	for _, row := range input {
		var tvars VARS
		if strings.Contains(row, " AT ") {
			//Its a communication variable
			if strings.Contains(row, "%RW") {
				//fmt.Println("reg")
				if regstr.MatchString(row) {
					//fmt.Println(row)
					rowdata := regstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					//iadress, err := strconv.Atoi(rowdata[2])
					tvars.adress = fmt.Sprintf("R%05s", rowdata[2])
					//fmt.Print(tvars.adress)
					tvars.datatype = regmap[strings.ToUpper(rowdata[3])]
					tvars.globaldatatype = regglobmap[strings.ToUpper(rowdata[3])]
					tvars.comment = rowdata[4]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					//fmt.Println(tvars)
					vars = append(vars, tvars)
				}

			}
			if strings.Contains(row, "%RX") {
				//fmt.Println("This is a bit")
				if bitstr.MatchString(row) {
					//fmt.Println(row)
					rowdata := bitstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					iadress, err := strconv.Atoi(rowdata[2])
					if err != nil {
						fmt.Println(err)
					}
					tvars.adress = fmt.Sprintf("%05o", iadress)
					tvars.datatype = bitmap[strings.ToUpper(rowdata[3])]
					tvars.globaldatatype = bitglobmap[strings.ToUpper(rowdata[3])]
					tvars.comment = rowdata[4]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					vars = append(vars, tvars)
				}
			}
		}
	}
	//fmt.Println(vars)
	return vars, nil
}
func GenerateAccessComli(s []string) ([]string, error) {
	//var res string
	var sres []string
	var rnum int = 1
	var bnum int = 1
	//fmt.Println("BEGIN\n", s)

	for _, row := range s {
		//fmt.Println(sres + "\n ENDS HERE")
		if strings.Contains(row, "BOOL") || strings.Contains(row, "bool") {
			//fmt.Println("bool")
			split := strings.Split(row, ":")
			//res += fmt.Sprintf("%s AT %%RX0.%v.0:%s\r", split[0], bnum, strings.ToUpper(split[1]))
			sres = append(sres, fmt.Sprintf("%s AT %%RX0.%v.0:%s", strings.TrimSpace(split[0]), bnum, strings.TrimSpace(strings.ToUpper(split[1]))))
			bnum++
		}
		if strings.Contains(row, "UINT") || strings.Contains(row, "uint") || strings.Contains(row, "WORD") || strings.Contains(row, "word") {
			//fmt.Println("reg")
			split := strings.Split(row, ":")
			//res += fmt.Sprintf("%s AT %%RW1.%v:%s\r", split[0], rnum, strings.ToUpper(split[1]))
			sres = append(sres, fmt.Sprintf("%s AT %%RW1.%v:%s", strings.TrimSpace(split[0]), rnum, strings.TrimSpace(strings.ToUpper(split[1]))))
			rnum++
		}
		if strings.Contains(row, "BJUMP") {
			split := strings.Split(row, " ")
			jumpnum, err := strconv.Atoi(split[2])
			if err != nil {
				return nil, err
			}
			bnum += jumpnum
		}
		if strings.Contains(row, "RJUMP") {
			split := strings.Split(row, " ")
			jumpnum, err := strconv.Atoi(split[2])
			if err != nil {
				return nil, err
			}
			rnum += jumpnum
		}
	}
	return sres, nil

}
