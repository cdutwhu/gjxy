package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestJSONChild(t *testing.T) {

	jsonBytes, _ := ioutil.ReadFile("./json/xapi.json") // only LF at end of line
	s := Str(jsonBytes)

	ok, tp, n, _ := IsJSONArrOnFmtL0(s.V())
	fPln(ok, tp, n)

	jsonSample := ` [ -1, "a", 123, ["ab"], -1, true, false, true, null, null ] `
	fPln(IsJSON(jsonSample))
	fPln(IsJSONArr(jsonSample))

	// {
	// 	// fPln(JSONXPathValue(s.V(), "abc ~ subject", " ~ ", 1, 1))
	// 	// fPln(JSONXPathValue(s.V(), "subject", " ~ ", 1))

	// 	// mapFT := &map[string][]string{}
	// 	// JSONFamilyTree(s, "xapi", " ~ ", mapFT)
	// 	// fPln(mapFT)

	// 	_, _, json := JSONWrapRoot(s.V(), "xapi")

	// 	fPln(" ----------------------------------------------- ")

	// 	mFT, mArrInfo := JSONArrInfo(json, "", " ~ ", "GUID", nil)
	// 	for k, v := range *mFT {
	// 		fPln(k, v)
	// 	}

	// 	fPln(" ----------------------------------------------- ")

	// 	for k, v := range *mArrInfo {
	// 		fPln(k, v)
	// 	}

	// 	fPln(" ----------------------------------------------- ")
	// }

	// jsonSample := `{
	// 	"ROOT" : {
	// 	    "root" : 12334.34,
	// 	    "test": [ 1, 3, 4 ]
	// 	}
	// }`
	// jsonSample = `{ "score" : 123, "age": 10  }`
	// fPf("IsJSONSingle :")
	// fPln(IsJSONSingle(jsonSample))
	// fPln(jsonSample.JSONChildValue("root"))
	// fPln(jsonSample.JSONWrapRoot("Root"))

	// if !s.IsJSON() {
	// 	t.Fatalf("JSON Format Error\n")
	// }
	// if s.JSONFstEle() != "test最最" {
	// 	t.Fatalf("JSONRoot() Error\n")
	// }
	// if ok, _ := s.IsJSONSingle(); ok {
	// 	t.Fatalf("IsJSONSingle() Error\n")
	// }

	// if root, ext, json := s.JSONWrapRoot("fakeRoot"); !ext {
	// 	t.Fatalf("JSONWrapRoot() Error <%s> <%v> <%s>\n", root, ext, json)
	// }

	// arr, vType := s.JSONChildValue("mixArr")
	// fPln(arr, JSONTypeDesc[vType])

	// fPln(s.JSONChildValueEx("numArr", 2))
	// fPln(s.JSONChildValueEx("objArr", 3))
	// fPln(s.JSONChildValueEx("arrArr", 4))
	// fPln(s.JSONChildValueEx("strArr", 2))

	// {
	// 	content, nArr := s.JSONXPathValue("actor.member.innerArr", ".", 1, 1, 0)
	// 	if content != `[1, 2, 3]` || nArr != 3 {
	// 		t.Fatalf("JSONXPathValue() Error <%s> <%d>\n", content, nArr)
	// 	}
	// 	content, nArr = s.JSONXPathValue("actor.member.innerArr", ".", 1, 2, 2)
	// 	if content != `5` || nArr != 4 {
	// 		t.Fatalf("JSONXPathValue() Error <%s> <%d>\n", content, nArr)
	// 	}
	// 	content, nArr = s.JSONXPathValue("actor.member.account.homePage", ".", 1, 3, 1, 1)
	// 	if content != `http://www.example3.com` || nArr != -1 {
	// 		t.Fatalf("JSONXPathValue() Error <%s> <%d>\n", content, nArr)
	// 	}
	// 	content, nArr = s.JSONXPathValue("actor.member.mbox_sha1sum", ".", 1, 2, 1)
	// 	if content != `ebd31e95054c018b1072最7ccffd2ef2ec3a016ee9222` || nArr != -1 {
	// 		t.Fatalf("JSONXPathValue() Error <%s> <%d>\n", content, nArr)
	// 	}
	// }

	// {
	// 	children := JSONObjChildren(s.V())
	// 	fPln(children)
	// 	// if len(children) != 10 {
	// 	// 	t.Fatalf("JSONChildren() Error <%d>\n", len(children))
	// 	// }
	// }
}

func TestJSONMake(t *testing.T) {

	mIPathObj := map[string]string{}

	JSONMakeIPath(mIPathObj, "ROOT", "StaffPersonal", "ROOT.StaffPersonal#1", true)
	JSONMakeIPath(mIPathObj, "ROOT", "StaffPersonal", "ROOT.StaffPersonal#2", true)
	JSONMakeIPath(mIPathObj, "ROOT", "NAME", "HELLO ROOT", true)
	JSONMakeIPath(mIPathObj, "ROOT", "GNAME", "HELLO", false)
	JSONMakeIPath(mIPathObj, "ROOT", "GNAME", "WORLD 1", false)
	JSONMakeIPath(mIPathObj, "ROOT", "GNAME", "WORLD 2", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1", "name", "hello world", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1", "fname", "world", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1", "gname", "hello", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2", "NAME", "HELLO WORLD", true)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2", "FNAME", "WORLD", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2", "GNAME", "HELLO", false)

	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "sub name", "hello sub world", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "sub fname", "sub world", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1", "test", "ROOT.StaffPersonal#1.SubTest", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "sub gname", "sub hello", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB NAME", "HELLO SUB WORLD A", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB NAME", "HELLO SUB WORLD B", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB NAME", "HELLO SUB WORLD C", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB NAME", "HELLO SUB WORLD D", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "SUB NAME", "HELLO SUB WORLD 1", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "SUB NAME", "HELLO SUB WORLD 2", false)
	JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#1.SubTest", "SUB NAME", "HELLO SUB WORLD 3", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB FNAME", "SUB WORLD", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2", "test", "ROOT.StaffPersonal#2.SubTest", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB GNAME", "SUB HELLO", false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB GNAME 1", true, false)
	// JSONMakeIPath(mIPathObj, "ROOT.StaffPersonal#2.SubTest", "SUB GNAME 1", false, false)

	//fPln(mIPathObj["ROOT"])
	//fPln("ROOT.StaffPersonal", mIPathObj["ROOT.StaffPersonal"])

	root := JSONMakeIPathRep(mIPathObj, ".")
	ioutil.WriteFile("temp.json", []byte(root), 0666)

	// mIPathStr["ROOT#1.StaffPersonal"] = JSONBuildObj("", "ROOT#1.StaffPersonal", "name", "hello", false)
	// fPln("11", mIPathStr["ROOT#1.StaffPersonal"])

	// mIPathStr["ROOT#1.StaffPersonal"] = JSONBuildObj(mIPathStr["ROOT#1.StaffPersonal"], "ROOT#1.StaffPersonal", "name1", "staff", false)
	// fPln("12", mIPathStr["ROOT#1.StaffPersonal"])

	// StaffPersonal, ok1 = JSONBuildObj(StaffPersonal, "ROOT#1.StaffPersonal", "age", 11, false)
	// fPln("13", StaffPersonal, ok1)

	// root := JSONBuildObj("", "ROOT", "StaffPersonal", "ROOT#1.StaffPersonal", false)
	// fPln("1", root)

	// root, ok = JSONBuildObj(root, "root", "StaffPersonal", StaffPersonal, true)
	// fPln("1", root, ok)
}

func TestJSONObjectMerge(t *testing.T) {

	json1 := `{
		"TeachingGroup1": {
			"RefId": "F47C2C6D-BD49-40E6-A430-111111111112"
		}
	}`
	json2 := `{
		"TeachingGroup1": {
			"RefId": "F47C2C6D-BD49-40E6-A430-111111111111",
			"SchoolYear": "2018",
			"LocalId": "2018-English-8-1-B",
			"ShortName": "8B English 1",
			"LongName": "Year 8B English 1",
			"TimeTableSubjectRefId": "FD3E4B1F-0FC6-4607-BB95-78791ABA8997",
			"TeacherList": {
				"TeachingGroupTeacher": {
					"StaffPersonalRefId": "D4A3C1E3-3F6E-4B31-ABA6-26809DF5FD63",
					"StaffLocalId": "kafaj506",
					"Name": {
						"Type": "LGL",
						"FamilyName": "Knoll",
						"GivenName": "Ina"
					},
					"Association": "Class Teacher"
				}
			},
			"TeachingGroupPeriodList": {
				"TeachingGroupPeriod": [
					{
						"RoomNumber": "171",
						"DayId": "Fr",
						"PeriodId": "12:00:00"
					},
					{
						"RoomNumber": "166",
						"DayId": "Fr",
						"PeriodId": "15:00:00"
					}
				]
			}
		}
	}`

	rst := JSONObjectMerge(json1, json2)
	ioutil.WriteFile("test_temp.json", []byte(rst), 0666)
}

func TestJSONRoot(t *testing.T) {
	// jsonbytes, _ := ioutil.ReadFile("./test1.json")
	json := `{
		"TeachingGroup": {
			"RefId": "F47C2C6D-BD49-40E6-A430-360333274DB2",
			"SchoolYear": "2018",
			"LocalId": "2018-English-8-1-B",
			"ShortName": "8B English 1",
			"LongName": "Year 8B English 1",
			"TimeTableSubjectRefId": "FD3E4B1F-0FC6-4607-BB95-78791ABA8997",
			"TeacherList": {
				"TeachingGroupTeacher": {
					"StaffPersonalRefId": "D4A3C1E3-3F6E-4B31-ABA6-26809DF5FD63",
					"StaffLocalId": "kafaj506",
					"Name": {
						"Type": "LGL",
						"FamilyName": "Knoll",
						"GivenName": "Ina"
					},
					"Association": "Class Teacher"
				}
			}
		}
	}`
	fPln(JSONFstEle(json))
	// root, ext, newJSON := Str(json).JSONRootEx("MyRoot")
	// fPln(root)
	// if ext {
	// 	json = newJSON
	// }

	// mapFT, mapAC := Str(json).JSONArrInfo("", " ~ ", "1234567890", nil)
	// if mapFT != nil && mapAC != nil {
	// 	for k, v := range *mapFT {
	// 		fPln(k, v)
	// 	}
	// 	fPln("<----------------------------------->")
	// 	for k, v := range *mapAC {
	// 		fPln(k, v)
	// 	}
	// }
}
