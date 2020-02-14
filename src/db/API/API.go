package API

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"db/Con"
	"STR"
	"reflect"
	"errors"
	"File"
)

//函数注册
var API = map[string]interface{} {

	"login":login,

	"newDB":newDB,
	"howDB":howDB,
	"delDB":delDB,

	"newTB":newTB,
	"howTB":howTB,
	"delTB":delTB,

	"addVA":addVA,
	"excVA":excVA,
	"howVA":howVA,
	"toVA":toVA,
	"toV":toV,
	"delVA":delVA,

}

func Goto(){

	http.HandleFunc("/", Handler)
	fmt.Println("Datebase is running")
	http.ListenAndServe(":" + Con.DB.Port, nil)

}

func setCode(Code int, Msg interface{}) string {

	type Ret struct {
		Code int
		Msg interface{}
	}

	var setReturn Ret

	setReturn.Code = Code
	setReturn.Msg = Msg

	setReturn_byte, _ := json.Marshal(setReturn)

	return string(setReturn_byte)

}

func Call(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func Handler(w http.ResponseWriter, r *http.Request){

	//获取执行参数
	Con.Site.Func = r.RequestURI

	//是否为POST
	if r.Method == "GET" {
		setCode(500, "The Method is not")
	}

	//存储 POST 数据
	str, _ := ioutil.ReadAll(r.Body)
	Con.Site.POST = string(str)

	//调出 PIN 和函数
	pare := STR.Cut(Con.Site.Func,"/")
	strs := len(pare)

	//缺少 GET 参数
	if strs != 2 {
		fmt.Fprintf(w, setCode(500, "The URL is not"))
		return
	}

	//存储与检查 函数
	Con.Site.Func = pare[1]

	if Con.Site.Func == "" {
		fmt.Fprintf(w, setCode(500, "The Func is null"))
		return
	}

	//函数错误
	if API[Con.Site.Func] == nil {
		fmt.Fprintf(w, setCode(404, "The Func is not found"))
		return
	}

	//API 返回 json
	if result, err := Call(API, Con.Site.Func); err == nil {
		for _, r := range result {
			str := r.String()
			fmt.Println(str)
			w.Header().Set("Content-type","application/x-www-form-urlencoded")
			fmt.Fprintf(w, str)
			return
		}
	}

}

//登录
func login() string {

	var POST Con.Login

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	if POST.ID == Con.DB.ID && POST.Password == Con.DB.Password {
		return setCode(200, "OK")
	}

	return setCode(500,"NO")
}

//新建数据库 ok
func newDB() string {

	var POST Con.FuncDB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == ""{
		return setCode(502, "The POST json string is null")
	}

	//数据库存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == true {
		return setCode(302, "The POST Name is on")
	} else {
		File.NewDir(Con.DB.Path + "/" + POST.Name)
	}

	return setCode(200, "OK")

}

//查询数据库 ok
func howDB() string {

	fileList := File.ListDir(Con.DB.Path)
	hows := len(fileList)
	var database []string
	for i := 0; i < hows; i++ {
		database = append(database,fileList[i].Name())
	}

	return setCode(200, database)

}

//删除数据库
func delDB() string {

	var POST Con.FuncDB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == ""{
		return setCode(502, "The POST json string is null")
	}

	//数据库存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		File.DelDir(Con.DB.Path + "/" + POST.Name)
	}

	return setCode(200, "OK")

}


//新建数据表 ok
func newTB() string {

	var POST Con.FuncTB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		if File.Is(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json") == true {
			return setCode(302, "The POST Table is not found")
		} else {
			File.New(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json","")
		}
	}

	return setCode(200, "OK")

}

//查询数据表 ok
func howTB() string {

	var POST Con.FuncDB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" {
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is on")
	}

	//返回所有数据表
	fileList := File.ListDir(Con.DB.Path + "/" + POST.Name)
	hows := len(fileList)
	var database []string
	for i := 0; i < hows; i++ {
		filename := fileList[i].Name()
		lang := len(filename)
		if filename[lang-5:lang] == ".json" {
			database = append(database,filename[0:lang-5])
		}
	}

	return setCode(200, database)

}

//删除数据表
func delTB() string {

	var POST Con.FuncTB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		if File.Is(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json") == false {
			return setCode(404, "The POST Table is not found")
		} else {
			File.DelDir(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
		}
	}

	return setCode(200, "OK")

}


//添加新记录 ok *bug 需要完整的字段哦～
func addVA() string {
	var POST Con.FuncVA

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == "" || len(POST.V) == 0 || len(POST.A) == 0 || len(POST.V) != len(POST.A){
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		//数据表不存在
		if File.Is(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json") == false {
			return setCode(404, "The POST Table is not found")
		}
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	var excV []string
	var addV []string

	var excA []string
	var addA []string

	howT := len(dbMap)
	howV := len(POST.V)

	for i := 0; i < howT ; i++ {
		for c := 0; c < howV ; c++ {
			_ , ok := dbMap[i][POST.V[c]]
			fmt.Println(ok)
			//字段存在，该字段需要修改
			if ok == true {
				fmt.Println("需要修改：" + POST.V[c] + "：" + POST.A[c])
				excV = append(excV,POST.V[c])
				excA = append(excA,POST.A[c])
			} else {
			//字段不存在，该字段需要新建
				fmt.Println("新建：" + POST.V[c] + "：" + POST.A[c])
				addV = append(addV,POST.V[c])
				addA = append(addA,POST.A[c])
			}
		}
	}

	dbMap[howT] = make(map[string]string)
	howexcV := len(excV)
	howaddV := len(addV)

	if howexcV > 0 {
		//对于需要修改的字段
		for v := 0; v < howexcV ; v++ {
			//dbMap[0]["id"]:"test"
			dbMap[howT][excV[v]] = excA[v]
		}
	}

	if howaddV > 0 {
		//对于需要新建的字段
		for t := 0; t < howT + 1; t++ {
			for v := 0; v < howaddV ; v++ {
				//dbMap[0]["id"]:"test"
				dbMap[t][addV[v]] = addA[v]
			}
		}
	}

	if howT == 0 {
		for i := 0; i < howV ; i++ {
			//dbMap[0]["id"]:"test"
			dbMap[0][POST.V[i]] = POST.A[i]
		}

	}

	tables_byte, err := json.Marshal(dbMap)
	File.New(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json", string(tables_byte))

	return setCode(200, "OK")
}

//修改指定记录 ok
func excVA() string {
	var POST Con.FunceVA

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == "" || len(POST.V) == 0 || len(POST.A) == 0 || len(POST.V) != len(POST.A){
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		//数据表不存在
		if File.Is(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json") == false {
			return setCode(404, "The POST Table is not found")
		}
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	howV := len(POST.V)

	//检查对应记录是否存在
	_ , ok := dbMap[POST.I]
	//字段存在，该字段需要修改
	if ok == false {
		return setCode(404, "The POST Table is not found")
	}

	//检查字段是否存在
	for c := 0; c < howV ; c++ {
		_ , ok := dbMap[POST.I][POST.V[c]]
		//字段存在，该字段需要修改
		if ok == false {
			return setCode(404, "The POST V is not found")
		}
	}

	for v := 0; v < howV ; v++ {
		//dbMap[0]["id"]:"test"
		dbMap[POST.I][POST.V[v]] = POST.A[v]
	}

	tables_byte, err := json.Marshal(dbMap)
	File.New(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json", string(tables_byte))

	return setCode(200, "OK")
}

//返回所有记录 ok
func howVA() string {

	var POST Con.FuncTB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	return setCode(200, dbMap)
}

//返回指定记录 ok
func toVA() string {

	var POST Con.FuncI

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	return setCode(200, dbMap[POST.I])
}

//返回所有字段 ok
func toV() string {

	var POST Con.FuncTB

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	var tableV []string
	for V, _ :=  range dbMap[0]{
		tableV = append(tableV,V)
	}

	return setCode(200, tableV)
}

//删除指定记录 ok
func delVA() string {

	var POST Con.FuncI

	err := json.Unmarshal([]byte(Con.Site.POST), &POST)

	//Json 错误 - 即 传参错误
	if err != nil {
		return setCode(502, "The POST json is not")
	}

	//创建的数据库名不符合规范
	if POST.Name == "" || POST.Table == ""{
		return setCode(502, "The POST json string is null")
	}

	//数据库不存在？
	if File.Is(Con.DB.Path + "/" + POST.Name) == false {
		return setCode(404, "The POST Name is not found")
	} else {
		//数据表不存在
		if File.Is(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json") == false {
			return setCode(404, "The POST Table is not found")
		}
	}

	dbMap := make(map[int]map[string]string)
	dbstring := File.Open(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json")
	json.Unmarshal([]byte(dbstring), &dbMap)

	//检查对应记录是否存在
	_ , ok := dbMap[POST.I]
	//字段存在，该字段需要修改
	if ok == false {
		return setCode(404, "The POST Table is not found")
	}

	delete(dbMap, POST.I)
	tables_byte, err := json.Marshal(dbMap)
	File.New(Con.DB.Path + "/" + POST.Name + "/" + POST.Table + ".json", string(tables_byte))

	return setCode(200, "OK")
}