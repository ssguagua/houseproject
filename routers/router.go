package routers

import (
	"github.com/astaxie/beego"
	"houseproject/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})   //默认调用Get()函数
	//获得区信息
    beego.Router("/api/v1.0/areas", &controllers.AreaController{},"get:GetArea")
	//
	beego.Router("/api/v1.0/houses/index", &controllers.HouseIndexController{},"get:GetHouseIndex")
	//设置session、删除session退出登录
	beego.Router("/api/v1.0/session", &controllers.SessionController{},"get:GetSessionData;delete:DeleteSession")
	//用户注册和登录
	beego.Router("/api/v1.0/users", &controllers.UserController{},"post:PostUser")

}

