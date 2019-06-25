package license

import (
	"github.com/kataras/iris"
	"github.com/spf13/viper"
)

func RegRouter(rg iris.Party) {

	license := rg.Party("/license", func(ctx iris.Context){
		if !viper.GetBool("zstack.license.enable"){
			ctx.StatusCode(600)
		}else {
			ctx.Next()
		}})
	{
		license.Handle("GET","/download/{md5:string regexp(^[0-9A-Fa-f]{32}$)}",download)
		license.Post("/generate/{type:string has([Trial,Paid,OEM,Free,TrialExt,Hybrid,AddOn,HybridTrialExt])}", generate)
	}
}