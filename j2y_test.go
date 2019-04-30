package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestJstr2Y(t *testing.T) {
	bytes, _ := ioutil.ReadFile("./files/content.json")
	ystr := Jstr2Y(string(bytes))
	ioutil.WriteFile("test.yaml", []byte(ystr), 0666)
}
