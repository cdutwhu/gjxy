package gjxy

// XMLTag : get the xml string's first tag
func XMLTag(xml string) string {
	XML := Str(xml).T(BLANK)
	PC(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Not a valid XML section"))
	return XML.S(XML.LIdx("</")+2, ALL-1).V()
}

// XMLTagEle : Looking for the first tag's xml string
func XMLTagEle(xml, tag string) (string, int, int) {
	XML := Str(xml).T(BLANK)
	PC(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Invalid XML"))

	XML = Str(xml) //                                                   *** we have to return original position, so use original xml ***
	s, sa := XML.Idx(fSf("<%s>", tag)), XML.Idx(fSf("<%s ", tag))
	if s < 0 && sa >= 0 {
		s = sa
	} else if s >= 0 && sa < 0 {
		// s = s
	} else if s >= 0 && sa >= 0 {
		min, _ := Min(I32s{s, sa}, "")
		s = min.(int)
	} else if s < 0 && sa < 0 {
		return "", -1, -1
	}

	if eR := XML.S(s, ALL).Idx(fSf("</%s>", tag)); eR > 0 {
		sNext := s + eR + Str(tag).L() + 3
		return XML.S(s, sNext).V(), s, sNext
	}

	panic("Invalid XML")
}

// XMLTagEleEx : idx from 1
func XMLTagEleEx(xml, tag string, idx int) string {
	esum := 0
	for i := 1; i <= idx; i++ {
		XML := Str(xml).S(esum, ALL)
		ele, _, e := XMLTagEle(XML.V(), tag)
		if e == -1 {
			return ""
		}
		if i == idx {
			return ele
		}
		esum += e
	}
	panic("Should not be here!")
}

// XMLXPathEle :
func XMLXPathEle(xml, xpath, del string, indices ...int) (ele string) {
	PC(xpath == "", fEf("At least one path must be provided"))

	segs := sSpl(xpath, del)
	PC(len(segs) != len(indices), fEf("path & seg's index count not match"))

	for i, seg := range segs {
		xml = IF(ele != "", ele, xml).(string)
		ele = XMLTagEleEx(xml, seg, indices[i])
	}
	return
}

// XMLAttributes is (ONLY LIKE  <SchoolInfo RefId="D3F5B90C-D85D-4728-8C6F-0D606070606C" Type="LGL">)
func XMLAttributes(xmlele, attrPF string) (attributes, attriValues []string) { //  *** 'map' may cause mis-order, so use slice
	XMLELE := Str(xmlele).T(BLANK)
	PC(XMLELE == "" || XMLELE.C(0) != '<' || XMLELE.C(LAST) != '>', fEf("Not a valid XML section"))

	tag := Str(XMLTag(xmlele))
	if eol := XMLELE.Idx(`">`) + 1; XMLELE.C(tag.L()+1) == ' ' && eol > tag.L() { //     *** has attributes
		kvs := sFF(XMLELE.S(tag.L()+2, eol).V(), func(c rune) bool { return c == ' ' })
		for _, kv := range kvs {
			kvstrs := sFF(kv, func(c rune) bool { return c == '=' })
			attributes = append(attributes, (attrPF + kvstrs[0])) //                     *** mark '-' before attribute for differentiating child
			attriValues = append(attriValues, Str(kvstrs[1]).RmQuotes(QDouble).V())
		}
	}
	return attributes, attriValues
}

// XMLChildren : (NOT search grandchildren)
func XMLChildren(xmlele string, fNArr bool) (children []string) {
	XMLELE := Str(xmlele).T(BLANK)
	PC(XMLELE == "" || XMLELE.C(0) != '<' || XMLELE.C(LAST) != '>', fEf("Invalid XML section"))

	L := XMLELE.L()
	skip, childpos, level, inflag := false, []int{}, 0, false

	for i := 0; i < L; i++ {
		c := XMLELE.C(i)

		if c == '<' && XMLELE.S(i, i+4) == "<!--" {
			skip = true
		}
		if c == '>' && XMLELE.S(i-2, i+1) == "-->" {
			skip = false
		}
		if skip {
			continue
		}

		if c == '<' && XMLELE.C(i+1) != '/' {
			level++
		}
		if c == '<' && XMLELE.C(i+1) == '/' {
			level--
			if level == 1 {
				inflag = false
			}
		}

		if level == 2 {
			if !inflag {
				childpos, inflag = append(childpos, i+1), true
			}
		}
	}

	if len(childpos) == 0 {
		return
	}

	for _, p := range childpos {
		pe, peA := XMLELE.S(p, ALL).Idx(">"), XMLELE.S(p, ALL).Idx(" ")
		pe = IF(peA > 0 && peA < pe, peA, pe).(int)
		child := XMLELE.S(p, p+pe)
		children = append(children, child.V())
	}

	children = IArrFoldRep(Strs(children), IF(fNArr, "[n]", "[]").(string)).([]string)
	return
}

// XMLFamilyTree : ******************************************************************************************
func XMLFamilyTree(xml, fName, del string, mapFT *map[string][]string) {
	PC(mapFT == nil, fEf("FamilyTree return map is not initialised !"))
	XML := Str(xml).T(BLANK)
	PC(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Invalid XML section"))

	fName = IF(fName == "", XMLTag(xml), fName).(string)
	if children := XMLChildren(xml, false); len(children) > 0 {
		// fPln(tag, children)

		(*mapFT)[fName] = children //                           *** record path ***

		for _, child := range children {
			if Str(child).HP("[") {				
				child = Str(child).RmHeadToLast("]").V() //     *** remove array symbol ***
			}
			nextPath := Str(fName + del + child).T(del).V()
			subxml := XMLTagEleEx(xml, child, 1)
			XMLFamilyTree(subxml, nextPath, del, mapFT)
		}
	}
}

// // XMLYieldArrInfo :
// func XMLYieldArrInfo(xmlstr string, ids, objs []string, mapkeyprefix, pathDel, childDel string, eleObjIDArrcnts *[]pathIDn) {
// 	if len(mapkeyprefix) > 0 {
// 		mapkeyprefix += pathDel
// 	}
// 	for i, obj := range objs {
// 		curPath := mapkeyprefix + obj

// 		xmlele := XMLEleStrByTag(xmlstr, obj, 1)
// 		uids, children, _, arrCnt := XMLFindChildren(xmlele, ids[i], childDel) /* uniform ids, children */
// 		attributes, _, _ := XMLFindAttributes(xmlele, childDel)                /* attributes */

// 		/* array children info */
// 		if arrCnt > 0 {
// 			(*eleObjIDArrcnts) = append((*eleObjIDArrcnts), pathIDn{arrPath: curPath + pathDel + children[0], objID: ids[i], arrCnt: arrCnt})
// 		}

// 		if len(children) == 0 && len(attributes) == 0 { /* attributes */
// 			continue
// 		} else {
// 			XMLYieldArrInfo(xmlele, uids, children, curPath, pathDel, childDel, eleObjIDArrcnts) /* recursive */
// 		}
// 	}
// }

/**********************************************************************************************************************************/

// XMLSegPos : level from 1, index from 1                                         &
func XMLSegPos(s Str, level, index int) (tag, str string, left, right int) {
	markS, markE1, markE2, markE3 := '<', '<', '/', '>'
	curLevel, curIndex, To := 0, 0, s.L()-1

	found := false
	i := 0
	for _, c := range s {
		if i < To {
			curLevel = IF(c == markS && s.C(i+1) != markE2, curLevel+1, curLevel).(int)
			curLevel = IF(c == markE1 && s.C(i+1) == markE2, curLevel-1, curLevel).(int)
			if curLevel == level && c == markS && s.C(i+1) != markE2 {
				left = i
			}
			if curLevel == level-1 && c == markE1 && s.C(i+1) == markE2 {
				right = i
				curIndex++
				if curIndex == index {
					found = true
					break
				}
			}
		}
		i++
	}

	if !found {
		return "", "", 0, 0
	}

	tagendRel := s.S(left+1, right).Idx(" ") // when tag has attribute(s)
	if tagendRel == -1 {
		tagendRel = s.S(left+1, right).Idx(string(markE3))
	}
	PC(tagendRel == -1, fEf("xml error"))

	tag = s.S(left+1, left+1+tagendRel).V()
	right += Str(tag).L() + 2
	return tag, s.S(left, right+1).V(), left, right
}

// XMLSegsCount : only count top level                                            &
func XMLSegsCount(s Str) (count int) {
	markS, markE1, markE2 := '<', '<', '/'

	level, inflag, To := 0, false, s.L()-1
	i := 0
	for _, c := range s {
		if i < To {
			if c == markS && s.C(i+1) != markE2 {
				level++
			}
			if c == markE1 && s.C(i+1) == markE2 {
				level--
				if level == 0 {
					inflag = false
				}
			}
			if level == 1 {
				if !inflag {
					count++
					inflag = true
				}
			}
		}
		i++
	}
	return count
}
