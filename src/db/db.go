package db

import (
	"encoding/json"
	"File"
	"db/Con"
	"fmt"
	"db/API"
)

func Config(path string) {

	var str_config string

	//配置文件是否存在
	if File.Is(path) == true {

		str_config = File.Open(path)
		err := json.Unmarshal([]byte(str_config), &Con.DB)

		//配置文件解析出错
		if err != nil {
			fmt.Println("Load config is not ok")
			return
		}

	} else {
		fmt.Println("File config is not found")
	}

	running()
}

func running(){
	API.Goto()
}