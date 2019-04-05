package gjxy

type JSONTYPE int

var (
	DigStrArr    = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	DigRuneArr   = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	JSONTypeDesc = map[JSONTYPE]string{
		JT_NULL: "Null",
		JT_OBJ:  "Object",
		JT_ARR:  "Array",
		JT_STR:  "String",
		JT_NUM:  "Number",
		JT_BOOL: "Boolean",
		JT_MIX:  "Mix",
		JT_UNK:  "Unknown",
	}
	mapJBPos  = make(map[string]int)
	mapJBKids = make(map[string][]string)
)

const (
	BLANK = " \t\n\r"

	JT_NULL JSONTYPE = 0
	JT_OBJ  JSONTYPE = 1
	JT_ARR  JSONTYPE = 2
	JT_STR  JSONTYPE = 3
	JT_NUM  JSONTYPE = 4
	JT_BOOL JSONTYPE = 5
	JT_MIX  JSONTYPE = 9
	JT_UNK  JSONTYPE = 99
)
