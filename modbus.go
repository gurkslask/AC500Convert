package AC500Convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractDataModbus(input []string) ([]VARS, error) {

	//Regexes
	bitstr := regexp.MustCompile(`^\s*(.*) AT %RX(\d*).(\d*)\.(\d*):\s*(\w*) *;(.*)`)
	regstr := regexp.MustCompile(`^\s*(.*) AT %RW(\d*).(\d*):\s*(\w*);(.*)`)

	//Translate datatypes
	regmap := map[string]string{"UINT": "UINT16", "WORD": "UINT16", "REAL": "REAL"}
	regglobmap := map[string]string{"UINT": "FLOAT", "WORD": "FLOAT", "REAL": "FLOAT"}
	bitmap := map[string]string{"BOOL": "BOOL"}
	bitglobmap := map[string]string{"BOOL": "DEFAULT"}

	var vars []VARS
	for _, row := range input {
		var tvars VARS
		if strings.Contains(row, " AT ") {
			//Its a communication variable
			if strings.Contains(row, "%RW") {
				if regstr.MatchString(row) {
					rowdata := regstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					prefixadress, err := strconv.Atoi(rowdata[2])
					if err != nil {
						fmt.Println(err)
					}
					adress, _ := strconv.Atoi(rowdata[3])
					adress += prefixadress * 32768
					tvars.adress = fmt.Sprintf("4%05s", strconv.Itoa(adress))
					tvars.datatype = regmap[strings.ToUpper(rowdata[4])]
					tvars.globaldatatype = regglobmap[strings.ToUpper(rowdata[4])]
					tvars.comment = rowdata[5]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					vars = append(vars, tvars)
				}

			}
			if strings.Contains(row, "%RX") {
				if bitstr.MatchString(row) {
					rowdata := bitstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					prefixadress, err := strconv.Atoi(rowdata[2])
					if err != nil {
						fmt.Println(err)
					}
					iadressh, err := strconv.Atoi(rowdata[3])
					if err != nil {
						fmt.Println(err)
					}
					iadressl, err := strconv.Atoi(rowdata[4])
					if err != nil {
						fmt.Println(err)
					}
					iadress := iadressh*8 + iadressl + (prefixadress * 32768)
					tvars.adress = fmt.Sprintf("%05v", iadress)
					tvars.datatype = bitmap[strings.ToUpper(rowdata[5])]
					tvars.globaldatatype = bitglobmap[strings.ToUpper(rowdata[5])]
					tvars.comment = rowdata[6]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					vars = append(vars, tvars)
				}
			}
		}
	}
	return vars, nil
}
func GenerateAccessModbus(s []string) ([]string, error) {
	var sres []string
	var rnum int = 1
	var bnumLow int = 0
	var bnumHigh int = 0
	var regstr = regexp.MustCompile(`AT %RW1\.(\d*)\s*:`)
	var regstrreplace = regexp.MustCompile(`AT %RW1\.\d*\s*`)
	var bitstr = regexp.MustCompile(`%RX0\.(\d*)\.(\d*)\s*:`)
	var bitstrreplace = regexp.MustCompile(`AT %RX0\.\d*\.\d*\s*`)
	var err error

	for key, row := range s {
		if strings.Contains(row, "BOOL") || strings.Contains(row, "bool") {
			if strings.Contains(row, " AT ") && bitstr.MatchString(row) {
				res := bitstr.FindStringSubmatch(row)
				bnumHigh, err = strconv.Atoi(res[1])
				bnumLow, err = strconv.Atoi(res[2])
				if err != nil {
					return nil, err
				}
				s[key] = bitstrreplace.ReplaceAllString(row, "")
				row = s[key]
			}
			split := strings.Split(row, ":")
			sres = append(sres, fmt.Sprintf("%s AT %%RX0.%v.%v:%s", strings.TrimSpace(split[0]), bnumHigh, bnumLow, strings.TrimSpace(strings.ToUpper(split[1]))))
			raiseModbusCounter(1, &bnumLow, &bnumHigh)
		}
		if strings.Contains(row, "UINT") || strings.Contains(row, "uint") || strings.Contains(row, "WORD") || strings.Contains(row, "word") {
			if strings.Contains(row, " AT ") && regstr.MatchString(row) {
				res := regstr.FindStringSubmatch(row)
				rnum, err = strconv.Atoi(res[1])
				if err != nil {
					return nil, err
				}
				s[key] = regstrreplace.ReplaceAllString(row, "")
				row = s[key]
			}
			split := strings.Split(row, ":")
			sres = append(sres, fmt.Sprintf("%s AT %%RW1.%v:%s", strings.TrimSpace(split[0]), rnum, strings.TrimSpace(strings.ToUpper(split[1]))))
			rnum++
		}
		if strings.Contains(row, "BJUMP") {
			split := strings.Split(row, " ")
			jumpnum, err := strconv.Atoi(split[2])
			if err != nil {
				return nil, err
			}
			raiseModbusCounter(jumpnum, &bnumLow, &bnumHigh)
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
