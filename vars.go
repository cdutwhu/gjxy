package gjxy

import (
	"fmt"
	"strings"

	u "github.com/cdutwhu/go-util"
	w "github.com/cdutwhu/go-wrappers"
)

type (
	Str = w.Str
)

const (
	BCurly  = w.BCurly
	BBox    = w.BBox
	QDouble = w.QDouble
	LAST    = w.LAST
	ALL     = w.ALL
)

var (
	IF = u.IF
	PC = u.PanicOnCondition

	fEf  = fmt.Errorf
	fSf  = fmt.Sprintf
	fPln = fmt.Println

	sCtn = strings.Contains
	sCnt = strings.Count
	sSpl = strings.Split
)
