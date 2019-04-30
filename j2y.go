package gjxy

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Jstr2Y : JSON string to YAML string
func Jstr2Y(jsonstr string) string {
	yamlbytes, err := yaml.JSONToYAML([]byte(jsonstr))
	PE1(err, "error on yaml.JSONToYAML")
	return string(yamlbytes)
	// return YAMLRmHangStr(string(yamlbytes)) /* avoid hanging string value */
}

// Jb2Yb : JSON Bytes to YAML Bytes
func Jb2Yb(jsonbytes []byte) []byte {
	yamlbytes, err := yaml.JSONToYAML(jsonbytes)
	PE1(err, "error on yaml.JSONToYAML")
	return yamlbytes
}

// Jfile2Y : JSON file to YAML string
func Jfile2Y(jsonfile string) string {
	jsonbytes, err := ioutil.ReadFile(jsonfile)
	PE1(err, "error on ioutil.ReadFile")
	return Jstr2Y(string(jsonbytes))
}

// Jfile2Yb : JSON file to YAML Bytes
func Jfile2Yb(jsonfile string) []byte {
	jsonbytes, err := ioutil.ReadFile(jsonfile)
	PE1(err, "error on ioutil.ReadFile")
	return Jb2Yb(jsonbytes)
}
