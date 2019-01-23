package AC500Convert

import "fmt"
import "io/ioutil"
import "regexp"
import "strconv"
import "strings"

func Openfile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ExtractData(input []string) ([]VARS, error) {
	// Extracts data from AC500 accessvariables

	//Regexes
	bitstr := regexp.MustCompile(`^\s*(.*) AT %RX0.(\d*)\.0:(\w*);(.*)`)
	regstr := regexp.MustCompile(`^\s*(.*) AT %RW1.(\d*):(\w*);(.*)`)

	//Translate datatypes
	regmap := map[string]string{"UINT": "UINT16", "WORD": "UINT16"}
	regglobmap := map[string]string{"UINT": "FLOAT", "WORD": "FLOAT"}
	bitmap := map[string]string{"BOOL": "BOOL"}
	bitglobmap := map[string]string{"BOOL": "DEFAULT"}

	var vars []VARS
	for _, row := range input {
		if strings.Contains(row, " AT ") {
			var tvars VARS
			//Its a communication variable
			if strings.Contains(row, "%RW") {
				if regstr.MatchString(row) {
					rowdata := regstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					iadress, err := strconv.Atoi(rowdata[2])
					if err != nil {
						return nil, err

					}
					tvars.adress = iadress
					tvars.datatype = regmap[strings.ToUpper(rowdata[3])]
					tvars.globaldatatype = regglobmap[strings.ToUpper(rowdata[3])]
					tvars.comment = rowdata[4]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					vars = append(vars, tvars)
				}

			}
			if strings.Contains(row, "%RX") {
				//fmt.Println("This is a bit")
				if bitstr.MatchString(row) {
					rowdata := bitstr.FindStringSubmatch(row)
					tvars.tag = rowdata[1]
					iadress, err := strconv.Atoi(rowdata[2])
					if err != nil {
						return nil, err

					}
					tvars.adress = iadress
					tvars.datatype = bitmap[strings.ToUpper(rowdata[3])]
					tvars.globaldatatype = bitglobmap[strings.ToUpper(rowdata[3])]
					tvars.comment = rowdata[4]
					tvars.comment = RmLeadSpace(RemoveStars(tvars.comment))
					fmt.Println(tvars)
					vars = append(vars, tvars)
				}
			}
		}
	}
	return vars, nil
}

func OutputToText(vars []VARS) string {
	s := "//\nName,DataType,GlobalDataType,Adress_1,Description //"
	for _, v := range vars {
		s += v.String()
	}
	return s
}

func RemoveStars(s string) string {
	s = strings.Replace(s, "*", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return s
}

type VARS struct {
	tag            string
	datatype       string
	globaldatatype string
	adress         int
	comment        string
}

func (v VARS) String() string {
	//return fmt.Sprintf("Tag: %s\nType: %s\nAdress: %v\nComment: %s\n\n", v.tag, v.datatype, v.adress, v.comment)
	return fmt.Sprintf("%s,%s,%s,%v,%s\n", v.tag, v.datatype, v.globaldatatype, v.adress, v.comment)
}

func RmLeadSpace(s string) string {
	r := regexp.MustCompile(`^\s*`)
	s = r.ReplaceAllString(s, "")
	return s
}
