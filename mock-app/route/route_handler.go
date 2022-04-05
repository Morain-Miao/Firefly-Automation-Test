package route

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mock-app/common"
	"mock-app/database"
	"mock-app/service"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const prefix = "/mock/"

var (
	once   sync.Once
	router *gin.Engine
)

// GetRouter 单例获取路由管理者
func GetRouter() *gin.Engine {
	once.Do(func() {
		router = gin.Default()
	})
	return router
}

// 初始化路由信息
func init() {
	routes, _ := database.QueryAllRoutes()
	for _, route := range routes {
		// 永不过期
		common.AddCache(route.HttpMethod+":"+route.Host+route.Path, route.Id, -1*time.Minute)
		AddRoute(prefix+route.Host+route.Path, route.HttpMethod)
	}
}

// AddRoute 在gin的路由信息上下文添加路由
func AddRoute(url string, httpMethod string) {
	router := GetRouter()
	router.Handle(httpMethod, url, service.MockRequestHandler)
}

// AddMockRouteHandler 添加mockApi
// 请求参数：
// host 请求主机头 不带http、https前缀
// port 端口号
// path 主机头后面的请求路径
// http_method 请求方法 例如：GET
// http_template_id http模板的uuid
func AddMockRouteHandler(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	for i := 0; i < n; i += 1 {
		if buf[i] == ' ' || buf[i] == '\x00' {
			buf = append(buf[:i], buf[(i+1):]...)
			n -= 1
		}
	}

	responseMap := make(map[string]string)

	err := json.Unmarshal(buf[:n], &responseMap)
	if err != nil {
		panic(common.JsonDecodeException("请求体Json解码错误"))
	}
	fmt.Println(responseMap)

	port, err := strconv.ParseInt(responseMap["port"], 10, 64)

	route := database.Route{
		Host:           responseMap["host"],
		Port:           port,
		Path:           responseMap["path"],
		HttpMethod:     responseMap["http_method"],
		HttpTemplateId: responseMap["http_template_id"],
	}

	err = checkRoute(route)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"data":    false,
			"message": err.Error(),
		})
		return
	}

	id, rowsAffected, err := route.Insert()
	if rowsAffected == 0 {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    false,
			"message": "路由信息插入错误，影响行数为0！",
		})
		return
	}
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"data":    false,
			"message": err.Error(),
		})
		return
	}
	common.AddCache(route.HttpMethod+":"+route.Host+route.Path, id, -1*time.Minute)
	AddRoute(prefix+route.Host+route.Path, route.HttpMethod)
	c.JSON(200, gin.H{
		"code":    200,
		"data":    true,
		"message": "添加路由成功",
	})
}

// DelMockRouteHandler 删除mock中的api  url = host+path 去掉http前缀
func DelMockRouteHandler(c *gin.Context) {
	url := c.Query("url")
	method := c.Query("http_method")
	router := GetRouter()
	router.DELETE(prefix+url, service.MockRequestHandler)
	common.DeleteCache(method + ":" + url)
	c.JSON(200, gin.H{
		"code":    200,
		"data":    true,
		"message": "ok",
	})
}

func checkRoute(route database.Route) error {
	cache, found := common.GetCache(route.HttpMethod + ":" + route.Host + route.Path)
	if found && cache != nil {
		return fmt.Errorf("route already registered for path %s", route.Host+route.Path)
	}
	host := route.Host
	// (.*)(?=:)
	if strings.ContainsAny(route.Host, ":") {
		host = strings.Split(route.Host, ":")[0]
	}
	count, err := database.QueryHostRoutes(host)
	if count > 0 {
		return fmt.Errorf("new path '%s' port conflicts with existing wildcard in existing prefix", route.Host+route.Path)
	}
	if err != nil {
		return fmt.Errorf("check sql execute error %s, Please contact the operation admin", err)
	}

	return nil
}
