package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestXML(t *testing.T) {
	sifbytes, _ := ioutil.ReadFile("./xml/sif.xml")
	sif := Str(sifbytes)
	sif.SetEnC()

	// str1, nArr := XMLTagEleEx(sif.V(), "CourseTitle", 3)
	// fPln(str1, nArr)
	// fPln(" --------------------------------------- ")

	str, nArr := XMLXPathEle(sif.V(), "SchoolCourseInfo", " ~ ", 1)
	fPln(str, nArr)
	if str == "" {
		return
	}
	fPln(" --------------------------------------- ")	

	mapFT := &map[string][]string{}
	XMLFamilyTree(str, "", " ~ ", mapFT)
	//fPln((*mapFT))
	for k, v := range (*mapFT) {
		fPln(k, v)
	}
	fPln(" --------------------------------------- ")

	XMLArrByIPath(str, "SchoolCourseInfo#1", " ~ ", mapFT)
	fPln(" --------------------------------------- ")

	return

	fPln(XMLAttributes(str, "-"))
	fPln(" --------------------------------------- ")

	fPln(XMLChildren(str, true))
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
