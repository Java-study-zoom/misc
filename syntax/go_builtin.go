package syntax

var goBuiltInTypes = []string{
	"int",
	"uint",
	"int64",
	"uint64",
	"int32",
	"uint32",
	"int16",
	"uint16",
	"int8",
	"uint8",
	"byte",
	"rune",
	"error",
	"string",
	"float32",
	"float64",
	"complex64",
	"complex128",
	"uintptr",
	"bool",
	"map",
	"true",
	"false",
	"nil",
	"iota",
}

var goBuiltInFuncs = []string{
	"len",
	"cap",
	"close",
	"complex",
	"delete",
	"imag",
	"panic",
	"print",
	"println",
	"real",
	"recover",
	"make",
	"append",
	"new",
	"copy",
}

var (
	goBuiltInFuncMap = makeMap(goBuiltInFuncs)
	goBuiltInTypeMap = makeMap(goBuiltInTypes)
)
