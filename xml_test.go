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

	str, nArr := XMLXPathEle(sif.V(), "SchoolInfo", " ~ ", 1)
	fPln(str, nArr)
	if str == "" {
		return
	}
	fPln(" --------------------------------------- ")

	// mapFT := &map[string][]string{}
	// XMLFamilyTree(str, "", " ~ ", mapFT)
	// //fPln((*mapFT))
	// for k, v := range (*mapFT) {
	// 	fPln(k, v)
	// }
	// fPln(" --------------------------------------- ")

	// XMLCntByIPath(str, "SchoolCourseInfo#1", " ~ ", mapFT)
	// fPln(" --------------------------------------- ")

	mFT, mArr := XMLCntInfo(str, "", " ~ ", "ArrGUID", nil)
	for k, v := range *mFT {
		fPln(k, v)
	}
	fPln("    -------------------- ")
	for k, v := range *mArr {
		fPln(k, v)
	}

	return

	fPln(XMLAttributes(str, "-"))
	fPln(" A --------------------------------------- ")

	fPln(XMLChildren(str, true))
	fPln(" --------------------------------------- ")

	tag, xml, l, r := XMLSegPos(sif.V(), 3, 1)
	fPln(tag)
	fPln(xml)
	fPln(l, r)

	fPln(XMLSegsCount(sif.V()))

	fPln(" --------------------------------------- ")

	tag, xml, l, r = XMLSegPos(sif.V(), 1, 726)
	fPln(tag)
	fPln(xml)
	fPln(l, r)

	fPln(XMLSegsCount(xml))
}
