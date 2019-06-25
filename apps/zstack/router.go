package zstack

import (
	"api_system_iris/apps/zstack/license"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
)

func RegRouter(app *iris.Application) {
	//注册 Macro
	app.Macros().Get("string").RegisterFunc("has", func(validNames []string) func(string) bool {
		return func(paramValue string) bool {
			for _, validName := range validNames {
				if validName == paramValue {
					return true
				}
			}
			return false
		}
	})
	//注册路由组
	routerGroup := app.Party("/zstack",func(ctx iris.Context){
		if !viper.GetBool("zstack.enable"){
			ctx.StatusCode(600)
		}else {
			ctx.Next()
		}})
	{
		license.RegRouter(routerGroup)
	}
}