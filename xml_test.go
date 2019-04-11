package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestXML(t *testing.T) {
	sifbytes, _ := ioutil.ReadFile("./xml/sif.xml")
	sif := Str(sifbytes)
	sif.SetEnC()

	// str := XMLTagEleEx(sif.V(), "CourseTitle", 3)
	// fPln(str)
	// fPln(" --------------------------------------- ")

	str := XMLXPathEle(sif.V(), "SchoolCourseInfo", " ~ ", 1)
	fPln(str)
	fPln(" --------------------------------------- ")

	mapFT := &map[string][]string{}
	XMLFamilyTree(str, "", " ~ ", mapFT)
	//fPln((*mapFT))
	for k, v := range (*mapFT) {
		fPln(k, v)
	}
	fPln(" --------------------------------------- ")

	return

	fPln(XMLAttributes(str, "-"))
	fPln(" --------------------------------------- ")

	fPln(XMLChildren(str))
	fPln(" --------------------------------------- ")

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
