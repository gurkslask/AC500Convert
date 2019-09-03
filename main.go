package AC500Convert

import "fmt"
import "io/ioutil"
import "regexp"
import "strings"

func Openfile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
	adress         string
	comment        string
}
type VARSModbus struct {
	tag            string
	datatype       string
	globaldatatype string
	adressHigh     string
	adressLow      string
	comment        string
}

func (v VARS) String() string {
	return fmt.Sprintf("%s,%s,%s,%v,%s\r", v.tag, v.datatype, v.globaldatatype, v.adress, v.comment)
}

func RmLeadSpace(s string) string {
	r := regexp.MustCompile(`^\s*`)
	s = r.ReplaceAllString(s, "")
	return s
}

func OutputToText(vars []VARS) []string {
	var s []string
	s = append(s, "//\nName,DataType,GlobalDataType,Address_1,Description //")
	for _, v := range vars {
		s = append(s, v.String())
	}
	return s
}
func raiseModbusCounter(steps int, low, high *int) {
	for i := 0; i < steps; i++ {
		*low++
		if *low > 7 {
			*low = 0
			*high++
		}
	}
}
