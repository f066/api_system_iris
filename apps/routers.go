package apps

import (
	"api_system_iris/apps/zstack"
	"api_system_iris/config"
	"github.com/kataras/iris"
)

func RegRouter(app *iris.Application,logger iris.Handler)  {
	//捕获所有错误
	app.OnAnyErrorCode(logger, func(ctx iris.Context) {
		ctx.JSON(config.GetInfo(ctx, "error", "Oups something went wrong, try again"))
		//ctx.JSON(iris.Map{"code": ctx.GetStatusCode(),"path":ctx.Path(),"ip":ctx.RemoteAddr(),"status": "error", "message": "Oups something went wrong, try again"})
	})
	//404 错误路由
	app.OnErrorCode(iris.StatusNotFound,logger, func(ctx iris.Context) {
		ctx.JSON(config.GetInfo(ctx,"error","404 Not Found"))
		//ctx.JSON(iris.Map{"code":"404","path":ctx.Path(),"ip":ctx.RemoteAddr(),"status":"error","message":"404 Not Found"})
	})
	//5xx 错误路由
	app.OnErrorCode(iris.StatusInternalServerError,logger, func(ctx iris.Context) {
		ctx.JSON(config.GetInfo(ctx,"error","服务器内部错误"))
		//ctx.JSON(iris.Map{"code":"500","path":ctx.Path(),"ip":ctx.RemoteAddr(),"status":"error","message":"服务器内部错误"})
	})
	//6xx 自定义错误
	app.OnErrorCode(600, func(ctx iris.Context) {
		ctx.JSON(config.GetInfo(ctx,"info","该功能暂时被关闭"))
	})
	//注册应用路由
	zstack.RegRouter(app)

}