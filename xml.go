package gjxy

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
