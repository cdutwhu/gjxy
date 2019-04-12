package gjxy

import "github.com/google/uuid"

// YAMLIsPath :
func YAMLIsPath(line string) bool {
	LINE := Str(line).T(BLANK)
	check1 := LINE.C(LAST) == ':'
	check2 := !LINE.Contains(": ")
	return check1 && check2
}

// YAMLIsValueLine :
func YAMLIsValueLine(line string) bool {
	return !YAMLIsPath(line)
}

// YAMLIsHangingLine :
func YAMLIsHangingLine(line string) bool {
	LINE := Str(line).T(BLANK)
	check1 := !LINE.HP("- ")
	check2 := !LINE.Contains(": ")
	check3 := LINE.C(LAST) != ':'
	if check1 && check2 && check3 {
		return true
	}
	return false
}

// YAMLTag :
func YAMLTag(line string) string {
	K := Str(line).T(BLANK)
	if YAMLIsValueLine(line) {
		K, _ = K.KeyValuePair(": ", "~", "~", true, true)
		if K.HP("- ") {
			K = K.S(2, ALL).RmQuotes(QDouble)
		}
		return K.V()
	}
	K = K.S(0, ALL-1)
	return K.TL(BLANK).RmQuotes(QDouble).V()
}

// YAMLValue :
func YAMLValue(line string) (value string, isArr bool) {
	if YAMLIsValueLine(line) {
		LINE := Str(line).T(BLANK)
		if p := LINE.Idx(": "); p >= 0 { //     ***  Normal 'Sub: Obj' line ****
			value := LINE.S(p+2, ALL)
			value = IF(value != `""`, value.RmQuotes(QDouble), value).(Str)
			return value.V(), false
		}
		if p := LINE.Idx("- "); p >= 0 { //     *** Pure Array Element '- Obj' line ***
			value := LINE.S(p+2, ALL)
			value = IF(value != `""`, value.RmQuotes(QDouble), value).(Str)
			return value.V(), true
		}
	}
	return "", false //                         *** Pure One Path Section ***
}

// ************ DO NOT TRIM LINE'S BLANK START HERE ************

// YAMLLevel :
func YAMLLevel(line string) int {
	LINE := Str(line)
	nLine := LINE.L()
	if nLine == 3 && LINE.HP("- ") {
		return 1
	}
	for i := 0; i < nLine-1; i++ {
		c, cn := LINE.C(i), LINE.C(i+1)
		if c != ' ' && cn != ' ' {
			return i / 2
		}
	}
	return -1
}

// YAMLLineInfo :
func YAMLLineInfo(line, idmark, id string) (tag, value, ID string, level int, isArr, IDByTxt bool) {
	LINE := Str(line)
	tag = YAMLTag(line)
	value, isArr = YAMLValue(line)
	level = YAMLLevel(line)
	ID, IDByTxt = id, false
	if LINE.HP(idmark+": ") || LINE.Contains(" "+idmark+": ") || LINE.Contains(" -"+idmark+": ") {
		PC(!Str(value).IsUUID(), fEf("%s is not a valid UUID", value))
		ID, IDByTxt = value, true
	}
	return
}

// YAMLInfo :
func YAMLInfo(yaml, idmark, pathdel string, onlyValues bool) *[]struct {
	// tag   string
	path  string
	value string
	ID    string
} {
	lines := sFF(yaml, func(c rune) bool { return c == '\n' })
	nLine := len(lines)
	rst := make([]struct {
		// tag   string
		path  string
		value string
		ID    string
	}, nLine)

	objGUID := ""
	objid := uuid.New().String()

	tag, value, ID, lvl, _, IDByTxt := YAMLLineInfo(lines[0], idmark, objid)
	rst[0].path, rst[0].value, rst[0].ID = tag, value, ID
	if objGUID == "" && IDByTxt {
		objGUID = ID
	}

	lvlPrev, lvlIdxPrev := lvl, make([]int, 64)	

	// fPln(lvl)
	// fPln(rst[0])

	for i, l := range lines[1:] {
		i++
		objid = IF(IDByTxt, ID, objid).(string)

		// if i == 4 {
		// 	fPln("break")
		// }

		_, value, ID, lvl, _, IDByTxt = YAMLLineInfo(l, idmark, objid)
		rst[i].value, rst[i].ID = value, ID
		if objGUID == "" && IDByTxt {
			objGUID = ID
		}

		lvlIdx := make([]int, lvl+1)
		copy(lvlIdx, lvlIdxPrev)

		switch {
		case lvl == lvlPrev+1: //                                         *** jump into ***
			lvlIdx[lvl-1] = i - 1
		case lvl == lvlPrev: //                                           *** next sibling ***
			// copy(levelIdx, levelIdxPrev)
		case lvl == lvlPrev-1, lvl == lvlPrev-2, lvl == lvlPrev-3: //     *** jump out ***
		default: //                                                       *** incorrect yaml ***
		}
		lvlIdx[lvl] = i

		for _, p := range lvlIdx {
			tag := YAMLTag(lines[p]) //                   *** parent, grand-parent... tag ***
			if Str(tag).HP("- ") {
				tag = Str(tag).S(2, ALL).V()
			}
			if tag == "" {
				continue
			}
			rst[i].path += (tag + pathdel)
		}
		rst[i].path = Str(rst[i].path).RmTailFromLast(pathdel).V()

		// fPln(lvl)
		// fPln(rst[i])

		copy(lvlIdxPrev, lvlIdx)
		lvlPrev = lvl
	}

	//     *** clean up ID, & make onlyValues rst ***

	rstOnlyHasV := []struct {
		// tag   string
		path  string
		value string
		ID    string
	}{}

	for i := 0; i < nLine; i++ {
		if objGUID != "" {
			rst[i].ID = objGUID
		}
		if rst[i].value != "" {
			rstOnlyHasV = append(rstOnlyHasV, rst[i])
		}
	}

	if onlyValues {
		return &rstOnlyHasV
	}
	return &rst
}
