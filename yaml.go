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
	return IF(check1 && check2 && check3, true, false).(bool)
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

	if LINE.HP("- ") && !LINE.HS(":") {
		return 1
	}
	if LINE.HP("  - ") && !LINE.HS(":") {
		return 2
	}
	if LINE.HP("    - ") && !LINE.HS(":") {
		return 3
	}
	if LINE.HP("      - ") && !LINE.HS(":") {
		return 4
	}

	for i := 0; i < nLine-1; i++ {
		c, cn := LINE.C(i), LINE.C(i+1)
		if c != ' ' && cn != ' ' {
			return i / 2
		}
	}

	pc(true, fEf("Getting YAMLLevel Error"))
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
		// pc(!Str(value).IsUUID(), fEf("%s is not a valid UUID", value))
		if Str(value).IsUUID() {
			ID, IDByTxt = value, true
		}
	}
	return
}

// YAMLInfo :
func YAMLInfo(yaml, idmark, pathdel string, onlyValues bool) *[]struct {
	// Tag   string
	Path  string
	Value string
	ID    string
} {
	lines := sFF(yaml, func(c rune) bool { return c == '\n' })
	nLine := len(lines)
	rst := make([]struct {
		// Tag   string
		Path  string
		Value string
		ID    string
	}, nLine)

	objGUID := ""
	objid := uuid.New().String()

	tag, value, ID, lvl, _, IDByTxt := YAMLLineInfo(lines[0], idmark, objid)
	rst[0].Path, rst[0].Value, rst[0].ID = tag, value, ID
	if objGUID == "" && IDByTxt {
		objGUID = ID
	}

	lvlPrev, lvlIdxPrev := lvl, make([]int, 64)

	// fPln(lvl)
	// fPln(rst[0])

	for i, l := range lines[1:] {
		i++
		objid = IF(IDByTxt, ID, objid).(string)

		// fPln(i)
		// if i == 93 {
		// 	fPln("debug break")
		// }

		_, value, ID, lvl, _, IDByTxt = YAMLLineInfo(l, idmark, objid)
		rst[i].Value, rst[i].ID = value, ID
		if objGUID == "" && IDByTxt {
			objGUID = ID
		}

		lvlIdx := make([]int, lvl+1)
		copy(lvlIdx, lvlIdxPrev)

		switch {
		case lvl == lvlPrev+1: //                                     *** jump into ***
			lvlIdx[lvl-1] = i - 1
		case lvl == lvlPrev: //                                       *** next sibling ***
			// copy(levelIdx, levelIdxPrev)
		case lvl == lvlPrev-1, lvl == lvlPrev-2, lvl == lvlPrev-3: // *** jump out ***
		default: //                                                   *** incorrect yaml ***
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
			rst[i].Path += (tag + pathdel)
		}
		rst[i].Path = Str(rst[i].Path).RmTailFromLast(pathdel).V()

		// fPln(lvl)
		// fPln(rst[i])

		copy(lvlIdxPrev, lvlIdx)
		lvlPrev = lvl
	}

	//     *** clean up ID, & make onlyValues rst ***

	rstOnlyHasV := []struct {
		// Tag   string
		Path  string
		Value string
		ID    string
	}{}

	for i := 0; i < nLine; i++ {
		if objGUID != "" {
			rst[i].ID = objGUID
		}
		if rst[i].Value != "" {
			rstOnlyHasV = append(rstOnlyHasV, rst[i])
		}
	}

	if onlyValues {
		return &rstOnlyHasV
	}
	return &rst
}

// YAMLGetSplittedLines : iLns : HangingLines List; mBelongto : HangingLine's first value line
func YAMLGetSplittedLines(yaml string) (iLns []int, mBelongto map[int]int) {
	mBelongto = make(map[int]int)
	lines := sFF(yaml, func(c rune) bool { return c == '\n' })
	for i, line := range lines {
		if YAMLIsHangingLine(line) {
			iLns = append(iLns, i)
			for j := i - 1; j >= 0; j-- {
				if !YAMLIsHangingLine(lines[j]) {
					mBelongto[i] = j
					break
				}
			}
		}
	}
	return
}

// YAMLJoinSplittedLines :
func YAMLJoinSplittedLines(yaml string) string {
	iLines, _ := YAMLGetSplittedLines(yaml)
	lines := sFF(yaml, func(c rune) bool { return c == '\n' })
	newLines := []string{}
	for i, line := range lines {
		if IArrEleIn(i, I32s(iLines)) {
			newLines[len(newLines)-1] += " " + Str(line).T(BLANK).V()
		} else {
			newLines = append(newLines, line)
		}
	}
	return sJ(newLines, "\n") + "\n"
}
