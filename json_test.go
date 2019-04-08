package gjxy

import (
	"io/ioutil"
	"testing"
)

func TestJSONChild(t *testing.T) {

	jsonBytes, _ := ioutil.ReadFile("./json/content.json") // only LF at end of line
	s := Str(jsonBytes)
	s.SetEnC()

	{
		fPln(JSONXPathValue(s.V(), "abc ~ subject", " ~ ", 1, 1))
		fPln(JSONXPathValue(s.V(), "subject", " ~ ", 1))

		// mapFT := &map[string][]string{}
		// JSONFamilyTree(s, "xapi", " ~ ", mapFT)
		// fPln(mapFT)

		fPln(" ----------------------------------------------- ")

		mapFT, mapArrInfo := JSONArrInfo(s.V(), "xapi", " ~ ", "id", nil)
		for k, v := range *mapFT {
			fPln(k, v)
		}

		fPln(" ----------------------------------------------- ")

		for k, v := range *mapArrInfo {
			fPln(k, v)
		}
	}

	// jsonSample := Str(` [ 1, "a", 123, ["ab"], true, false, true, null, null ] `)
	// fPln(jsonSample.IsJSON())
	// fPln(jsonSample.IsJSONArray())

	// jsonSample = Str(` {"root" : 12334.34, "test": [ 1, 3, 4 ] } `)
	// fPln(jsonSample.IsJSONSingle())
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
	// 	children := s.JSONObjChildren()
	// 	fPln(children)
	// 	if len(children) != 10 {
	// 		t.Fatalf("JSONChildren() Error <%d>\n", len(children))
	// 	}
	// }
}

func TestJSONArrInfo(t *testing.T) {
	jsonBytesTemp, _ := ioutil.ReadFile("xapi.json")
	sTemp := Str(jsonBytesTemp)
	fPln(sTemp.SetEnC())

	{
		mapFT := &map[string][]string{}
		JSONFamilyTree(sTemp.V(), "", " ~ ", mapFT)
		for k, v := range *mapFT {
			fPln(k, "::", v)
		}

		jsonBytes, _ := ioutil.ReadFile("xapi.1.json")
		s := Str(jsonBytes)
		fPln(s.SetEnC())

		// mapFT = &map[string][]string{}
		// root, _, newJSON := JSONRootEx(s, "fakeRoot")
		// Str(newJSON).JSONFamilyTree(root, ".", mapFT)
		// fPln(*mapFT)

		mapFT, mapAI := JSONArrInfo(s.V(), "", " ~ ", "535e966a-931e-430f-a809-d90401147864", mapFT)
		// for k, v := range *mapFT {
		// 	fPln(k, "::", v)
		// }
		fPln("----------------------------------------------------------------------")
		if mapAI != nil {
			for k, v := range *mapAI {
				fPln(k, "::", v)
			}
		}
	}
}

func TestJSONMake(t *testing.T) {

	mIPathObj := map[string]string{}

	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#1", "name", "hello world")
	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#1", "fname", "world")
	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#1", "gname", "hello")
	JSONBuildIPath(mIPathObj, "ROOT", "StaffPersonal", "ROOT.StaffPersonal#1")
	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#2", "NAME", "HELLO WORLD")
	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#2", "FNAME", "WORLD")
	JSONBuildIPath(mIPathObj, "ROOT.StaffPersonal#2", "GNAME", "HELLO")
	JSONBuildIPath(mIPathObj, "ROOT", "StaffPersonal", "ROOT.StaffPersonal#2")
	JSONBuildIPath(mIPathObj, "ROOT", "StaffPersonal", "{}")

	//fPln(mIPathObj["ROOT"])
	//fPln("ROOT.StaffPersonal", mIPathObj["ROOT.StaffPersonal"])

AGAIN:
	for k, v := range mIPathObj {
		for k1, v1 := range mIPathObj {
			if k == k1 {
				continue
			}
			old := Str(v).RmQuotes(QDouble).V()
			if sCtn(old, k1) {
				mIPathObj[k] = sRep(v, "\""+k1+"\"", v1, -1)
				goto AGAIN
			}
		}
	}

	fPln(mIPathObj)
	root := mIPathObj["ROOT"]

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

	// // //json, ok = Str(json).JSONBuild("StaffPersonal", ".", 1, "-RefId", "{}")
	// json, ok = JSONBuild(Str(json), "StaffPersonal#1", ".", "LocalId", "946379881")
	// json, ok = JSONBuild(Str(json), "StaffPersonal#1", ".", "LocalId", "946379882")
	// json, ok = JSONBuild(Str(json), "StaffPersonal#1", ".", "LocalIdTest", "tttttttt")
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "StateProvinceId", "C2345681", 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "OtherIdList", "{}", 1)            //                                ***
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList", ".", "OtherId", "{}", 1, 1) //                     *** 1
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "-Type", "0004", 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "#content", "333333333", 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "PersonInfo", "{}", 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo", ".", "Name", "{}", 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.Name", ".", "-Type", "LGL", 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo", ".", "OtherNames", "{}", 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.OtherNames", ".", "Name", "[{},{}]", 1, 1, 1) //      ***
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.OtherNames.Name", ".", "-Type", "AKA", 1, 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.OtherNames.Name", ".", "-Type", "PRF", 1, 1, 1, 2) // ***
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo", ".", "Demographics", "{}", 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.Demographics", ".", "CountriesOfCitizenship", "{}", 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.Demographics.CountriesOfCitizenship", ".", "CountryOfCitizenship", "\"8104\"", 1, 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.Demographics.CountriesOfCitizenship", ".", "CountryOfCitizenship", "\"1101\"", 1, 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "LocalId", "946379883", 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "LocalIdTest", "iiiiiiii", 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.PersonInfo.Demographics.CountriesOfCitizenship", ".", "CountryOfCitizenship", "\"2202\"", 1, 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList", ".", "OtherId", "{}", 1, 1) //                     *** 2
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "-Type", "0005", 1, 1, 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "-Type1", "0008", 1, 1, 2)
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "#content", "44444444", 1, 1, 2)
	// // // json, ok = Str(json).JSONBuild("StaffPersonal", ".", 1, "LocalId", "{}")
	// json, ok = Str(json).JSONBuild("StaffPersonal", ".", "LocalId", "test2", 1)
	// json, ok = Str(json).JSONBuild("StaffPersonal.OtherIdList.OtherId", ".", "#content", "44444444", 1, 1, 2)
	// // // // fPln(Str(json).JSONXPath("StaffPersonal.PersonInfo.OtherNames.Name", ".", 1))

	// fPln(json, ok)
	ioutil.WriteFile("temp.json", []byte(root), 0666)
}

func TestJSONObjectMerge(t *testing.T) {

	json1 := ``
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
	ioutil.WriteFile("debug_temp.json", []byte(rst), 0666)
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

func TestGQLBuild(t *testing.T) {
	s := Str(`
type StaffPersonal {
	-RefId: String
	LocalId: String
	StateProvinceId: String
	OtherIdList: OtherIdList
}
	
type OtherIdList {
	OtherId: OtherId
}

type OtherIdList1 {
	OtherId1: OtherId1
}

type OtherIdList2 {
	OtherId2: OtherId2
}
	
type OtherId {
	-Type: String
}`)

	s = Str(SchemaBuild(s, "OtherIdList2", "OtherId2", "String"))

	// s := Str("")
	// s = Str(s.GQLBuild("StaffPersonal", "RefId", "String"))
	// s = Str(s.GQLBuild("StaffPersonal", "LocalId", "String"))
	// s = Str(s.GQLBuild("Recent", "SchoolLocalId", "String"))
	// s = Str(s.GQLBuild("Recent", "LocalCampusId", "String"))
	// s = Str(s.GQLBuild("StaffPersonal", "StateProvinceId", "String"))
	// s = Str(s.GQLBuild("NAPLANClassListType", "ClassCode", "[String]"))
	// s = Str(s.GQLBuild("StaffPersonal", "OtherIdList", "OtherIdList"))

	fPln(s)
}
