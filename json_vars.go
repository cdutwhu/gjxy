package gjxy

type JTYPE int

var (
	DigStrArr    = [...]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-"}
	DigRuneArr   = [...]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-'}
	JSONTypeDesc = map[JTYPE]string{
		J_NULL: "Null",
		J_OBJ:  "Object",
		J_ARR:  "Array",
		J_STR:  "String",
		J_NUM:  "Number",
		J_BOOL: "Boolean",
		J_MIX:  "Mix",
		J_UNK:  "Unknown",
	}
	mapJBPos  = make(map[string]int)
	mapJBKids = make(map[string][]string)
)

const (
	BLANK = " \t\n\r"

	J_NULL JTYPE = 0
	J_OBJ  JTYPE = 1
	J_ARR  JTYPE = 2
	J_STR  JTYPE = 3
	J_NUM  JTYPE = 4
	J_BOOL JTYPE = 5
	J_MIX  JTYPE = 9
	J_UNK  JTYPE = 99
)
