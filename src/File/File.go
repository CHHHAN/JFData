package File

import (
	"os"
	"io/ioutil"
	"fmt"
)

func Is(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
	return true
	}

	if os.IsNotExist(err) {
	return false
	}

	return false
}

func New(path string, str string) bool {
	fmt.Println("new")
	if Is(path) == true {
		fmt.Println("true")
		err := os.Remove(path)
		fmt.Println(err)
	}

	f,err := os.Create(path)
	defer f.Close()

	if err !=nil {
		//创建文件不成功
		return false
	} else {
		_,err=f.Write([]byte(str))
		return true
	}

}

func Open(path string) string {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		//读取文件失败
		return ""
	}

	return string(b)

}

func NewDir(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false
	} else {
		return true
	}
}

func DelDir(path string)bool{
	os.RemoveAll(path)
	return true
}

func ListDir(path string) []os.FileInfo {

	dir, err := ioutil.ReadDir(path)

	if err != nil {
		return nil
	}

	return dir
}
