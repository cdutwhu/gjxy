package gjxy

// YAMLIsPath :
func YAMLIsPath(line string) bool {
	return Str(line).C(LAST) == ':' && !sCtn(line, ": ")
}

// YAMLIsValueLine :
func YAMLIsValueLine(line string) bool {
	return !YAMLIsPath(line)
}

// YAMLIsHangingLine :
func YAMLIsHangingLine(line string) bool {
	LINE := Str(line).T(BLANK)
	if !LINE.HP("- ") && !sCtn(LINE.V(), ": ") && LINE.C(LAST) != ':' {
		return true
	}
	return false
}

// YAMLTag :
func YAMLTag(line string) string {
	K := Str("")
	if YAMLIsValueLine(line) {
		K, _ = Str(line).KeyValuePair(": ", "~", "~", true, true)
		if K.HP("- ") {
			K = K.S(2, ALL).RmQuotes(QDouble)
		}
		return K.V()
	}
	return K.S(0, ALL-1).TL(BLANK).RmQuotes(QDouble).V()
}

// YAMLValue :
func YAMLValue(line string) (value string, isArr bool) {

	if YAMLIsValueLine(line) {
		LINE := Str(line)
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
