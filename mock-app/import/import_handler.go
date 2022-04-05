package _import

import (
	"github.com/gin-gonic/gin"
	"log"
	"mock-app/common"
	"path/filepath"
	"plugin"
	"sync"
)

// ImportPlugins 我是存储插件的容器
var (
	ImportPlugins sync.Map
	Dir           = filepath.Dir("./resource/plugins")
)

func init() {
	files, _ := common.GetAllFiles(Dir)
	for _, file := range files {
		p, err := plugin.Open(file)
		if err != nil {
			log.Println(file + "插件加载错误！")
			continue
		}
		ImportPlugins.Store(filepath.Base(file), p)
	}

}

func ImportHandler(c *gin.Context) {
	pluginName := c.Param("plugin")
	p, _ := ImportPlugins.Load(pluginName + ".so")
	if p == nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    false,
			"message": "插件不存在，请先注册插件！",
		})
		return
	}
	//p, err := plugin.Open("testplugin.so")
	f, err := p.(*plugin.Plugin).Lookup("Hello")
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    false,
			"message": "自定义插件中不存在XXX方法！",
		})
		return
	}
	f.(func())()
}

func GetImportPlugins(c *gin.Context) {

}

func DelImportPlugin(c *gin.Context) {

}

func UpdateImportPlugin(c *gin.Context) {

}
