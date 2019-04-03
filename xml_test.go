package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestXML(t *testing.T) {
	sifbytes, _ := ioutil.ReadFile("./xml/sif.xml")
	sif := Str(sifbytes)
	sif.SetEnC()

	tag, xml, l, r := XMLSegPos(sif, 3, 1)
	fPln(tag)
	fPln(xml)
	fPln(l, r)

	fPln(XMLSegsCount(sif))

	fPln(" --------------------------------------- ")

	tag, xml, l, r = XMLSegPos(sif, 1, 726)
	fPln(tag)
	fPln(xml)
	fPln(l, r)

	fPln(XMLSegsCount(Str(xml)))
}
