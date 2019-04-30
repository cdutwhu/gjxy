package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestXfile2Y(t *testing.T) {
	y := Xfile2Y("./files/nswdig.xml")
	ioutil.WriteFile("./files/nswdig.yaml", []byte(y), 0666)
}
