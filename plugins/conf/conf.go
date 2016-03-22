// Package conf 从配置文件（JSON格式）中加载配置。
package conf

import (
	"encoding/json"
	"io/ioutil"
)

// Options 配置内容
var Options struct {
	Port             int        // 端口
	AutoOpenBrowser  bool       // 是否自动打开浏览器
	DocumentsService []struct { // 请根据特定项目的配置文件结构进行配置
		Project string // 项目标题
		Path    string // 项目文档目录
	}
}

// Load load configuration from .json file
func Load(filepath string) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(content, &Options); err != nil {
		panic(err)
	}
}
