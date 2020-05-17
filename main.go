package main

import (
	"github.com/astaxie/beego"
	_ "houseproject/routers"
	"github.com/astaxie/beego/context"
	"strings"
	"net/http"
	_"houseproject/models"
)

func ignoreStaticPath() {
	//透明static
	//beego框架默认放在static目录下面
	//通过下面设置静态路径,重定向
	//beego.SetStaticPath("group1/M00/","static/upload/avatar")

	//beego.InsertFilter实现路由过滤
	//参数分别为路由规则、执行过滤的地方(BeforeRouter 寻找路由之前)、跳转的方法、还有一个可选参数默认为false(过滤输出后仍执行其他过滤)
	//过滤网站请求，所有8080端口及以下的地址都会转向TransparentStatic函数
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)

}
func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
	//将原来的路径+"static/html/",然后再匹配
}
func main() {
	ignoreStaticPath()

	//设置session需要设置：
	//main.go中设置beego.BConfig.WebConfig.Session.SessionOn = true
	//或者在conf文件中设置sessionon=true

	beego.Run()
	//beego.Run(":8899")     //可以修改端口号，并相对于conf文件中优先使用
}

