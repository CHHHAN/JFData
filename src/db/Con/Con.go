package Con

var DB mains
var Site set

type mains struct {
	ID string
	Password string
	Path string
	Port string
}

type set struct {
	Func string
	POST string
}

type Login struct {
	ID string
	Password string
}

type FuncDB struct {
	Name string
}

type FuncTB struct {
	Name string
	Table string
}

//添加新记录
type FuncVA struct {
	Name string
	Table string
	V []string
	A []string
}

//修改指定记录
type FunceVA struct {
	Name string
	Table string
	I int
	V []string
	A []string
}

//返回指定记录
type FuncI struct {
	Name string
	Table string
	I int
}