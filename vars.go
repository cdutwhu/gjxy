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
	C32s = w.C32s
	I32s = w.I32s
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
	PE = u.PanicOnError

	GetMapKeys    = w.GetMapKeys
	IArrSearchOne = w.IArrSearchOne
	IArrEleIn     = w.IArrEleIn
	IArrIsSameEle = w.IArrIsSameEle
	IArrFoldRep   = w.IArrFoldRep
	Min           = w.Min

	sortByLess = sort.Sort

	fEf  = fmt.Errorf
	fSf  = fmt.Sprintf
	fPln = fmt.Println
	fPf  = fmt.Printf

	// sCtn = strings.Contains
	sCnt = strings.Count
	sSpl = strings.Split
	sJ   = strings.Join
	sRep = strings.Replace
	sFF  = strings.FieldsFunc
)
