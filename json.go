package gjxy

import "encoding/json"

// IsJSON :                                                                                     &
func IsJSON(s Str) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s.V()), &js) == nil
}

// JSONFstEle : The first json child                                                            &
func JSONFstEle(s Str) string {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	str := s.T(BLANK).RmBrackets(BCurly).T(BLANK)
	if p := str.Idx(":"); p > 0 {
		str = str.S(0, p)
		str = str.T(BLANK)
	}
	return str.RmQuotes(QDouble).V()
}

// IsJSONSingle :
func IsJSONSingle(s Str) (ok bool, tag string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	s = s.T(BLANK)
	if s.C(0) == '{' && s.C(LAST) == '}' {
		s = s.S(1, ALL-1).T(BLANK)
		if s.QuotePairCount(QDouble) == 1 {
			ok, tag = true, s.S(0, s.Idx(":")).RmQuotes(QDouble).V()
		}
	}
	return
}

// IsJSONArray : Array Info be valid only on single-type elements.
func IsJSONArray(s Str) (ok bool, eleType JSONTYPE, n int, eles []string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	s = s.T(BLANK)
	if s.C(0) == '[' && s.C(LAST) == ']' {

		ok = true

		inBox := s.RmBrackets(BBox).T(BLANK)
		if inBox == "" {
			eleType, n, eles = JT_UNK, 0, []string{}
			return
		}

		switch inBox.C(0) {
		case '{':
			eleType, n = JT_OBJ, inBox.BracketPairCount(BCurly)
			for i := 1; i <= n; i++ {
				obj, _, r := inBox.BracketsPos(BCurly, 1, 1)
				eles = append(eles, obj.V())
				inBox = inBox.S(r+1, ALL)
			}
		case '[':
			eleType, n = JT_ARR, inBox.BracketPairCount(BBox)
			for i := 1; i <= n; i++ {
				arr, _, r := inBox.BracketsPos(BBox, 1, 1)
				eles = append(eles, arr.V())
				inBox = inBox.S(r+1, ALL)
			}
		case '"':
			eleType, n = JT_STR, inBox.QuotePairCount(QDouble)
			for i := 1; i <= n; i++ {
				str, _, r := inBox.QuotesPos(QDouble, 1)
				eles = append(eles, str.V())
				inBox = inBox.S(r+1, ALL)
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			eleType, n = JT_NUM, sCnt(inBox.V(), ",")+1
			for _, num := range sSpl(inBox.V(), ",") {
				eles = append(eles, Str(num).T(BLANK).V())
			}
		case 't', 'f':
			eleType, n = JT_BOOL, sCnt(inBox.V(), ",")+1
			for _, b := range sSpl(inBox.V(), ",") {
				eles = append(eles, Str(b).T(BLANK).V())
			}
		default:
			panic("Not implemented")
		}

	} else {
		ok, eleType, n, eles = false, JT_UNK, -1, nil
	}
	return
}

// JSONWrapRoot : if Not Single JSON, add Single "root", return the modified JSON.              &
func JSONWrapRoot(s Str, rootExt string) (root string, ext bool, extJSON string) {
	if ok, tag := IsJSONSingle(s); ok {
		return tag, false, s.V()
	}
	root, ext, extJSON = rootExt, true, fSf(`{ "%s": %s }`, rootExt, s.V())
	PC(!IsJSON(Str(extJSON)), fEf("JSONWrapRoot error"))
	return
}

// JSONChildValue : only return the value content                                               ?
func JSONChildValue(s Str, child string) (content string, cType JSONTYPE) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	child = Str(child).MkQuotes(QDouble).V()
	Lc := Str(child).L()

AGAIN:
	L := s.L()
	if ok, start, _ := s.SearchStrsIgnore(child, ":", BLANK); ok {
		above := s.S(0, start, L).V()
		sBelow := s.S(start, ALL, L)
		if sCnt(above, "{")-sCnt(above, "}") == 1 { // *** TRUELY FOUND ( Object OR Value ) ***
			if ok, s, _ := sBelow.SearchStrsIgnore(":", "{", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { //         *** object ***

				str, _, _ := sBelow.BracketsPos(BCurly, 1, 1)
				content, cType = str.V(), JT_OBJ

			} else if ok, s, _ := sBelow.SearchStrsIgnore(":", "[", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { //  *** array ***

				str, _, _ := sBelow.BracketsPos(BBox, 1, 1)
				content, cType = str.V(), JT_ARR

			} else if ok, s, _ := sBelow.SearchStrsIgnore(":", "\"", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { // *** string ***

				str, _, _ := sBelow.QuotesPos(QDouble, 2)
				content, cType = str.RmQuotes(QDouble).V(), JT_STR

			} else if ok, s, _ := sBelow.SearchAny2StrsIgnore([]string{":"}, DigStrArr, BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { // *** number ***

				_, value := sBelow.KeyValuePair(":", BLANK+",{", BLANK+",}", true, true)
				content, cType = value.V(), JT_NUM

			} else {
				panic("not implemented")
			}
		}

		// *** FAKE FOUND, Maybe above string's sub-element's same tag ***
		s = s.SegRep(start, start+2, "\"*")
		goto AGAIN
	}
	return
}

// JSONChildValueEx : if this child value is single-type array, e.g. [{},{}], return array count
// idx is only applicable on array-child. Normally from 1 to get an array-element.
// If idx is 0, get whole array.
func JSONChildValueEx(s Str, child string, idx int) (ele string, nArr int) {
	content, cType := JSONChildValue(s, child)
	if cType == JT_ARR {
		_, eType, n, eles := IsJSONArray(Str(content))
		nArr = n
		if idx == 0 || eType == JT_UNK {
			ele = content
			return
		}
		ele = eles[idx-1]
	} else {
		ele, nArr = content, -1
	}
	return
}

// JSONXPathValue : cannot use empty path to get current. indices is from the 1st path-seg's array index.    &
// if it's not array, use 1; if it is 0 and array, return the whole array
func JSONXPathValue(s Str, xpath, del string, indices ...int) (content string, nArr int) {
	PC(xpath == "", fEf("at least one path must be provided"))

	segs := sSpl(xpath, del)
	PC(len(segs) != len(indices), fEf("path & seg's index count not match"))
	for i := 0; i < len(indices)-1; i++ {
		PC(indices[i] <= 0, fEf("Only Last index can be 0 to get the whole array content"))
	}

	for i, seg := range segs {
		s = IF(content != "", Str(content), s).(Str)
		content, nArr = JSONChildValueEx(s, seg, indices[i])
	}
	return
}

// JSONObjChildren :
func JSONObjChildren(s Str) (children []string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	inCurly := s.T(BLANK).RmBrackets(BCurly).T(BLANK)

	posList := []int{}
	starts, _ := inCurly.Indices2StrsIgnore("\"", ":", BLANK) //         *** slow, flexible ***
	// starts := inCurly.Indices(`":`) //                                *** fast, but must be `":` ***
	for _, p := range starts {
		if inCurly.BracketDepth(BCurly, p) == 0 {
			posList = append(posList, p)
		}
	}
	for _, pe := range posList {
		str := inCurly.S(0, pe)
		if ps := str.LIdx("\""); ps >= 0 {
			children = append(children, str.S(ps+1, ALL).V())

			// *** search "[" for array, add prefix "[]" to array child ***
			if inCurly.S(pe+1, ALL).TL(BLANK + ":").HP("[") {
				children[len(children)-1] = "[]" + children[len(children)-1]
			}
		}
	}
	return
}

// JSONFamilyTree : Use One Sample to get FamilyTree, DO NOT use long array data                                   &
func JSONFamilyTree(s Str, fName, del string, mapFT *map[string][]string) {
	PC(mapFT == nil, fEf("FamilyTree return map is not initialised !"))
	// fPln(fName) // DEBUG

	if s.T(BLANK).C(0) != '{' && s.T(BLANK).C(LAST) != '}' {
		return
	}

	if children := JSONObjChildren(s); len(children) > 0 {
		fPln(fName, children) // DEBUG

		(*mapFT)[fName] = children //                         *** record path ***

		for _, child := range children {
			if Str(child).HP("[]") {
				child = Str(child).S(2, ALL).V()
			}
			nextPath := Str(fName + del + child).T(del).V()
			indices := []int{}
			for range sSpl(nextPath, del) {
				indices = append(indices, 1)
			}

			str, _ := JSONXPathValue(s, nextPath, del, indices...)
			JSONFamilyTree(Str(str), nextPath, del, mapFT)
		}
	}
}

// // ** mapPC : map[path]children
// func mapSortByPathKey(mapPC *map[string][]string) []string {

// }

// // JSONArrInfo1 :
// func (s Str) JSONArrInfo1(xpath, del, id string, mapFT *map[string][]string) (*map[string][]string, *map[string]struct {
// 	Count int
// 	ID    string
// }) {
// 	if mapFT == nil {
// 		mapFT = &map[string][]string{}
// 		s.JSONFamilyTree(xpath, del, mapFT)
// 	}

// }

// JSONArrInfo : Only deal with SAME TYPE element array                                                            ?
func JSONArrInfo(s Str, xpath, del, id string, mapFT *map[string][]string) (*map[string][]string, *map[string]struct {
	Count int
	ID    string
}) {
	if mapFT == nil {
		mapFT = &map[string][]string{}
		JSONFamilyTree(s, xpath, del, mapFT)
	}

	// for k, v := range *mapFT {
	// 	fPln(k, " : ", v)
	// }
	// return nil, nil
	// fPln(" ------------------------------------------------- ")

	mapA := map[string]bool{}
	for attr, children := range *mapFT {
		for _, child := range children {
			sChild := Str(child)
			if sChild.HP("[]") {
				sChild = sChild.S(2, ALL)
				path := Str(attr + del + sChild.V()).T(del).V()
				mapA[path] = true
			}
		}
	}

	for k, v := range mapA {
		fPln(k, " : ", v)
	}
	//return nil, nil

	mapAC := &map[string]struct {
		Count int
		ID    string
	}{}

	for k := range mapA {
		if len(sSpl(k, del)) == 1 {
			if _, n := JSONChildValueEx(s, k, 0); n >= 0 {
				(*mapAC)[k] = struct {
					Count int
					ID    string
				}{Count: n, ID: id}
			}
		}
	}

	// for k := range mapA {
	// 	ss := sSpl(k, del)
	// 	if len(ss) == 2 {
	// 		s1, s2, s12, s12ns := ss[0], ss[1], sJ(ss, del), []string{}
	// 		if cntid1, ok := (*mapAC)[s1]; ok { //                                ** get upper level's count
	// 			for i := 1; i <= cntid1.Count; i++ {
	// 				s12ns = append(s12ns, fSf("%s#%d%s%s", s1, i, del, s2))
	// 			}
	// 			for i, ns := range s12ns {
	// 				idx := Str(sSpl(sSpl(ns, del)[0], "#")[1]).ToInt()
	// 				_, _, _, n := s.JSONXPathValue(s12, del, []int{idx, 0}...) // ** get this level's count
	// 				(*mapAC)[s12ns[i]] = struct {
	// 					Count int
	// 					ID    string
	// 				}{Count: n, ID: id}
	// 			}
	// 		} else {
	// 			if _, _, _, n := s.JSONXPathValue(s12, del, []int{1, 0}...); n >= 0 {
	// 				(*mapAC)[s12] = struct {
	// 					Count int
	// 					ID    string
	// 				}{Count: n, ID: id}
	// 			}
	// 		}
	// 	}
	// }

	// for k := range mapA {
	// 	ss := sSpl(k, del)
	// 	if len(ss) == 3 {
	// 		s1, s1nArr := ss[0], []string{}
	// 		if n1, ok := (*mapAC)[s1]; ok {
	// 			for i := 1; i <= n1.Count; i++ {
	// 				s1nArr = append(s1nArr, fSf("%s#%d", s1, i))
	// 			}
	// 		} else {

	// 		}

	// 		s2, s1ns2nArr := ss[1], []string{}
	// 		for _, s1n := range s1nArr {
	// 			s1ns2 := s1n + del + s2
	// 			if n12, ok := (*mapAC)[s1ns2]; ok {
	// 				for i := 1; i <= n12.Count; i++ {
	// 					s1ns2nArr = append(s1ns2nArr, fSf("%s#%d", s1ns2, i))
	// 				}
	// 			} else {

	// 			}
	// 		}

	// 		s3 := ss[2]
	// 		for _, sn := range s1ns2nArr {
	// 			sArr, indices := []string{}, []int{}
	// 			for _, sn := range sSpl(sn, del) {
	// 				strnum := sSpl(sn, "#")
	// 				sArr = append(sArr, strnum[0])
	// 				indices = append(indices, Str(strnum[1]).ToInt())
	// 			}
	// 			path := sJ(sArr, del) + del + s3
	// 			indices = append(indices, 0)
	// 			if _, _, _, n := s.JSONXPathValue(path, del, indices...); n >= 0 { // ** get this level's count
	// 				(*mapAC)[sn+del+s3] = struct {
	// 					Count int
	// 					ID    string
	// 				}{Count: n, ID: id}
	// 			}
	// 		}
	// 	}
	// }

	// ****************************************************************************************

	//for k := range mapA {
	//ss := sSpl(k, del)
	//lss := len(ss)

	// n := -1
	// if lss == 1 {

	// 	indices := []int{0}
	// 	_, _, _, n = s.JSONXPathValue(k, del, indices...)
	// 	(*mapAC)[k] = struct {
	// 		Count int
	// 		ID    string
	// 	}{Count: n, ID: id}

	// } else if lss == 2 {

	// 	s1, s2, s12, s12ns := ss[0], ss[1], sJ(ss, del), []string{}
	// 	_, _, _, n = s.JSONXPathValue(s1, del, []int{0}...)
	// 	for i := 1; i <= n; i++ {
	// 		s12ns = append(s12ns, fSf("%s#%d%s%s", s1, i, del, s2))
	// 	}
	// 	for i, ns := range s12ns {
	// 		idx := Str(sSpl(sSpl(ns, del)[0], "#")[1]).ToInt()
	// 		_, _, _, n = s.JSONXPathValue(s12, del, []int{idx, 0}...)
	// 		(*mapAC)[s12ns[i]] = struct {
	// 			Count int
	// 			ID    string
	// 		}{Count: n, ID: id}
	// 	}

	// } else if lss == 3 {

	// }

	/******************************************/

	// if lss < 3 {

	// 	indices := TerOp(len(ss) == 1, []int{0}, []int{1, 0}).([]int)
	// 	_, _, _, n := s.JSONXPathValue(k, del, indices...)
	// 	(*mapAC)[k] = struct {
	// 		Count int
	// 		ID    string
	// 	}{Count: n, ID: id}

	// } else if lss == 3 {

	// 	s12, s3, s123, s123ns := sJ(ss[:2], del), ss[2], sJ(ss, del), []string{}
	// 	_, _, _, n := s.JSONXPathValue(s12, del, []int{1, 0}...)
	// 	for i := 1; i <= n; i++ {
	// 		s123ns = append(s123ns, fSf("%s#%d%s%s", s12, i, del, s3))
	// 	}
	// 	// fPln(s123ns)

	// 	for i, ns := range s123ns {
	// 		idx := Str(sSpl(sSpl(ns, del)[1], "#")[1]).ToInt()
	// 		_, _, _, n = s.JSONXPathValue(s123, del, []int{1, idx, 0}...)
	// 		(*mapAC)[s123ns[i]] = struct {
	// 			Count int
	// 			ID    string
	// 		}{Count: n, ID: id}
	// 	}

	// } else if lss == 4 {

	// 	// s123, s4, s1234, s1234ns := sJ(ss[:3], del), ss[3], sJ(ss, del), []string{}
	// 	// _, _, _, n := s.JSONXPathValue(s12, del, []int{1, 0}...)

	// } else {

	// 	panic("haven't implemented this level nested array function")

	// }
	//}

	return mapFT, mapAC
}

// JSONObjectMerge :
func JSONObjectMerge(s Str, json string) (rst string) {
	jsonStr := Str(json)
	PC((s.V() != "" && !IsJSON(s)) || !IsJSON(jsonStr), fEf("Error: Invalid JSON string"))

	if s.V() == "" {
		return json
	}

	root1, root2 := JSONFstEle(s), JSONFstEle(jsonStr)
	PC(root1 != root2, fEf("Error: Different JSON object"))

	content1, _ := JSONChildValueEx(s, root1, 0)
	content2, _ := JSONChildValueEx(jsonStr, root2, 0)
	content1, content2 = Str(content1).TL("["+BLANK).V(), Str(content2).TL("["+BLANK).V()
	content1, content2 = Str(content1).TR("]"+BLANK).V(), Str(content2).TR("]"+BLANK).V()

	rst = fSf(`{ "%s": [ %s, %s ] }`, root1, content1, content2)
	PC(!IsJSON(Str(rst)), fEf("Error: Json Merge Result"))
	return
}

// JSONBuild : NOT support mixed (atomic & object) types in one array                                              &
// func (s Str) JSONBuild(xpath, del, property, value string, indices ...int) (string, bool) {
// 	if s.T(BLANK) == "" {
// 		s = Str("{}")
// 	}

// 	PC(len(indices) != len(sSpl(xpath, del)), fEf("indices count must be xpath seg's count"))

// 	property = Str(property).MkQuotes(QDouble).V() + ": "
// 	sValue := Str(value)
// 	if sValue.C(0) != '{' && sValue.C(0) != '[' {
// 		value = sValue.MkQuotes(QDouble).V()
// 	}

// 	if s == "{}" {
// 		s = Str(fSf(`{ %s%s}`, property, value))
// 		return s.V(), s.IsJSON()
// 	}

// 	if sub, start, end, _ := s.JSONXPathValue(xpath, del, indices...); start != -1 {

// 		for _, p0 := range Str(sub).Indices(property) { //                               ** incoming p-v's property already exists **
// 			sub02p0 := Str(sub).S(0, p0).V()
// 			if sCnt(sub02p0, "{")-sCnt(sub02p0, "}") == 1 { //                           ** 1 level child property **

// 				Subp02end := Str(sub).S(p0, ALL)
// 				if p1 := Subp02end.Idx(property + "["); p1 == 0 { //                     ** already array format, 3rd, 4th... coming **
// 					box, _, _ := Str(sub).BracketsPos(BBox, 1, 1)
// 					inBox := box.RmBrackets(BBox).TrimAllLMR(BLANK).V()
// 					ss := append(sSpl(inBox, ","), value)
// 					newBox := Str(sJ(ss, ", ")).MkBrackets(BBox).V()
// 					sub = sRep(sub, box.V(), newBox, 1)
// 				} else { //                                                              ** only one exists, the 2nd coming, change to array format **
// 					k, v := Subp02end.KeyValuePair(":", " ", " ", false, false)
// 					if !v.HP("{") { //                                                   ** atomic value **
// 						v = v.TR(",")
// 						sub = sRep(sub, property+v.V(), property+"["+v.V()+", "+value+"]", 1)
// 					} else { //                                                          ** object value **
// 						sub = sRep(sub, k.V()+":", k.V()+": [", 1)
// 						sub = Str(sub).S(0, ALL-1).T(BLANK).V() + ", " + "{}" + " ] " + "}"
// 					}
// 				}

// 				left, right := s.S(0, start).V(), s.S(end+1, ALL).V()
// 				json := left + sub + right
// 				return json, Str(json).IsJSON()
// 			}
// 		}

// 		// **********************************************

// 		left, right := s.S(0, end).T(BLANK).V(), s.S(end, ALL).T(BLANK).V()
// 		left = TerOp(!Str(left).HS("{"), left+",", left).(string)
// 		json := left + " " + Str(property+value).T(BLANK).V() + " " + right
// 		return json, Str(json).IsJSON()
// 	}
// 	return "", false
// }
