package main

import (
	"api_system_iris/apps"
	"api_system_iris/config"
	"api_system_iris/utils"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/spf13/viper"
	"strings"
)

func newApp() (app *iris.Application) {
	app = iris.New()
	loglevel := viper.GetString("logLevel")
	if len(loglevel) == 0 {
		loglevel = "info"
	}
	app.Logger().SetLevel(loglevel)
	c := logger.Config{
		Status:   true,
		IP:       true,
		Method:   true,
		Path:     true,
		Query:    true,
		Columns:  false,
		MessageHeaderKeys: []string{"User-Agent"},
	}
	var excludeExtensions = [...]string{".js", ".css", ".jpg", ".png", ".ico", ".svg",}
	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _,ext := range excludeExtensions {
			if strings.HasSuffix(path,ext) {
				return true
			}
		}
		return false
	})
	logger := logger.New(c)
	//注册日志记录中间件
	app.Use(recover.New())
	app.Use(logger)
	//注册路由
	apps.RegRouter(app,logger)

	return
}

func main() {
	fmt.Println(utils.GetBuildInfo())
	if err := config.Init(""); err != nil {
		panic(err)
	}
	app := newApp()
	listenAddr := viper.GetString("listenAddr")
	listenPort := viper.GetString("listenPort")
	if listenPort == "" {listenPort = "8080"}
	addr :=  listenAddr + ":" + listenPort
	app.Run(iris.Addr(addr))
}