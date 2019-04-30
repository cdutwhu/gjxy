package gjxy

import (
	"io/ioutil"
	"strings"

	xj "github.com/basgys/goxml2json"
)

// Xstr2J : XML string to JSON string
func Xstr2J(xmlstr string) string {
	xml := strings.NewReader(xmlstr)
	json, err := xj.Convert(xml)
	PE1(err, "error on xj.Convert")
	return json.String()
}

// Xfile2J : XML file to JSON string
func Xfile2J(xmlfile string) string {
	xmlbytes, err := ioutil.ReadFile(xmlfile)
	PE1(err, "error on ioutil.ReadFile")
	return Xstr2J(string(xmlbytes))
}
