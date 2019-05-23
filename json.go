package gjxy

import (
	"encoding/json"
	"sort"
)

// IsJSON :                                                                                     &
func IsJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

// JSONFstEle : The first json child                                                            &
func JSONFstEle(s string) string {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	str := Str(s).T(BLANK).RmBrackets(BCurly).T(BLANK)
	str = str.STo(":").T(BLANK)
	return str.RmQuotes(QDouble).V()
}

// IsJSONSingle :
func IsJSONSingle(s string) (ok bool, tag string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	S := Str(s).T(BLANK)
	if S.C(0) == '{' && S.C(LAST) == '}' {
		S = S.S(1, ALL-1).T(BLANK)
		switch S.BracketPairCount(BCurly) {
		case 0:
			if sCnt(S.V(), ":") == 1 {
				ok, tag = true, S.STo(":").T(BLANK).RmQuotes(QDouble).V()
			}
		case 1:
			leftS := S.STo("{").V()
			if sCnt(leftS, ":") == 1 && sCnt(leftS, ",") == 0 && S.C(LAST) == '}' {
				ok, tag = true, S.STo(":").T(BLANK).RmQuotes(QDouble).V()
			}
		}
	}
	return
}

// IsJSONArray : Array Info be valid only on single-type elements.
func IsJSONArray(s string) (ok bool, eleType JSONTYPE, n int, eles []string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	S := Str(s).T(BLANK)
	if S.C(0) == '[' && S.C(LAST) == ']' {

		ok = true

		inBox := S.RmBrackets(BBox).T(BLANK)
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
func JSONWrapRoot(s, rootExt string) (root string, ext bool, extJSON string) {
	if ok, tag := IsJSONSingle(s); ok {
		return tag, false, s
	}
	root, ext, extJSON = rootExt, true, fSf(`{ "%s": %s }`, rootExt, s)
	PC(!IsJSON(extJSON), fEf("JSONWrapRoot error"))
	return
}

// JSONChildValue :                                                                                   ?
func JSONChildValue(s, child string) (content string, cType JSONTYPE) {
	if s == "" {
		return "", JT_UNK
	}

	PC(!IsJSON(s), fEf("Invalid JSON"))
	child = Str(child).MkQuotes(QDouble).V()
	Lc := Str(child).L()

AGAIN:
	L := Str(s).L()
	if ok, start, _ := Str(s).SearchStrsIgnore(child, ":", BLANK); ok {
		above := Str(s).S(0, start, L).V()
		sBelow := Str(s).S(start, ALL, L)
		if sCnt(above, "{")-sCnt(above, "}") == 1 { // *** TRUELY FOUND ( Object OR Value ) ***
			if ok, s, _ := sBelow.SearchStrsIgnore(":", "{", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { //         *** object ***
				str, _, _ := sBelow.BracketsPos(BCurly, 1, 1)
				content, cType = str.V(), JT_OBJ
			} else if ok, s, _ := sBelow.SearchStrsIgnore(":", "[", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { //  *** array ***
				str, _, _ := sBelow.BracketsPos(BBox, 1, 1)
				content, cType = str.V(), JT_ARR
			} else if ok, s, _ := sBelow.SearchStrsIgnore(":", "\"", BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { // *** string ***
				str, _, _ := sBelow.QuotesPos(QDouble, 2)
				//content, cType = str.RmQuotes(QDouble).V(), JT_STR
				content, cType = str.V(), JT_STR //                      ** keep string value's quotations **
			} else if ok, s, _ := sBelow.SearchAny2StrsIgnore([]string{":"}, DigStrArr, BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { // *** number ***
				_, value := sBelow.KeyValuePair(":", BLANK+",{", BLANK+",}", true, true)
				content, cType = value.V(), JT_NUM
			} else if ok, s, _ := sBelow.SearchAny2StrsIgnore([]string{":"}, []string{"true", "false"}, BLANK); ok && sBelow.S(0, s).T(BLANK).L() == Lc { // *** bool ***
				_, value := sBelow.KeyValuePair(":", BLANK+",{", BLANK+",}", true, true)
				content, cType = value.V(), JT_BOOL
			} else {
				panic("not implemented")
			}
		}

		// *** FAKE FOUND, Maybe above string's sub-element's same tag ***
		s = Str(s).SegRep(start, start+2, "\"*").V()
		goto AGAIN
	}
	return
}

// JSONChildValueEx : if this child value is single-type array, e.g. [{},{}], return array count
// idx is only applicable on array-child. Normally from 1 to get an array-element.
// If idx is 0, get whole array.
func JSONChildValueEx(s, child string, idx int) (ele string, nArr int) {
	content, cType := JSONChildValue(s, child)
	if cType == JT_ARR {
		_, eType, n, eles := IsJSONArray(content)
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
func JSONXPathValue(s, xpath, del string, indices ...int) (content string, nArr int) {
	PC(xpath == "", fEf("at least one path must be provided"))

	segs := sSpl(xpath, del)
	PC(len(segs) != len(indices), fEf("path & seg's index count not match"))
	for i := 0; i < len(indices)-1; i++ {
		PC(indices[i] <= 0, fEf("Only Last index can be 0 to get the whole array content"))
	}

	for i, seg := range segs {
		s = IF(content != "", content, s).(string)
		content, nArr = JSONChildValueEx(s, seg, indices[i])
	}
	return
}

// JSONObjChildren :
func JSONObjChildren(s string) (children []string) {
	PC(!IsJSON(s), fEf("Invalid JSON"))
	inCurly := Str(s).T(BLANK).RmBrackets(BCurly).T(BLANK)

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

// JSONFamilyTree : Use One Sample to get FamilyTree, DO NOT use long array data                             &
func JSONFamilyTree(s, fName, del string, mFT *map[string][]string) {
	PC(mFT == nil, fEf("FamilyTree return map is not initialised !"))
	// fPln(fName) // DEBUG

	if Str(s).T(BLANK).C(0) != '{' && Str(s).T(BLANK).C(LAST) != '}' {
		return
	}

	if children := JSONObjChildren(s); len(children) > 0 {
		// fPln(fName, children) // DEBUG

		(*mFT)[fName] = children //                          *** record path ***

		for _, child := range children {
			if Str(child).HP("[") {
				child = Str(child).RmHeadToLast("]").V() //  *** remove array symbol ***
			}
			nextPath := Str(fName + del + child).T(del).V()
			indices := []int{}
			for range sSpl(nextPath, del) {
				indices = append(indices, 1)
			}

			content, _ := JSONXPathValue(s, nextPath, del, indices...)
			JSONFamilyTree(content, nextPath, del, mFT)
		}
	}
}

// --------------------------------------------------------------------------------------------- //

// IPathToPathIndices :
func IPathToPathIndices(iPath, del string) (string, []int) {
	if iPath == "" {
		return "", nil
	}
	segs, indices := []string{}, []int{}
	for _, s := range sSpl(iPath, del) {
		si := sSpl(s, "#")
		segs, indices = append(segs, si[0]), append(indices, Str(si[1]).ToInt())
	}
	return sJ(segs, del), indices
}

// JSONArrByIPath :
func JSONArrByIPath(s, iPath, del string, mFT *map[string][]string) (arrNames []string, arrCnts []int, nextIPaths []string) {
	path, indices := IPathToPathIndices(iPath, del)
	// fPln("indices:", indices)
	leaves := (*mFT)[path]
	for _, leaf := range leaves {
		// fPln(leaf)
		LEAF := Str(leaf)
		if LEAF.HP("[]") {
			arrName := LEAF.S(2, ALL).V()
			// fPln(arrName)
			_, nArr := JSONXPathValue(s, path+del+arrName, del, append(indices, 1)...)
			// fPln(nArr)

			arrNames = append(arrNames, arrName)
			arrCnts = append(arrCnts, nArr)

			for i := 1; i <= nArr; i++ {
				nextIPath := iPath + del + arrName + fSf("#%d", i)
				// fPln(nextIPath)
				nextIPaths = append(nextIPaths, nextIPath)
			}
		}
	}
	return
}

// JSONWholeArrByIPathByR :
func JSONWholeArrByIPathByR(s, iPath, del, id string, mFT *map[string][]string, mIPathNID *map[string]struct {
	Count int
	ID    string
}) {

	PC(mIPathNID == nil, fEf("result mIPathNID is not initialized"))
	arrNames, arrCnts, subIPaths := JSONArrByIPath(s, iPath, del, mFT)
	PC(len(arrNames) != len(arrCnts), fEf("error in JSONArrByIPath"))

	if len(arrNames) == 0 { //  *** iPath is not array, BUT maybe has children, make new iPaths, then recursive ***

		// fPln(iPath, "no array child")
		r1, _ := Str(iPath).SplitEx(del, "#", "string", "int")
		path := sJ(r1.([]string), del)
		// fPln(path)
		for _, child := range (*mFT)[path] {
			subIPath := iPath + del + child + "#1"
			// fPln("DO MORE THINGS:", subIPath)
			JSONWholeArrByIPathByR(s, subIPath, del, id, mFT, mIPathNID)
		}

	} else {

		nNames := len(arrNames)
		for i := 0; i < nNames; i++ {
			(*mIPathNID)[iPath+del+arrNames[i]] = struct {
				Count int
				ID    string
			}{
				Count: arrCnts[i],
				ID:    id,
			}
		}
		for _, subIPath := range subIPaths {
			JSONWholeArrByIPathByR(s, subIPath, del, id, mFT, mIPathNID)
		}

	}
}

// JSONArrInfo :
func JSONArrInfo(s, xpath, del, id string, mFT *map[string][]string) (*map[string][]string, *map[string]struct {
	Count int
	ID    string
}) {
	if mFT == nil {
		mFT = &map[string][]string{}
		JSONFamilyTree(s, xpath, del, mFT)
	}

	delete(*mFT, "")

	// fPln(" ------------------------------------------------- ")

	keys := GetMapKeys(*mFT).([]string)
	// w.FunSortLess = func(a, b interface{}) bool { //                ** must directly touch package's function variable **
	// 	return sCnt(a.(string), del) < sCnt(b.(string), del)
	// }
	// sortByLess(Strs(keys))

	ok, _, root := IArrSearchOne(Strs(keys), func(i int, a interface{}) (bool, interface{}) {
		return !Str(a.(string)).Contains(del), a
	})
	// fPf("ROOT is <%s>\n", root) //                                  *** DEBUG ***
	PC(!ok, fEf("Invalid path"))

	iRoot := root.(string) + "#1"
	mIPathNID := &map[string]struct {
		Count int
		ID    string
	}{}
	JSONWholeArrByIPathByR(s, iRoot, del, id, mFT, mIPathNID)
	return mFT, mIPathNID
}

// --------------------------------------------------------------------------------------------- //

// JSONMakeObj :
func JSONMakeObj(s, obj, property string, value interface{}, overwrite, mustArr bool) string {
	defer func() { mapJBPos[obj] = Str(s).L() - 1 }()

	property = Str(property).MkQuotes(QDouble).V()

	switch value.(type) {
	case string:
		if !IArrEleIn(Str(value.(string)).C(0), C32s{'{', '['}) {
			value = Str(value.(string)).MkQuotes(QDouble).V()
		}
	}

	s = IF(Str(s).T(BLANK) == "", "{}", s).(string)
	if s == "{}" { //                                                     *** first element ***
		s = fSf(IF(mustArr, `{%s: [%v]}`, `{%s: %v}`).(string), property, value)
		mapJBKids[obj] = append(mapJBKids[obj], property)
		return s
	}

	if start, ok := mapJBPos[obj]; ok { //                                *** existing iPath to add ***

		if IArrEleIn(property, Strs(mapJBKids[obj])) { //                 *** same property, to merge / overwrite ***
			// fPln(property, "did before, merge / overwrite into array")

			for _, find := range Str(s).Indices(property) {
				if Str(s).BracketDepth(BCurly, find) == 1 { //            *** correct insert position ***
					start = find + Str(property).L() + 2 //               *** move to behind ": ", so plus 2 ***
					content, _ := JSONChildValue(s, property)
					// fPln(content)

					if !overwrite { //                                    *** merge into array ***
						oriVL := Str(content).L()
						content = Str(content).RmBrackets(BBox).V()
						content = fSf(`[%s, %v]`, content, value)
						s = Str(s).SegRep(start, start+oriVL, content).V()
						return s
					}

					//                                                    *** overwrite existing content ***
					proval := property + ": " + content
					s = sRep(s, proval, fSf("%s: %v", property, value), 1)
					return s
				}
			}
		}

		//                                                                *** new property ***
		seg := fSf(IF(mustArr, `, %s: [%v]}`, `, %s: %v}`).(string), property, value)
		s = Str(s).SegRep(start, start+1, seg).V()
		mapJBKids[obj] = append(mapJBKids[obj], property)
		return s
	}

	return s
}

// JSONMakeIPath :
func JSONMakeIPath(mIPathObj map[string]string, iPath, property string, value interface{}, mustArr bool) string {
	PC(mIPathObj == nil, fEf("mIPathObj is nil"))
	if content, ok := mIPathObj[iPath]; ok {
		mIPathObj[iPath] = JSONMakeObj(content, iPath, property, value, false, mustArr)
	} else {
		mIPathObj[iPath] = JSONMakeObj("", iPath, property, value, false, mustArr)
	}
	PC(!IsJSON(mIPathObj[iPath]), fEf("<%s>: <%s> is not valid JSON string", iPath, mIPathObj[iPath]))
	return mIPathObj[iPath]
}

// JSONMakeIPathRep :
func JSONMakeIPathRep(mIPathObj map[string]string, del string) string {
	RootKey, keys := "", GetMapKeys(mIPathObj).([]string)
	for _, k := range keys {
		if !Str(k).Contains(del) {
			RootKey = k
			break
		}
	}
	PC(RootKey == "", fEf("ROOT error"))

	sort.SliceStable(keys, func(i, j int) bool {
		return sCnt(keys[i], del) > sCnt(keys[j], del)
	})
	for _, repKey := range keys {
		for k, v := range mIPathObj {
			if Str(v).Contains(repKey) {
				if repValue, ok := mIPathObj[repKey]; ok {
					mIPathObj[k] = sRep(v, "\""+repKey+"\"", repValue, -1)
				}
			}
		}
	}

	// AGAIN:
	// 	for k, v := range mIPathObj {
	// 		for k1, v1 := range mIPathObj {
	// 			if k == k1 {
	// 				continue
	// 			}
	// 			old := Str(v).RmQuotes(QDouble).V()
	// 			if Str(old).Contains(k1) {
	// 				mIPathObj[k] = sRep(v, "\""+k1+"\"", v1, -1)
	// 				goto AGAIN
	// 			}
	// 		}
	// 	}

	return mIPathObj[RootKey]
}

// // JSONMake :
// func JSONMake(s, iPath, del, property string, value interface{}) (string, bool) {
// 	path, indices := IPathToPathIndices(iPath, del)
// 	ss := sSpl(path, del)
// 	exPath, obj := sJ(ss[:len(ss)-1], del), ss[len(ss)-1]
// 	content, _ := JSONXPathValue(s, path, del, indices...)
// 	objstr, ok := JSONMakeObj(content, obj, property, value, false)
// 	root, ok := JSONMakeObj(s, exPath, obj, objstr, true)
// 	return root, ok
// }

// JSONObjectMerge :
func JSONObjectMerge(s, json string) (rst string) {
	PC((s != "" && !IsJSON(s)) || !IsJSON(json), fEf("Error: Invalid JSON string"))

	if s == "" {
		return json
	}

	root1, root2 := JSONFstEle(s), JSONFstEle(json)
	PC(root1 != root2, fEf("Error: Different JSON object"))

	content1, _ := JSONChildValueEx(s, root1, 0)
	content2, _ := JSONChildValueEx(json, root2, 0)
	content1, content2 = Str(content1).TL("["+BLANK).V(), Str(content2).TL("["+BLANK).V()
	content1, content2 = Str(content1).TR("]"+BLANK).V(), Str(content2).TR("]"+BLANK).V()

	rst = fSf(`{ "%s": [ %s, %s ] }`, root1, content1, content2)
	PC(!IsJSON(rst), fEf("Error: Json Merge Result"))
	return
}
