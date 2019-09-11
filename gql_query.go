package gjxy

import (
	"errors"
)

// RmQryCmts :
func RmQryCmts(qry string) string {
	i := 0
	qrybuf := make([]rune, len(qry)+1)
	flagAdd := true
	for _, uc := range qry {
		flagAdd = matchAssign(uc, '#', '\n', false, true, flagAdd).(bool)
		if flagAdd {
			qrybuf[i] = uc
			i++
		}
	}
	return string(qrybuf)
}

// Get1stObjInQry :
func Get1stObjInQry(qry string) (obj string, err error) {
	Sq := Str(RmQryCmts(qry)).T(BLANK)
	_, curly1, _ := Sq.BracketsPos(BCurly, 1, 1)
	_, curly2, _ := Sq.BracketsPos(BCurly, 2, 1)
	_, round1, _ := Sq.BracketsPos(BRound, 1, 1)
	_, round2, _ := Sq.BracketsPos(BRound, 1, 2)

	if curly1 == -1 || curly2 == -1 {
		return "", errors.New("NOT A VALID QUERY TEXT")
	}

	to := -1
	if round1 > 0 { //    *** has query params ***
		to = trueAssign(!Sq.HP("query"), Sq.HP("query") && round2 > 0, round1, round2, to).(int)
	} else { //           *** no query params ***
		to = curly2
	}

	field := Sq.S(curly1+1, to).T(BLANK).V()
	if segs := sSpl(field, ":"); len(segs) == 2 {
		field = segs[1]
	}
	field = Str(field).T(BLANK).V()
	return field, nil
}
