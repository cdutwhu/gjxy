package gjxy

import (
	"fmt"
	"sort"
	"strings"

	u "github.com/cdutwhu/go-util"
	w "github.com/cdutwhu/go-wrappers"
)

type (
	Str  = w.Str
	Strs = w.Strs
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

	GetMapKeys = w.GetMapKeys

	sortByLess = sort.Sort

	fEf  = fmt.Errorf
	fSf  = fmt.Sprintf
	fPln = fmt.Println

	sCtn = strings.Contains
	sCnt = strings.Count
	sSpl = strings.Split
	sJ   = strings.Join	
)
